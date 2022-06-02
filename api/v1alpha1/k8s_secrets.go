package v1alpha1

import (
	"context"

	"github.com/kuptan/terraform-operator/internal/kube"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// createSecretForOutputs creates a secret to store the the Terraform output of the workflow/run
func createSecretForOutputs(namespacedName types.NamespacedName, t *Terraform) (*corev1.Secret, error) {
	secrets := kube.ClientSet.CoreV1().Secrets(namespacedName.Namespace)

	secretName := getUniqueResourceName(namespacedName.Name, t.Status.RunID)

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

	secret, err := secrets.Create(context.Background(), obj, metav1.CreateOptions{})

	if err != nil {
		return nil, err
	}

	return secret, nil
}
