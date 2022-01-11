package controllers

import (
	"context"
	"fmt"

	"github.com/kube-champ/terraform-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *TerraformReconciler) watchRun(run *v1alpha1.Terraform, namespacedName types.NamespacedName) (bool, error) {
	l := log.FromContext(context.Background())

	job, err := run.GetJobByRun()

	l.Info("watching job", "name", job.Name)

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
			run.Status.RunStatus = v1alpha1.RunRunning

			r.Status().Update(context.Background(), run)
			r.Recorder.Event(run, "Normal", "Running", fmt.Sprintf("Run(%s) waiting for run job to finish", run.Status.RunId))
		}

		return true, nil
	}

	// job is successful
	if job.Status.Succeeded > 0 {
		if run.Spec.DeleteCompletedJobs {
			if err := run.DeleteAfterCompletion(); err != nil {
				l.Error(err, "failed to delete run job after completion")
			} else {
				r.Recorder.Event(run, "Normal", "Cleanup", fmt.Sprintf("Run(%s) kubernetes job was deleted", run.Status.RunId))
			}
		}

		if !run.Spec.Destroy {
			if err := r.updateOutputStatus(run, namespacedName); err != nil {
				r.Recorder.Event(run, "Normal", "Warn", "Failed to read run outputs from the generated run secret")
			}

			r.Recorder.Event(run, "Normal", "Completed", fmt.Sprintf("Run(%s) completed", run.Status.RunId))
		} else {
			r.Recorder.Event(run, "Normal", "Destroyed", fmt.Sprintf("Run(%s) completed with terraform destroy", run.Status.RunId))
		}

		run.Status.RunStatus = v1alpha1.RunCompleted
		r.Status().Update(context.Background(), run)

		return false, nil
	}

	// if it got here, then the job is failed -- sadly .... :( :( :(
	run.Status.RunStatus = v1alpha1.RunFailed
	r.Recorder.Event(run, "Error", "Failed", fmt.Sprintf("Run(%s) failed", run.Status.RunId))

	r.Status().Update(context.Background(), run)

	return false, nil
}
