package v1alpha1

import (
	"context"
	"fmt"
	"strconv"

	"github.com/kube-champ/terraform-operator/pkg/kube"
	"github.com/kube-champ/terraform-operator/pkg/utils"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	tfVarsMountPath          string = "/tmp/tfvars"
	moduleWorkingDir         string = "/tmp/tfmodule"
	conifgMapModuleMountPath string = "/terraform/modules"
	emptyDirVolumeName       string = "tfmodule"
)

func getTerraformDockerImage() string {
	return fmt.Sprintf("%s/%s:%s", utils.Env.DockerRepository, utils.Env.TerraformRunnerImage, utils.Env.TerraformRunnerImageTag)
}

// return the volumes to be mounted
func getJobVolumes(varFiles []VariableFile) []corev1.Volume {
	volumes := []corev1.Volume{}

	for _, file := range varFiles {
		volumes = append(volumes, getVolumeSpec(file.Key, *file.ValueFrom))
	}

	return volumes
}

// return the volumes mounts
func getJobVolumeMounts(varFiles []VariableFile) []corev1.VolumeMount {
	mounts := []corev1.VolumeMount{}

	for _, file := range varFiles {
		mounts = append(mounts, getVolumeMountSpec(file.Key, tfVarsMountPath, true))
	}

	return mounts
}

// returns a string that could or could not start with a TF_VAR_ prefix for the container
func getEnvVarKey(v Variable) string {
	prefix := ""

	if !v.EnvironmentVariable {
		prefix = "TF_VAR_"
	}

	return fmt.Sprintf("%s%s", prefix, v.Key)
}

// returns Kubernetes Pod environment variables to be passed to the run job
func getEnvVariables(variables []Variable) []corev1.EnvVar {
	vars := []corev1.EnvVar{}

	for _, v := range variables {
		if v.ValueFrom != nil {
			vars = append(vars, corev1.EnvVar{
				Name:      getEnvVarKey(v),
				ValueFrom: v.ValueFrom,
			})
		}

		if v.Value != "" {
			vars = append(vars, corev1.EnvVar{
				Name:  getEnvVarKey(v),
				Value: v.Value,
			})
		}
	}

	return vars
}

// returns a list of environment variables to injected to the job runner
// these environment variables are specific to the terraform runner container
func (t *Terraform) getRunnerSpecificEnvVars(secret *corev1.Secret) []corev1.EnvVar {
	envVars := []corev1.EnvVar{}

	envVars = append(envVars, getEnvVariable("TERRAFORM_VERSION", t.Spec.TerraformVersion))
	envVars = append(envVars, getEnvVariable("TERRAFORM_WORKING_DIR", moduleWorkingDir))
	envVars = append(envVars, getEnvVariable("TERRAFORM_VAR_FILES_PATH", tfVarsMountPath))
	envVars = append(envVars, getEnvVariable("OUTPUT_SECRET_NAME", secret.ObjectMeta.Name))
	envVars = append(envVars, getEnvVariable("TERRAFORM_DESTROY", strconv.FormatBool(t.Spec.Destroy)))

	envVars = append(envVars, getEnvVariableFromFieldSelector("POD_NAMESPACE", "metadata.namespace"))

	if t.Spec.Workspace != "" {
		envVars = append(envVars, getEnvVariable("TERRAFORM_WORKSPACE", t.Spec.Workspace))
	}

	return envVars
}

func (t *Terraform) getRunnerSpecificVolumes(configMap *corev1.ConfigMap) []corev1.Volume {
	volumes := []corev1.Volume{}

	volumes = append(volumes, getEmptyDirVolume(emptyDirVolumeName))

	return append(volumes, corev1.Volume{
		Name: configMap.ObjectMeta.Name,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: configMap.ObjectMeta.Name,
				},
			},
		},
	})
}

func (t *Terraform) getRunnerSpecificVolumeMounts(configMap *corev1.ConfigMap) []corev1.VolumeMount {
	mounts := []corev1.VolumeMount{}

	mounts = append(mounts, getVolumeMountSpec(emptyDirVolumeName, moduleWorkingDir, false))
	mounts = append(mounts, getVolumeMountSpec(configMap.Name, conifgMapModuleMountPath, false))

	return mounts
}

func getInitContainersSpec(configMap *corev1.ConfigMap, t *Terraform) []corev1.Container {
	containers := []corev1.Container{}

	containers = append(containers, corev1.Container{
		Name:         "busybox",
		Image:        "busybox",
		VolumeMounts: t.getRunnerSpecificVolumeMounts(configMap),
		Command: []string{
			"/bin/sh",
			"-c",
			fmt.Sprintf("cp -a %s/. %s/", conifgMapModuleMountPath, moduleWorkingDir),
		},
	})

	return containers
}

// returns a Kubernetes job struct to run the terraform
func getJobSpecForRun(
	t *Terraform,
	configMap *corev1.ConfigMap,
	secret *corev1.Secret,
	owner metav1.OwnerReference) *batchv1.Job {

	envVars := append(getEnvVariables(t.Spec.Variables), t.getRunnerSpecificEnvVars(secret)...)
	volumes := append(getJobVolumes(t.Spec.VariableFiles), t.getRunnerSpecificVolumes(configMap)...)
	mounts := append(getJobVolumeMounts(t.Spec.VariableFiles), t.getRunnerSpecificVolumeMounts(configMap)...)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getUniqueResourceName(t.Name, t.Status.RunId),
			Namespace: t.Namespace,
			Labels:    getCommonLabels(t.Name, t.Status.RunId),
			OwnerReferences: []metav1.OwnerReference{
				owner,
			},
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: getCommonLabels(t.Name, t.Status.RunId),
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: "terraform-runner",
					InitContainers:     getInitContainersSpec(configMap, t),
					Containers: []corev1.Container{
						{
							Name:            "terraform",
							Image:           getTerraformDockerImage(),
							VolumeMounts:    mounts,
							Env:             envVars,
							ImagePullPolicy: corev1.PullIfNotPresent,
						},
					},
					Volumes:       volumes,
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
		},
	}

	job.Spec.BackoffLimit = &t.Spec.RetryLimit

	return job
}

// Gets the kubernetes job for a specific run
func getJobForRun(runName string, namespace string, runId string) (*batchv1.Job, error) {
	jobs := kube.ClientSet.BatchV1().Jobs(namespace)

	name := getUniqueResourceName(runName, runId)

	job, err := jobs.Get(context.Background(), name, metav1.GetOptions{})

	if err != nil {
		return nil, err
	}

	return job, err
}

func createJobForRun(run *Terraform, configMap *corev1.ConfigMap, secret *corev1.Secret) (*batchv1.Job, error) {
	jobs := kube.ClientSet.BatchV1().Jobs(run.Namespace)

	ownerRef := run.GetOwnerReference()

	job := getJobSpecForRun(run, configMap, secret, ownerRef)

	if _, err := jobs.Create(context.TODO(), job, metav1.CreateOptions{}); err != nil {
		return nil, err
	}

	return job, nil
}

// Deletes the job run
func deleteJobByRun(runName string, namespace string, runId string) error {
	jobs := kube.ClientSet.BatchV1().Jobs(namespace)

	resourceName := getUniqueResourceName(runName, runId)

	deletePolicy := metav1.DeletePropagationForeground

	if err := jobs.Delete(context.Background(), resourceName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		return err
	}

	return nil
}
