package v1alpha1

import (
	"context"

	"github.com/kube-champ/terraform-operator/internal/kube"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// List all pods owned by the Run
func (t *Terraform) GetSecretById(namespacedName types.NamespacedName) (*corev1.Secret, error) {
	secrets := kube.ClientSet.CoreV1().Secrets(namespacedName.Namespace)

	name := getUniqueResourceName(namespacedName.Name, t.Status.RunId)

	secret, err := secrets.Get(context.Background(), name, metav1.GetOptions{})

	if err != nil {
		return nil, err
	}

	return secret, err
}

func createSecretForOutputs(namespacedName types.NamespacedName, t *Terraform) (*corev1.Secret, error) {
	secrets := kube.ClientSet.CoreV1().Secrets(namespacedName.Namespace)

	secretName := getUniqueResourceName(namespacedName.Name, t.Status.RunId)

	obj := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:   secretName,
			Labels: getCommonLabels(namespacedName.Name, t.Status.RunId),
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
