package v1alpha1

import (
	"context"

	"github.com/kube-champ/terraform-operator/pkg/kube"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const runnerRBACName string = "terraform-runner"

func createServiceAccount(namespace string) (*corev1.ServiceAccount, error) {
	key := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      runnerRBACName,
			Namespace: namespace,
		},
	}

	sa, err := kube.ClientSet.CoreV1().ServiceAccounts(namespace).Create(context.Background(), key, metav1.CreateOptions{})

	if err != nil {
		return nil, err
	}

	return sa, nil
}

func createRoleBinding(namespace string) (*rbacv1.RoleBinding, error) {
	key := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      runnerRBACName,
			Namespace: namespace,
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     runnerRBACName,
			APIGroup: "rbac.authorization.k8s.io",
		},
		Subjects: []rbacv1.Subject{
			rbacv1.Subject{
				Kind:      "ServiceAccount",
				Name:      runnerRBACName,
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

func isServiceAccountExist(namespace string) (bool, error) {
	_, err := kube.ClientSet.CoreV1().ServiceAccounts(namespace).Get(context.Background(), runnerRBACName, metav1.GetOptions{})

	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func isRoleBindingExist(namespace string) (bool, error) {
	_, err := kube.ClientSet.RbacV1().RoleBindings(namespace).Get(context.Background(), runnerRBACName, metav1.GetOptions{})

	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func createRbacConfigIfNotExist(namespace string) error {
	saExist, err := isServiceAccountExist(namespace)

	if err != nil {
		return err
	}

	roleBindingExist, err := isRoleBindingExist(namespace)

	if err != nil {
		return err
	}

	if !saExist && !roleBindingExist {
		if _, err := createServiceAccount(namespace); err != nil {
			return err
		}

		if _, err := createRoleBinding(namespace); err != nil {
			return err
		}
	}

	return nil
}
