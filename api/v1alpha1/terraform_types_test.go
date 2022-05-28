package v1alpha1

import (
	"context"

	"github.com/kube-champ/terraform-operator/internal/kube"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("TerraformRun", func() {
	BeforeEach(func() {
		// Add any setup steps that needs to be executed before each test
	})

	AfterEach(func() {
		// Add any teardown steps that needs to be executed after each test
	})

	Context("Create run", func() {
		var created, fetched *Terraform

		key := types.NamespacedName{
			Name:      "foo",
			Namespace: "default",
		}

		It("should create a Terraform Object", func() {
			created = &Terraform{
				ObjectMeta: metav1.ObjectMeta{
					Name:      key.Name,
					Namespace: key.Namespace,
				},
				Spec: TerraformSpec{
					TerraformVersion: "1.0.2",
					Module: Module{
						Source:  "IbraheemAlSaady/test/module",
						Version: "0.0.1",
					},
					Variables: []Variable{
						Variable{
							Key:   "length",
							Value: "16",
						},
					},
					Destroy:             false,
					DeleteCompletedJobs: false,
				},
			}

			By("creating a terraform run")
			Expect(k8sClient.Create(context.TODO(), created)).To(Succeed())

			fetched = &Terraform{}
			Expect(k8sClient.Get(context.TODO(), key, fetched)).To(Succeed())
			Expect(fetched).To(Equal(created))

			By("deleting the created object")
			Expect(k8sClient.Delete(context.TODO(), created)).To(Succeed())
			Expect(k8sClient.Get(context.TODO(), key, created)).ToNot(Succeed())
		})

		It("should correctly handle run statuses", func() {
			run1 := &Terraform{
				Status: TerraformStatus{
					RunId: "",
				},
			}

			By("run was just submitted")
			Expect(run1.IsSubmitted()).To(BeTrue())

			run2 := &Terraform{
				ObjectMeta: metav1.ObjectMeta{
					Generation: 1,
				},
				Status: TerraformStatus{
					RunStatus: RunStarted,
				},
			}

			run2.SetRunId()
			Expect(run2.Status.RunId).To(HaveLen(8))

			By("run is now in a Started state")
			Expect(run2.IsSubmitted()).To(BeFalse())
			Expect(run2.IsStarted()).To(BeTrue())

			run2.Status.RunStatus = RunRunning
			By("run is now in a Running state")
			Expect(run2.IsStarted()).To(BeTrue())
			Expect(run2.IsRunning()).To(BeTrue())

			run2.Status.RunStatus = RunFailed
			By("run is now in a Failed state")
			Expect(run2.HasErrored()).To(BeTrue())

			run2.Status.ObservedGeneration = run2.Generation
			run2.Generation = 2

			By("run generation was updated")
			Expect(run2.IsUpdated()).To(BeTrue())

			run2.Status.RunStatus = RunWaitingForDependency
			By("run is now in a waiting state")
			Expect(run2.IsWaiting()).To(BeTrue())

			run2.Status.RunStatus = RunCompleted
			By("run is now in a Completed state")
			Expect(run2.IsStarted()).To(BeFalse())
		})

		It("should handle a terraform run job", func() {
			run := fetched.DeepCopy()

			run.Status.RunId = "1234"

			job, err := run.CreateTerraformRun(key)
			Expect(err).ToNot(HaveOccurred(), "failed to create a terraform run")
			Expect(job.Name).ToNot(BeEmpty())

			job, err = run.GetJobByRun()

			Expect(err).ToNot(HaveOccurred(), "run job was not found")
			Expect(job.Name).ToNot(BeEmpty())

			err = run.DeleteAfterCompletion()

			Expect(err).ToNot(HaveOccurred(), "run job could not be deleted")
		})

		It("should get the owner preference", func() {
			run := fetched.DeepCopy()

			owner := run.GetOwnerReference()

			Expect(owner).ToNot(BeNil())
		})

		It("should handle previous statuses", func() {
			run := &Terraform{
				Status: TerraformStatus{
					RunId:     "1234",
					RunStatus: RunCompleted,
				},
			}

			run.PrepareForUpdate()

			Expect(run.Status.PreviousRuns).To(HaveLen(1))
			Expect(run.Status.PreviousRuns[0].RunId).To(Equal("1234"))
			Expect(run.Status.PreviousRuns[0].Status).To(Equal(RunCompleted))
		})
	})

	Context("Run Errors", func() {
		key := types.NamespacedName{
			Name:      "barbar",
			Namespace: "default",
		}

		run := &Terraform{
			ObjectMeta: metav1.ObjectMeta{
				Name:      key.Name,
				Namespace: key.Namespace,
			},
			Spec: TerraformSpec{
				TerraformVersion: "1.0.2",
				Module: Module{
					Source:  "IbraheemAlSaady/test/module",
					Version: "0.0.1",
				},
				Variables: []Variable{
					Variable{
						Key:   "length",
						Value: "16",
					},
				},
				Destroy:             false,
				DeleteCompletedJobs: false,
			},
			Status: TerraformStatus{
				RunId: "1234",
			},
		}

		name := getUniqueResourceName(run.Name, run.Status.RunId)

		It("should fail to create a run due to existing configmap", func() {
			cfg := corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: "default",
				},
				Data: make(map[string]string),
			}

			kube.ClientSet.CoreV1().ConfigMaps("default").Create(context.Background(), &cfg, metav1.CreateOptions{})

			job, err := run.CreateTerraformRun(key)

			Expect(err).To(HaveOccurred())
			Expect(job).To(BeNil())

			kube.ClientSet.CoreV1().ConfigMaps("default").Delete(context.Background(), name, metav1.DeleteOptions{})
		})

		It("should fail to create a run due to existing secret", func() {
			secret := corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: "default",
				},
				Data: make(map[string][]byte),
			}

			kube.ClientSet.CoreV1().Secrets("default").Create(context.Background(), &secret, metav1.CreateOptions{})

			job, err := run.CreateTerraformRun(key)

			Expect(err).To(HaveOccurred())
			Expect(job).To(BeNil())

			kube.ClientSet.CoreV1().Secrets("default").Delete(context.Background(), name, metav1.DeleteOptions{})
		})

		It("should fail to create a run due to existing job", func() {
			j := batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: "default",
				},
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								corev1.Container{
									Name:  "busybox",
									Image: "busybox",
								},
							},
							RestartPolicy: corev1.RestartPolicyNever,
						},
					},
				},
			}

			kube.ClientSet.BatchV1().Jobs("default").Create(context.Background(), &j, metav1.CreateOptions{})

			job, err := run.CreateTerraformRun(key)

			Expect(err).To(HaveOccurred())
			Expect(job).To(BeNil())

			kube.ClientSet.BatchV1().Jobs("default").Delete(context.Background(), name, metav1.DeleteOptions{})
		})

		It("should return error if the job does not exist", func() {
			job, err := run.GetJobByRun()

			Expect(err).To(HaveOccurred())
			Expect(job).To(BeNil())
		})

		It("should fail to delete a job that does not exist", func() {
			err := run.DeleteAfterCompletion()

			Expect(err).To(HaveOccurred())
		})
	})
})
