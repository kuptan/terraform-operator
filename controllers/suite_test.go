/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/kuptan/terraform-operator/api/v1alpha1"
	"github.com/kuptan/terraform-operator/internal/kube"
	"github.com/kuptan/terraform-operator/internal/utils"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var k8sClient client.Client
var testEnv *envtest.Environment

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	os.Setenv("DOCKER_REGISTRY", "docker.io")
	os.Setenv("TERRAFORM_RUNNER_IMAGE", "ibraheemalsaady/terraform-runner")
	os.Setenv("TERRAFORM_RUNNER_IMAGE_TAG", "0.0.3")
	os.Setenv("KNOWN_HOSTS_CONFIGMAP_NAME", "operator-known-hosts")

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
	}

	cfg, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = v1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})

	Expect(err).ToNot(HaveOccurred())

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	err = (&TerraformReconciler{
		Client:   k8sClient,
		Recorder: k8sManager.GetEventRecorderFor("terraform-controller"),
		Log:      ctrl.Log.WithName("controllers").WithName("TerraformController"),
	}).SetupWithManager(k8sManager)

	Expect(err).NotTo(HaveOccurred(), "failed to setup controller in test")

	kube.ClientSet = fake.NewSimpleClientset()
	utils.LoadEnv()

	err = prepareRunnerRBAC()
	Expect(err).ToNot(HaveOccurred(), "could not prepare rbac for terraform runner")

	go func() {
		defer GinkgoRecover()

		err = k8sManager.Start(ctrl.SetupSignalHandler())
		Expect(err).ToNot(HaveOccurred(), "failed to start manager")

		gexec.KillAndWait(4 * time.Second)

		err := testEnv.Stop()
		Expect(err).ToNot(HaveOccurred())
	}()

}, 60)

// var _ = AfterSuite(func() {
// 	By("tearing down the test environment")
// 	err := testEnv.Stop()
// 	Expect(err).NotTo(HaveOccurred())
// })

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

	name := getRunName(r.Name, r.Status.RunID)

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

	name := getRunName(r.Name, r.Status.RunID)

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

	name := getRunName(r.Name, r.Status.RunID)

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

	name := getRunName(r.Name, r.Status.RunID)

	_, err := jobsClient.Get(context.Background(), name, metav1.GetOptions{})

	return errors.IsNotFound(err)
}

func getRunName(name string, runID string) string {
	return fmt.Sprintf("%s-%s", name, runID)
}
