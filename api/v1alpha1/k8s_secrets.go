package v1alpha1

import (
	"context"

	"github.com/kuptan/terraform-operator/internal/kube"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// isSecretExist checks whether a Secret exist
func isSecretExist(ctx context.Context, name string, namespace string) (*corev1.Secret, error) {
	secret, err := kube.ClientSet.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})

	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return secret, nil
}

// createSecretForOutputs creates a secret to store the the Terraform output of the workflow/run
func createSecretForOutputs(ctx context.Context, namespacedName types.NamespacedName, t *Terraform) (*corev1.Secret, error) {
	secretName := getOutputSecretname(namespacedName.Name)

	exist, err := isSecretExist(ctx, secretName, namespacedName.Namespace)

	if err != nil {
		return nil, err
	}

	if exist != nil {
		return exist, nil
	}

	secrets := kube.ClientSet.CoreV1().Secrets(namespacedName.Namespace)

	obj := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:   secretName,
			Labels: getCommonLabels(namespacedName.Name, t.Status.RunID),
			OwnerReferences: []metav1.OwnerReference{
				t.GetOwnerReference(),
			},
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{},
	}

	secret, err := secrets.Create(ctx, obj, metav1.CreateOptions{})

	if err != nil {
		return nil, err
	}

	return secret, nil
}
