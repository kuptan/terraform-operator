package controllers

import (
	"context"
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kube-champ/terraform-operator/api/v1alpha1"
	"github.com/kube-champ/terraform-operator/pkg/kube"
)

func prepareRunnerRBAC() error {
	namespace := "default"
	serviceAccountName := "terraform-runner"

	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceAccountName,
			Namespace: namespace,
		},
	}

	clusterRole := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceAccountName,
		},
		Rules: []rbacv1.PolicyRule{
			rbacv1.PolicyRule{
				APIGroups: []string{""},
				Resources: []string{"secrets"},
				Verbs:     []string{"get", "update"},
			},
		},
	}

	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceAccountName,
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     serviceAccountName,
			APIGroup: "rbac.authorization.k8s.io",
		},
		Subjects: []rbacv1.Subject{
			rbacv1.Subject{
				Kind:      "ServiceAccount",
				Name:      serviceAccountName,
				Namespace: namespace,
			},
		},
	}

	ctx := context.Background()

	_, err := kube.ClientSet.CoreV1().ServiceAccounts("default").Create(ctx, serviceAccount, metav1.CreateOptions{})

	if err != nil {
		return err
	}

	_, err = kube.ClientSet.RbacV1().ClusterRoles().Create(ctx, clusterRole, metav1.CreateOptions{})

	if err != nil {
		return err
	}

	_, err = kube.ClientSet.RbacV1().ClusterRoleBindings().Create(ctx, clusterRoleBinding, metav1.CreateOptions{})

	if err != nil {
		return err
	}

	return nil
}

func makeRunJobRunning(r *v1alpha1.Terraform) {
	jobsClient := kube.ClientSet.BatchV1().Jobs("default")

	name := getRunName(r.Name, r.Status.RunId)

	job, _ := jobsClient.Get(context.Background(), name, metav1.GetOptions{})

	if job == nil {
		return
	}

	job.Status = batchv1.JobStatus{
		Active:    1,
		Succeeded: 0,
		Failed:    0,
	}

	jobsClient.Update(context.Background(), job, metav1.UpdateOptions{})

}

func makeRunJobSucceed(r *v1alpha1.Terraform) {
	jobsClient := kube.ClientSet.BatchV1().Jobs("default")

	name := getRunName(r.Name, r.Status.RunId)

	job, _ := jobsClient.Get(context.Background(), name, metav1.GetOptions{})

	if job == nil {
		return
	}

	job.Status = batchv1.JobStatus{
		Active:    0,
		Succeeded: 1,
		Failed:    0,
	}

	jobsClient.Update(context.Background(), job, metav1.UpdateOptions{})
}

func makeRunJobFail(r *v1alpha1.Terraform) {
	jobsClient := kube.ClientSet.BatchV1().Jobs("default")

	name := getRunName(r.Name, r.Status.RunId)

	job, _ := jobsClient.Get(context.Background(), name, metav1.GetOptions{})

	if job == nil {
		return
	}

	job.Status = batchv1.JobStatus{
		Active:    0,
		Succeeded: 0,
		Failed:    1,
	}

	jobsClient.Update(context.Background(), job, metav1.UpdateOptions{})

}

func isJobDeleted(r *v1alpha1.Terraform) bool {
	jobsClient := kube.ClientSet.BatchV1().Jobs("default")

	name := getRunName(r.Name, r.Status.RunId)

	_, err := jobsClient.Get(context.Background(), name, metav1.GetOptions{})

	return errors.IsNotFound(err)
}

func getRunName(name string, runId string) string {
	return fmt.Sprintf("%s-%s", name, runId)
}
