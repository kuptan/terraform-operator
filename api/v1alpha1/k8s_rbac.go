package v1alpha1

import (
	"context"

	"github.com/kube-champ/terraform-operator/internal/kube"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createServiceAccount(name string, namespace string) (*corev1.ServiceAccount, error) {
	key := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}

	sa, err := kube.ClientSet.CoreV1().ServiceAccounts(namespace).Create(context.Background(), key, metav1.CreateOptions{})

	if err != nil {
		return nil, err
	}

	return sa, nil
}

func createRoleBinding(name string, namespace string) (*rbacv1.RoleBinding, error) {
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

	role, err := kube.ClientSet.RbacV1().RoleBindings(namespace).Create(context.Background(), key, metav1.CreateOptions{})

	if err != nil {
		return nil, err
	}

	return role, nil
}

func isServiceAccountExist(name string, namespace string) (bool, error) {
	_, err := kube.ClientSet.CoreV1().ServiceAccounts(namespace).Get(context.Background(), name, metav1.GetOptions{})

	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func isRoleBindingExist(name string, namespace string) (bool, error) {
	_, err := kube.ClientSet.RbacV1().RoleBindings(namespace).Get(context.Background(), name, metav1.GetOptions{})

	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func createRbacConfigIfNotExist(name string, namespace string) error {
	saExist, err := isServiceAccountExist(name, namespace)

	if err != nil {
		return err
	}

	roleBindingExist, err := isRoleBindingExist(name, namespace)

	if err != nil {
		return err
	}

	if !saExist && !roleBindingExist {
		if _, err := createServiceAccount(name, namespace); err != nil {
			return err
		}

		if _, err := createRoleBinding(name, namespace); err != nil {
			return err
		}
	}

	return nil
}
