package v1alpha1

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

	key := types.NamespacedName{
		Name:      "foo",
		Namespace: "default",
	}

	var created, fetched *Terraform

	Context("Create run", func() {
		It("should create a Terraform Object", func() {
			created = &Terraform{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
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

			run2.Status.Generation = run2.Generation
			run2.Generation = 2

			By("run generation was updated")
			Expect(run2.IsUpdated()).To(BeTrue())

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
})
