package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Kubernetes Secrets", func() {
	BeforeEach(func() {
		// Add any setup steps that needs to be executed before each test
	})

	AfterEach(func() {
		// Add any teardown steps that needs to be executed after each test
	})

	Context("Secrets", func() {
		key := types.NamespacedName{
			Name:      "bar",
			Namespace: "default",
		}

		run := &Terraform{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bar",
				Namespace: "default",
			},
			Spec: TerraformSpec{
				TerraformVersion: "1.0.2",
				Module: Module{
					Source:  "IbraheemAlSaady/test/module",
					Version: "0.0.2",
				},
				Destroy:             false,
				DeleteCompletedJobs: false,
			},
			Status: TerraformStatus{
				RunID: "1234",
			},
		}

		It("should create the secret successfully", func() {
			secret, err := createSecretForOutputs(key, run)

			expectedName := "bar-1234"

			Expect(err).ToNot(HaveOccurred())
			Expect(secret).ToNot(BeNil())
			Expect(secret.Name).To(Equal(expectedName))
		})

		It("should fail to create a secret that already exist", func() {
			secret, err := createSecretForOutputs(key, run)

			Expect(err).To(HaveOccurred())
			Expect(secret).To(BeNil())
		})

		It("should delete the resource for run successfully", func() {
			err := deleteSecretByRun(key.Name, key.Namespace, run.Status.RunID)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return an error if the secret does not exist", func() {
			err := deleteSecretByRun(key.Name, key.Namespace, run.Status.RunID)
			Expect(err).To(HaveOccurred())
			Expect(errors.IsNotFound(err)).To(BeTrue())
		})
	})
})
