package controllers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kube-champ/terraform-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	requeueJobWatch   time.Duration = 20 * time.Second
	requeueDependency time.Duration = 25 * time.Second
)

func updateRunStatus(r *TerraformReconciler, run *v1alpha1.Terraform, status v1alpha1.TerraformRunStatus) {
	run.Status.RunStatus = status
	r.Status().Update(context.Background(), run)
}

func (r *TerraformReconciler) create(run *v1alpha1.Terraform, namespacedName types.NamespacedName) (ctrl.Result, error) {
	err := r.checkDependencies(*run)

	run.Status.ObservedGeneration = run.Generation

	if err != nil {
		if !run.IsWaiting() {
			r.Recorder.Event(run, "Normal", "Waiting", "Dependencies are not yet completed")
			updateRunStatus(r, run, v1alpha1.RunWaitingForDependency)
		}

		return ctrl.Result{
			RequeueAfter: requeueDependency,
		}, nil
	}

	run.SetRunId()

	_, err = run.CreateTerraformRun(namespacedName)

	if err != nil {
		r.Log.Error(err, "failed create a terraform run")

		updateRunStatus(r, run, v1alpha1.RunFailed)

		return ctrl.Result{}, err
	}

	updateRunStatus(r, run, v1alpha1.RunStarted)

	return ctrl.Result{}, nil
}

func (r *TerraformReconciler) update(run *v1alpha1.Terraform, namespacedName types.NamespacedName) (ctrl.Result, error) {
	run.PrepareForUpdate()

	r.Recorder.Event(run, "Normal", "Updated", "Creating a new run job")

	return r.create(run, namespacedName)
}

func (r *TerraformReconciler) watchRun(run *v1alpha1.Terraform, namespacedName types.NamespacedName) (ctrl.Result, error) {
	job, err := run.GetJobByRun()

	r.Log.Info("watching job run to complete", "name", job.Name)

	if err != nil {
		return ctrl.Result{}, err
	}

	// job hasn't started
	if job.Status.Active == 0 && job.Status.Succeeded == 0 && job.Status.Failed == 0 {
		return ctrl.Result{RequeueAfter: requeueJobWatch}, nil
	}

	// job is still running
	if job.Status.Active > 0 {
		if !run.IsRunning() {
			updateRunStatus(r, run, v1alpha1.RunRunning)

			r.Recorder.Event(run, "Normal", "Running", fmt.Sprintf("Run(%s) waiting for run job to finish", run.Status.RunId))
		}

		return ctrl.Result{RequeueAfter: requeueJobWatch}, nil
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

		return ctrl.Result{}, nil
	}

	// if it got here, then the job is failed -- sadly .... :( :( :(
	r.Recorder.Event(run, "Warning", "Failed", fmt.Sprintf("Run(%s) failed", run.Status.RunId))
	r.Log.Error(errors.New("job failed"), "terraform run job failed to complete", "name", job.Name)

	updateRunStatus(r, run, v1alpha1.RunFailed)

	return ctrl.Result{}, nil
}

func (r *TerraformReconciler) checkDependencies(run v1alpha1.Terraform) error {
	for _, d := range run.Spec.DependsOn {

		if d.Namespace == "" {
			d.Namespace = run.Namespace
		}

		dName := types.NamespacedName{
			Namespace: d.Namespace,
			Name:      d.Name,
		}

		var dRun v1alpha1.Terraform

		err := r.Get(context.Background(), dName, &dRun)

		if err != nil {
			return fmt.Errorf("unable to get '%s' dependency: %w", dName, err)
		}

		if dRun.Status.RunStatus != v1alpha1.RunCompleted {
			return fmt.Errorf("dependency '%s' is not ready", dName)
		}
	}

	return nil
}
