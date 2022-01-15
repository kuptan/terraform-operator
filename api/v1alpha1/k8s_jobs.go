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
	tfVarsMountPath           string = "/tmp/tfvars"
	moduleWorkingDirMountPath string = "/tmp/tfmodule"
	conifgMapModuleMountPath  string = "/terraform/modules"
	gitSSHKeyMountPath        string = "/root/.ssh"

	knownHostsVolumeName string = "known-hosts"
	emptyDirVolumeName   string = "tfmodule"
	gitSSHKeyVolumeName  string = "git-ssh"
)

func getTerraformRunnerDockerImage() string {
	return fmt.Sprintf("%s/%s:%s", utils.Env.DockerRepository, utils.Env.TerraformRunnerImage, utils.Env.TerraformRunnerImageTag)
}

func getBusyboxDockerImage() string {
	return fmt.Sprintf("%s/%s", utils.Env.DockerRepository, "busybox")
}

// returns a string that could or could not start with a TF_VAR_ prefix for the container
func getEnvVarKey(v Variable) string {
	prefix := ""

	if !v.EnvironmentVariable {
		prefix = "TF_VAR_"
	}

	return fmt.Sprintf("%s%s", prefix, v.Key)
}

// returns a list of environment variables to injected to the job runner
// these environment variables are specific to the terraform runner container
func (t *Terraform) getRunnerSpecificEnvVars() []corev1.EnvVar {
	envVars := []corev1.EnvVar{}

	envVars = append(envVars, getEnvVariable("TERRAFORM_VERSION", t.Spec.TerraformVersion))
	envVars = append(envVars, getEnvVariable("TERRAFORM_WORKING_DIR", moduleWorkingDirMountPath))
	envVars = append(envVars, getEnvVariable("TERRAFORM_VAR_FILES_PATH", tfVarsMountPath))
	envVars = append(envVars, getEnvVariable("OUTPUT_SECRET_NAME", getUniqueResourceName(t.Name, t.Status.RunId)))
	envVars = append(envVars, getEnvVariable("TERRAFORM_DESTROY", strconv.FormatBool(t.Spec.Destroy)))

	envVars = append(envVars, getEnvVariableFromFieldSelector("POD_NAMESPACE", "metadata.namespace"))

	if t.Spec.Workspace != "" {
		envVars = append(envVars, getEnvVariable("TERRAFORM_WORKSPACE", t.Spec.Workspace))
	}

	return envVars
}

// returns Kubernetes Pod environment variables to be passed to the run job
func (t *Terraform) getEnvVariables() []corev1.EnvVar {
	vars := []corev1.EnvVar{}

	for _, v := range t.Spec.Variables {
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

	vars = append(vars, t.getRunnerSpecificEnvVars()...)

	return vars
}

func (t *Terraform) getRunnerSpecificVolumes() []corev1.Volume {
	volumes := []corev1.Volume{}

	name := getUniqueResourceName(t.Name, t.Status.RunId)

	volumes = append(volumes, getEmptyDirVolume(emptyDirVolumeName))
	volumes = append(volumes, getVolumeSpecFromConfigMap(name, name))

	if t.Spec.GitSSHKey != nil && t.Spec.GitSSHKey.ValueFrom != nil {
		volumes = append(volumes, getVolumeSpec(gitSSHKeyVolumeName, *t.Spec.GitSSHKey.ValueFrom))
		volumes = append(volumes, getVolumeSpecFromConfigMap(knownHostsVolumeName, utils.Env.KnownHostsConfigMapName))
	}

	return volumes
}

// return the volumes to be mounted
func (t *Terraform) getJobVolumes() []corev1.Volume {
	volumes := []corev1.Volume{}

	for _, file := range t.Spec.VariableFiles {
		volumes = append(volumes, getVolumeSpec(file.Key, *file.ValueFrom))
	}

	volumes = append(volumes, t.getRunnerSpecificVolumes()...)

	return volumes
}

func (t *Terraform) getRunnerSpecificVolumeMounts() []corev1.VolumeMount {
	mounts := []corev1.VolumeMount{}

	mounts = append(mounts, getVolumeMountSpec(emptyDirVolumeName, moduleWorkingDirMountPath, false))
	mounts = append(mounts, getVolumeMountSpec(getUniqueResourceName(t.Name, t.Status.RunId), conifgMapModuleMountPath, false))

	if t.Spec.GitSSHKey != nil && t.Spec.GitSSHKey.ValueFrom != nil {
		sshKeyFileName := "id_rsa"
		sshKnownHostsFileName := "known_hosts"

		sshKeyMountPath := fmt.Sprintf("%s/%s", gitSSHKeyMountPath, sshKeyFileName)
		sshKnownHostsMountPath := fmt.Sprintf("%s/%s", gitSSHKeyMountPath, sshKnownHostsFileName)

		mounts = append(mounts, getVolumeMountSpecWithSubPath(gitSSHKeyVolumeName, sshKeyMountPath, sshKeyFileName, false))
		mounts = append(mounts, getVolumeMountSpecWithSubPath(knownHostsVolumeName, sshKnownHostsMountPath, sshKnownHostsFileName, false))
	}

	return mounts
}

// return the volumes mounts
func (t *Terraform) getJobVolumeMounts() []corev1.VolumeMount {
	mounts := []corev1.VolumeMount{}

	for _, file := range t.Spec.VariableFiles {
		mounts = append(mounts, getVolumeMountSpec(file.Key, tfVarsMountPath, true))
	}

	mounts = append(mounts, t.getRunnerSpecificVolumeMounts()...)

	return mounts
}

// returnrs the initContainers definition for the run job
func getInitContainersSpec(t *Terraform) []corev1.Container {
	containers := []corev1.Container{}

	cpModule := fmt.Sprintf("cp %s/main.tf %s/main.tf", conifgMapModuleMountPath, moduleWorkingDirMountPath)

	commands := []string{
		"/bin/sh",
		"-c",
	}

	args := []string{
		cpModule,
	}

	containers = append(containers, corev1.Container{
		Name:         "busybox",
		Image:        getBusyboxDockerImage(),
		VolumeMounts: t.getRunnerSpecificVolumeMounts(),
		Command:      commands,
		Args:         args,
	})

	return containers
}

// returns a Kubernetes job struct to run the terraform
func getJobSpecForRun(t *Terraform, owner metav1.OwnerReference) *batchv1.Job {

	envVars := t.getEnvVariables()
	volumes := t.getJobVolumes()
	mounts := t.getJobVolumeMounts()

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
					InitContainers:     getInitContainersSpec(t),
					Containers: []corev1.Container{
						{
							Name:            "terraform",
							Image:           getTerraformRunnerDockerImage(),
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

	job := getJobSpecForRun(run, ownerRef)

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
