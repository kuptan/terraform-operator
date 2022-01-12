package controllers

import (
	"context"

	"github.com/kube-champ/terraform-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
)

func (r *TerraformReconciler) create(run *v1alpha1.Terraform, namespacedName types.NamespacedName) error {
	run.SetRunId()
	run.Status.Generation = run.Generation

	_, err := run.CreateTerraformRun(namespacedName)

	if err != nil {
		r.Log.Error(err, "failed create a terraform run")

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
