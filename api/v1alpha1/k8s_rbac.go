package v1alpha1

import (
	"context"

	"github.com/kuptan/terraform-operator/internal/kube"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// createServiceAccount creates a Kubernetes ServiceAccount for the Terraform Runner
func createServiceAccount(ctx context.Context, name string, namespace string) (*corev1.ServiceAccount, error) {
	key := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}

	sa, err := kube.ClientSet.CoreV1().ServiceAccounts(namespace).Create(ctx, key, metav1.CreateOptions{})

	if err != nil {
		return nil, err
	}

	return sa, nil
}

// createRoleBinding creates a Kubernetes RoleBinding for the Terraform Runner
func createRoleBinding(ctx context.Context, name string, namespace string) (*rbacv1.RoleBinding, error) {
	key := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     name,
			APIGroup: "rbac.authorization.k8s.io",
		},
		Subjects: []rbacv1.Subject{
			rbacv1.Subject{
				Kind:      "ServiceAccount",
				Name:      name,
				Namespace: namespace,
			},
		},
	}

	role, err := kube.ClientSet.RbacV1().RoleBindings(namespace).Create(ctx, key, metav1.CreateOptions{})

	if err != nil {
		return nil, err
	}

	return role, nil
}

// isServiceAccountExist checks whether the ServiceAccount for the Terraform Runner exist
func isServiceAccountExist(ctx context.Context, name string, namespace string) (bool, error) {
	_, err := kube.ClientSet.CoreV1().ServiceAccounts(namespace).Get(ctx, name, metav1.GetOptions{})

	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// isRoleBindingExist checks if the RoleBinding for the Terraform Runner exists
func isRoleBindingExist(ctx context.Context, name string, namespace string) (bool, error) {
	_, err := kube.ClientSet.RbacV1().RoleBindings(namespace).Get(ctx, name, metav1.GetOptions{})

	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// createRbacConfigIfNotExist validates if RBAC exist for the Terraform Runner and creates it if not exist
func createRbacConfigIfNotExist(ctx context.Context, name string, namespace string) error {
	saExist, err := isServiceAccountExist(ctx, name, namespace)

	if err != nil {
		return err
	}

	roleBindingExist, err := isRoleBindingExist(ctx, name, namespace)

	if err != nil {
		return err
	}

	if !saExist {
		if _, err := createServiceAccount(ctx, name, namespace); err != nil {
			return err
		}
	}

	if !roleBindingExist {
		if _, err := createRoleBinding(ctx, name, namespace); err != nil {
			return err
		}
	}

	return nil
}
