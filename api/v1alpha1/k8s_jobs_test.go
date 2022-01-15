package v1alpha1

import (
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Kubernetes RBAC", func() {
	BeforeEach(func() {
		// Add any setup steps that needs to be executed before each test
	})

	AfterEach(func() {
		// Add any teardown steps that needs to be executed after each test
	})

	Context("Jobs", func() {
		var job *batchv1.Job

		run := &Terraform{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bar",
				Namespace: "default",
			},
			Spec: TerraformSpec{
				TerraformVersion: "1.0.2",
				Module: Module{
					Source:  "IbraheemAlSaady/test/module",
					Version: "0.0.1",
				},
				GitSSHKey: &GitSSHKeySpec{
					ValueFrom: &corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: "mysecret",
						},
					},
				},
				Destroy:             false,
				DeleteCompletedJobs: false,
			},
			Status: TerraformStatus{
				RunId: "12345",
			},
		}

		ownerRef := metav1.OwnerReference{
			APIVersion: fmt.Sprintf("%s/%s", GroupVersion.Group, GroupVersion.Version),
			Kind:       "terraform",
			Name:       "foot",
			UID:        "1234",
		}

		It("returns the job spec and should not be null", func() {
			jobSpec := getJobSpecForRun(run, ownerRef)

			Expect(jobSpec).ToNot(BeNil())

			job = jobSpec
		})

		It("should contain a volume for the git ssh", func() {
			var sshVolume *corev1.Volume

			for _, v := range job.Spec.Template.Spec.Volumes {
				if v.Name == gitSSHKeyVolumeName {
					sshVolume = &v
					break
				}
			}

			Expect(sshVolume).ToNot(BeNil())
			Expect(sshVolume.Name).To(Equal(gitSSHKeyVolumeName))
			Expect(sshVolume.VolumeSource.Secret.SecretName).To(Equal(run.Spec.GitSSHKey.ValueFrom.Secret.SecretName))
		})
	})
})
