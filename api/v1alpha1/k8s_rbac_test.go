package v1alpha1

import (
	"context"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kube-champ/terraform-operator/pkg/kube"
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

	rbacName := "terraform-runner"
	namespace := "default"

	Context("RBAC", func() {
		It("should create a cluster role", func() {
			kube.ClientSet.RbacV1().ClusterRoles().Create(context.Background(), &rbacv1.ClusterRole{
				ObjectMeta: metav1.ObjectMeta{
					Name: rbacName,
				},
				Rules: []rbacv1.PolicyRule{
					rbacv1.PolicyRule{
						APIGroups: []string{},
						Resources: []string{"secrets"},
						Verbs:     []string{"get", "update"},
					},
				},
			}, metav1.CreateOptions{})
		})

		It("service account should not be found", func() {
			found, err := isServiceAccountExist(namespace)

			Expect(err).ToNot(HaveOccurred())
			Expect(found).To(BeFalse())
		})

		It("role binding should not be found", func() {
			found, err := isRoleBindingExist(namespace)

			Expect(err).ToNot(HaveOccurred())
			Expect(found).To(BeFalse())
		})

		It("should create service account and role binding", func() {
			err := createRbacConfigIfNotExist(namespace)
			Expect(err).ToNot(HaveOccurred())

			sa, err := kube.ClientSet.CoreV1().ServiceAccounts(namespace).Get(context.Background(), rbacName, metav1.GetOptions{})

			Expect(err).ToNot(HaveOccurred())
			Expect(sa.Name).To(Equal(rbacName))

			roleBinding, err := kube.ClientSet.RbacV1().RoleBindings(namespace).Get(context.Background(), rbacName, metav1.GetOptions{})

			Expect(err).ToNot(HaveOccurred())
			Expect(roleBinding.Name).To(Equal(rbacName))
		})
	})
})
