package controllers

import (
	"context"
	"time"

	"github.com/kube-champ/terraform-operator/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Terraform Controller", func() {
	const timeout = time.Second * 30
	const interval = time.Second * 1

	BeforeEach(func() {
		// Add any setup steps that needs to be executed before each test
	})

	AfterEach(func() {
		// Add any teardown steps that needs to be executed after each test
	})

	Context("Terraform Run", func() {
		key := types.NamespacedName{
			Name:      "run-sample",
			Namespace: "default",
		}

		created := &v1alpha1.Terraform{
			ObjectMeta: metav1.ObjectMeta{
				Name:      key.Name,
				Namespace: key.Namespace,
			},
			Spec: v1alpha1.TerraformSpec{
				TerraformVersion: "1.0.2",
				Module: v1alpha1.Module{
					Source:  "IbraheemAlSaady/test/module",
					Version: "0.0.1",
				},
				Variables: []v1alpha1.Variable{
					v1alpha1.Variable{
						Key:   "length",
						Value: "16",
					},
				},
				Destroy:             false,
				DeleteCompletedJobs: false,
			},
		}

		It("Should create successfully", func() {
			// Create
			Expect(k8sClient.Create(context.Background(), created)).Should(Succeed())

			Eventually(func() string {
				r := &v1alpha1.Terraform{}
				k8sClient.Get(context.Background(), key, r)

				return r.Status.RunId
			}, timeout, interval).ShouldNot(BeEmpty())

			By("expect status to be started")
			Eventually(func() v1alpha1.TerraformRunStatus {
				r := &v1alpha1.Terraform{}
				k8sClient.Get(context.Background(), key, r)

				return r.Status.RunStatus
			}, timeout, interval).Should(Equal(v1alpha1.RunStarted))

			By("expect status to be running")
			Eventually(func() v1alpha1.TerraformRunStatus {
				r := &v1alpha1.Terraform{}
				k8sClient.Get(context.Background(), key, r)

				makeRunJobRunning(r)

				return r.Status.RunStatus
			}, timeout, interval).Should(Equal(v1alpha1.RunRunning))

			By("expect status to be completed")
			Eventually(func() v1alpha1.TerraformRunStatus {
				r := &v1alpha1.Terraform{}
				k8sClient.Get(context.Background(), key, r)

				makeRunJobSucceed(r)

				return r.Status.RunStatus
			}, timeout, interval).Should(Equal(v1alpha1.RunCompleted))
		})

		It("should update successfully", func() {
			// Update
			updated := &v1alpha1.Terraform{}
			Expect(k8sClient.Get(context.Background(), key, updated)).Should(Succeed())

			updated.Spec.Workspace = "dev"
			Expect(k8sClient.Update(context.Background(), updated)).Should(Succeed())

			By("expect status to be running")
			Eventually(func() v1alpha1.TerraformRunStatus {
				r := &v1alpha1.Terraform{}
				k8sClient.Get(context.Background(), key, r)

				makeRunJobRunning(r)

				return r.Status.RunStatus
			}, timeout, interval).Should(Equal(v1alpha1.RunRunning))

			By("expect status to be completed")
			Eventually(func() v1alpha1.TerraformRunStatus {
				r := &v1alpha1.Terraform{}
				k8sClient.Get(context.Background(), key, r)

				makeRunJobSucceed(r)

				return r.Status.RunStatus
			}, timeout, interval).Should(Equal(v1alpha1.RunCompleted))
		})

		It("should have a failed status if job failed", func() {
			updated := &v1alpha1.Terraform{}
			Expect(k8sClient.Get(context.Background(), key, updated)).Should(Succeed())

			updated.Spec.Workspace = "prod"
			Expect(k8sClient.Update(context.Background(), updated)).Should(Succeed())

			Eventually(func() v1alpha1.TerraformRunStatus {
				r := &v1alpha1.Terraform{}
				k8sClient.Get(context.Background(), key, r)

				makeRunJobFail(r)

				return r.Status.RunStatus
			}, timeout, interval).Should(Equal(v1alpha1.RunFailed))
		})

		It("should create with destroy and deletes completed jobs", func() {
			updated := &v1alpha1.Terraform{}
			Expect(k8sClient.Get(context.Background(), key, updated)).Should(Succeed())

			updated.Spec.Destroy = true
			updated.Spec.DeleteCompletedJobs = true
			Expect(k8sClient.Update(context.Background(), updated)).Should(Succeed())

			By("Expecting status to be started")
			Eventually(func() v1alpha1.TerraformRunStatus {
				r := &v1alpha1.Terraform{}
				k8sClient.Get(context.Background(), key, r)

				return r.Status.RunStatus
			}, timeout, interval).Should(Equal(v1alpha1.RunStarted))

			By("Expecting status to be completed")
			Eventually(func() v1alpha1.TerraformRunStatus {
				r := &v1alpha1.Terraform{}
				k8sClient.Get(context.Background(), key, r)

				makeRunJobSucceed(r)

				return r.Status.RunStatus
			}, timeout, interval).Should(Equal(v1alpha1.RunCompleted))

			By("Expecting job to be deleted")
			Eventually(func() bool {
				r := &v1alpha1.Terraform{}
				k8sClient.Get(context.Background(), key, r)

				return isJobDeleted(r)
			}, timeout, interval).Should(BeTrue())
		})

		It("should delete successfully", func() {
			// Delete
			By("Expecting to delete successfully")
			Eventually(func() error {
				r := &v1alpha1.Terraform{}
				k8sClient.Get(context.Background(), key, r)
				return k8sClient.Delete(context.Background(), r)
			}, timeout, interval).Should(Succeed())

			By("Expecting to delete finish")
			Eventually(func() error {
				r := &v1alpha1.Terraform{}
				return k8sClient.Get(context.Background(), key, r)
			}, timeout, interval).ShouldNot(Succeed())
		})
	})

	Context("Terraform Run Dependencies", func() {
		run1Key := types.NamespacedName{
			Name:      "run-dep1",
			Namespace: "default",
		}

		run2Key := types.NamespacedName{
			Name:      "run-dep2",
			Namespace: "default",
		}

		run1 := &v1alpha1.Terraform{
			ObjectMeta: metav1.ObjectMeta{
				Name:      run1Key.Name,
				Namespace: run1Key.Namespace,
			},
			Spec: v1alpha1.TerraformSpec{
				TerraformVersion: "1.0.2",
				Module: v1alpha1.Module{
					Source:  "IbraheemAlSaady/test/module",
					Version: "0.0.1",
				},
			},
		}

		run2 := &v1alpha1.Terraform{
			ObjectMeta: metav1.ObjectMeta{
				Name:      run2Key.Name,
				Namespace: run2Key.Namespace,
			},
			Spec: v1alpha1.TerraformSpec{
				TerraformVersion: "1.0.2",
				Module: v1alpha1.Module{
					Source:  "IbraheemAlSaady/test/module",
					Version: "0.0.1",
				},
				DependsOn: []*v1alpha1.DependsOn{
					&v1alpha1.DependsOn{
						Name:      run1Key.Name,
						Namespace: run1Key.Namespace,
					},
				},
			},
		}

		It("Should create runs successfully with the correct dependency flow", func() {
			// Create
			Expect(k8sClient.Create(context.Background(), run1)).Should(Succeed())

			By("evaluating run1 runId")
			Eventually(func() string {
				r := &v1alpha1.Terraform{}
				k8sClient.Get(context.Background(), run1Key, r)

				return r.Status.RunId
			}, timeout, interval).ShouldNot(BeEmpty())

			// Create
			Expect(k8sClient.Create(context.Background(), run2)).Should(Succeed())

			By("evaluating run2 status")
			Eventually(func() v1alpha1.TerraformRunStatus {
				r := &v1alpha1.Terraform{}
				k8sClient.Get(context.Background(), run2Key, r)

				return r.Status.RunStatus
			}, timeout, interval).Should(Equal(v1alpha1.RunWaitingForDependency))

			By("expect run1 status to be completed")
			Eventually(func() v1alpha1.TerraformRunStatus {
				r := &v1alpha1.Terraform{}
				k8sClient.Get(context.Background(), run1Key, r)

				makeRunJobSucceed(r)

				return r.Status.RunStatus
			}, timeout, interval).Should(Equal(v1alpha1.RunCompleted))

			By("evaluating run2 runId")
			Eventually(func() string {
				r := &v1alpha1.Terraform{}
				k8sClient.Get(context.Background(), run2Key, r)

				return r.Status.RunId
			}, timeout, interval).ShouldNot(BeEmpty())

			By("expect run2 status to be completed")
			Eventually(func() v1alpha1.TerraformRunStatus {
				r := &v1alpha1.Terraform{}
				k8sClient.Get(context.Background(), run2Key, r)

				makeRunJobSucceed(r)

				return r.Status.RunStatus
			}, timeout, interval).Should(Equal(v1alpha1.RunCompleted))
		})
	})
})
