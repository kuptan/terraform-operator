package v1alpha1

import (
	"context"

	"github.com/kube-champ/terraform-operator/internal/kube"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// creates a k8s ConifgMap for the terraform module string
// This configmap will be mounted in the pod so it can run the Terraform

func getConfigMapSpecForModule(name string, namespace string, module string, runId string, owner metav1.OwnerReference) *corev1.ConfigMap {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getUniqueResourceName(name, runId),
			Namespace: namespace,
			Labels:    getCommonLabels(name, runId),
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

func createConfigMapForModule(namespacedName types.NamespacedName, run *Terraform) (*corev1.ConfigMap, error) {
	configMaps := kube.ClientSet.CoreV1().ConfigMaps(namespacedName.Namespace)

	tpl, err := getTerraformModuleFromTemplate(run)

	if err != nil {
		return nil, err
	}

	configMap := getConfigMapSpecForModule(
		namespacedName.Name,
		namespacedName.Namespace,
		string(tpl), run.Status.RunId,
		run.GetOwnerReference())

	if _, err := configMaps.Create(context.TODO(), configMap, metav1.CreateOptions{}); err != nil {
		return nil, err
	}

	return configMap, nil
}
