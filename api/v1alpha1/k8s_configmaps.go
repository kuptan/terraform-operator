package v1alpha1

import (
	"context"

	"github.com/kube-champ/terraform-operator/internal/kube"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// getConfigMapSpecForModule returns a Kubernetes ConifgMap spec for the terraform module
// This configmap will be mounted in the Terraform Runner pod
func getConfigMapSpecForModule(name string, namespace string, module string, runID string, owner metav1.OwnerReference) *corev1.ConfigMap {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getUniqueResourceName(name, runID),
			Namespace: namespace,
			Labels:    getCommonLabels(name, runID),
			OwnerReferences: []metav1.OwnerReference{
				owner,
			},
		},
		Data: map[string]string{
			"main.tf": module,
		},
	}

	return cm
}

// createConfigMapForModule creates the ConfigMap for the Terraform workflow/run
func createConfigMapForModule(namespacedName types.NamespacedName, run *Terraform) (*corev1.ConfigMap, error) {
	configMaps := kube.ClientSet.CoreV1().ConfigMaps(namespacedName.Namespace)

	tpl, err := getTerraformModuleFromTemplate(run)

	if err != nil {
		return nil, err
	}

	configMap := getConfigMapSpecForModule(
		namespacedName.Name,
		namespacedName.Namespace,
		string(tpl), run.Status.RunID,
		run.GetOwnerReference())

	if _, err := configMaps.Create(context.TODO(), configMap, metav1.CreateOptions{}); err != nil {
		return nil, err
	}

	return configMap, nil
}
