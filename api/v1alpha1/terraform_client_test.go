package v1alpha1

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("TerraformClient", func() {
	BeforeEach(func() {
		// Add any setup steps that needs to be executed before each test
	})

	AfterEach(func() {
		// Add any teardown steps that needs to be executed after each test
	})

	Context("Kubernetes Terraform Client", func() {
		ns := "default"

		key := types.NamespacedName{
			Name:      "run-client",
			Namespace: ns,
		}

		run := &Terraform{
			ObjectMeta: metav1.ObjectMeta{
				Name:      key.Name,
				Namespace: key.Namespace,
				Labels: map[string]string{
					"createdWith": "terraform-client",
				},
			},
			Spec: TerraformSpec{
				TerraformVersion: "1.0.2",
				Module: Module{
					Source:  "IbraheemAlSaady/test/module",
					Version: "0.0.1",
				},
			},
		}

		It("should create a Terraform object", func() {
			tf, err := terraformKubeClient.Terraforms(ns).Create(context.Background(), run, metav1.CreateOptions{})

			Expect(err).ToNot(HaveOccurred())
			Expect(tf).ToNot(BeNil())
		})

		It("should get the Terraform object", func() {
			tf, err := terraformKubeClient.Terraforms(ns).Get(context.Background(), run.Name, metav1.GetOptions{})

			Expect(err).ToNot(HaveOccurred())
			Expect(tf).ToNot(BeNil())
		})

		It("should list the Terraform object", func() {
			tf, err := terraformKubeClient.Terraforms(ns).List(context.Background(), metav1.ListOptions{
				LabelSelector: "createdWith",
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(tf.Items).To(HaveLen(1))
		})

		It("should watch the Terraform object", func() {
			tf, err := terraformKubeClient.Terraforms(ns).Watch(context.Background(), metav1.ListOptions{
				LabelSelector: "createdWith",
			})

			Expect(err).ToNot(HaveOccurred())
			tf.Stop()
		})
	})
})
