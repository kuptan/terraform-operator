package controllers

import (
	"context"

	"github.com/kube-champ/terraform-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *TerraformReconciler) create(run *v1alpha1.Terraform, namespacedName types.NamespacedName) error {
	l := log.FromContext(context.Background())

	run.SetRunId()
	run.Status.Generation = run.Generation

	_, err := run.CreateTerraformRun(namespacedName)

	if err != nil {
		l.Error(err, "failed create a terraform run")

		run.Status.RunStatus = v1alpha1.RunFailed

		r.Status().Update(context.Background(), run)

		return err
	}

	run.Status.RunStatus = v1alpha1.RunStarted

	r.Status().Update(context.Background(), run)

	return nil
}

func (r *TerraformReconciler) update(run *v1alpha1.Terraform, namespacedName types.NamespacedName) error {
	run.PrepareForUpdate()

	r.Recorder.Event(run, "Normal", "Updated", "Creating a new run job")

	return r.create(run, namespacedName)
}
