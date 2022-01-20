package controllers

import (
	"context"
	"errors"
	"fmt"

	"github.com/kube-champ/terraform-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
)

func updateRunStatus(r *TerraformReconciler, run *v1alpha1.Terraform, status v1alpha1.TerraformRunStatus) {
	run.Status.RunStatus = status
	r.Status().Update(context.Background(), run)
}

func (r *TerraformReconciler) create(run *v1alpha1.Terraform, namespacedName types.NamespacedName) (bool, error) {
	isDependencyReady, err := run.DependenciesCompleted()

	run.Status.Generation = run.Generation

	if err != nil {
		r.Log.Error(err, "failed create a terraform run while trying to check dependencies")

		updateRunStatus(r, run, v1alpha1.RunFailed)

		return false, err
	}

	if !isDependencyReady {
		if !run.IsWaiting() {
			r.Recorder.Event(run, "Normal", "Waiting", "Run has dependencies that are not completed")
			updateRunStatus(r, run, v1alpha1.RunWaitingForDependency)
		}

		return true, nil
	}

	run.SetRunId()

	_, err = run.CreateTerraformRun(namespacedName)

	if err != nil {
		r.Log.Error(err, "failed create a terraform run ")

		updateRunStatus(r, run, v1alpha1.RunFailed)

		return false, err
	}

	updateRunStatus(r, run, v1alpha1.RunStarted)

	return false, nil
}

func (r *TerraformReconciler) update(run *v1alpha1.Terraform, namespacedName types.NamespacedName) (bool, error) {
	run.PrepareForUpdate()

	r.Recorder.Event(run, "Normal", "Updated", "Creating a new run job")

	return r.create(run, namespacedName)
}

func (r *TerraformReconciler) watchRun(run *v1alpha1.Terraform, namespacedName types.NamespacedName) (bool, error) {
	job, err := run.GetJobByRun()

	r.Log.Info("watching job run to complete", "name", job.Name)

	if err != nil {
		return false, err
	}

	// job hasn't started
	if job.Status.Active == 0 && job.Status.Succeeded == 0 && job.Status.Failed == 0 {
		return true, nil
	}

	// job is still running
	if job.Status.Active > 0 {
		if !run.IsRunning() {
			updateRunStatus(r, run, v1alpha1.RunRunning)

			r.Recorder.Event(run, "Normal", "Running", fmt.Sprintf("Run(%s) waiting for run job to finish", run.Status.RunId))
		}

		return true, nil
	}

	// job is successful
	if job.Status.Succeeded > 0 {
		r.Log.Info("terraform run job completed successfully")

		if run.Spec.DeleteCompletedJobs {
			r.Log.Info("deleting completed job")

			if err := run.DeleteAfterCompletion(); err != nil {
				r.Log.Error(err, "failed to delete run job after completion", "name", job.Name)
			} else {
				r.Recorder.Event(run, "Normal", "Cleanup", fmt.Sprintf("Run(%s) kubernetes job was deleted", run.Status.RunId))
			}
		}

		if !run.Spec.Destroy {
			r.Recorder.Event(run, "Normal", "Completed", fmt.Sprintf("Run(%s) completed", run.Status.RunId))
		} else {
			r.Recorder.Event(run, "Normal", "Destroyed", fmt.Sprintf("Run(%s) completed with terraform destroy", run.Status.RunId))
		}

		updateRunStatus(r, run, v1alpha1.RunCompleted)

		return false, nil
	}

	// if it got here, then the job is failed -- sadly .... :( :( :(
	r.Recorder.Event(run, "Warning", "Failed", fmt.Sprintf("Run(%s) failed", run.Status.RunId))
	r.Log.Error(errors.New("job failed"), "terraform run job failed to complete", "name", job.Name)

	updateRunStatus(r, run, v1alpha1.RunFailed)

	return false, nil
}
