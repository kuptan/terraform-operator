package controllers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kuptan/terraform-operator/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	requeueJobWatch   time.Duration = 20 * time.Second
	requeueDependency time.Duration = 25 * time.Second
)

func (r *TerraformReconciler) updateRunStatus(ctx context.Context, run *v1alpha1.Terraform, status v1alpha1.TerraformRunStatus) {
	run.Status.RunStatus = status

	if err := r.Status().Update(ctx, run); err != nil {
		r.Log.Error(err, "failed to update status")
	}
}

func (r *TerraformReconciler) handleRunCreate(ctx context.Context, run *v1alpha1.Terraform, namespacedName types.NamespacedName) (ctrl.Result, error) {
	dependencies, err := r.checkDependencies(ctx, *run)

	run.Status.ObservedGeneration = run.Generation

	if err != nil {
		if !run.IsWaiting() {
			r.Recorder.Event(run, "Normal", "Waiting", "Dependencies are not yet completed")
			r.updateRunStatus(ctx, run, v1alpha1.RunWaitingForDependency)
		}

		return ctrl.Result{
			RequeueAfter: requeueDependency,
		}, nil
	}

	run.SetRunID()

	setVariablesFromDependencies(run, dependencies)

	job, err := run.CreateTerraformRun(namespacedName)

	if err != nil {
		r.Log.Error(err, "failed create a terraform run")

		r.updateRunStatus(ctx, run, v1alpha1.RunFailed)

		return ctrl.Result{}, err
	}

	run.Status.OutputSecretName = job.ObjectMeta.Name
	r.updateRunStatus(ctx, run, v1alpha1.RunStarted)

	return ctrl.Result{}, nil
}

func (r *TerraformReconciler) handleRunUpdate(ctx context.Context, run *v1alpha1.Terraform, namespacedName types.NamespacedName) (ctrl.Result, error) {
	run.PrepareForUpdate()

	r.Recorder.Event(run, "Normal", "Updated", "Creating a new run job")

	return r.handleRunCreate(ctx, run, namespacedName)
}

func (r *TerraformReconciler) handleRunJobWatch(ctx context.Context, run *v1alpha1.Terraform) (ctrl.Result, error) {
	job, err := run.GetJobByRun()

	r.Log.Info("waiting for terraform job run to complete", "name", job.Name)

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
			r.updateRunStatus(ctx, run, v1alpha1.RunRunning)

			r.Recorder.Event(run, "Normal", "Running", fmt.Sprintf("Run(%s) waiting for run job to finish", run.Status.RunID))
		}

		return ctrl.Result{RequeueAfter: requeueJobWatch}, nil
	}

	// job is successful
	if job.Status.Succeeded > 0 {
		r.Log.Info("terraform run job completed successfully")

		if run.Spec.DeleteCompletedJobs {
			r.Log.Info("deleting completed job")

			if err := run.DeleteAfterCompletion(); err != nil {
				r.Log.Error(err, "failed to delete terraform run job after completion", "name", job.Name)
			} else {
				r.Recorder.Event(run, "Normal", "Cleanup", fmt.Sprintf("Run(%s) kubernetes job was deleted", run.Status.RunID))
			}
		}

		if !run.Spec.Destroy {
			r.Recorder.Event(run, "Normal", "Completed", fmt.Sprintf("Run(%s) completed", run.Status.RunID))
		} else {
			r.Recorder.Event(run, "Normal", "Destroyed", fmt.Sprintf("Run(%s) completed with terraform destroy", run.Status.RunID))
		}

		r.updateRunStatus(ctx, run, v1alpha1.RunCompleted)

		return ctrl.Result{}, nil
	}

	// if it got here, then the job is failed -- sadly .... :( :( :(
	r.Recorder.Event(run, "Warning", "Failed", fmt.Sprintf("Run(%s) failed", run.Status.RunID))
	r.Log.Error(errors.New("job failed"), "terraform run job failed to complete", "name", job.Name)

	r.updateRunStatus(ctx, run, v1alpha1.RunFailed)

	return ctrl.Result{}, nil
}

func (r *TerraformReconciler) checkDependencies(ctx context.Context, run v1alpha1.Terraform) ([]v1alpha1.Terraform, error) {
	dependencies := []v1alpha1.Terraform{}

	for _, d := range run.Spec.DependsOn {

		if d.Namespace == "" {
			d.Namespace = run.Namespace
		}

		dName := types.NamespacedName{
			Namespace: d.Namespace,
			Name:      d.Name,
		}

		var dRun v1alpha1.Terraform

		err := r.Get(ctx, dName, &dRun)

		if err != nil {
			return dependencies, fmt.Errorf("unable to get '%s' dependency: %w", dName, err)
		}

		if dRun.Status.RunStatus != v1alpha1.RunCompleted {
			return dependencies, fmt.Errorf("dependency '%s' is not ready", dName)
		}

		dependencies = append(dependencies, dRun)
	}

	return dependencies, nil
}

// setVariablesFromDependencies sets the variable from the output of a dependency
// this currently only works with runs within the same namespace
func setVariablesFromDependencies(run *v1alpha1.Terraform, dependencies []v1alpha1.Terraform) {
	if len(dependencies) == 0 {
		return
	}

	for _, v := range run.Spec.Variables {
		if v.DependencyRef == nil {
			continue
		}

		for index, d := range dependencies {
			if d.Name == v.DependencyRef.Name && d.Namespace == run.Namespace {
				tfVarRef := &v1.EnvVarSource{
					SecretKeyRef: &v1.SecretKeySelector{
						Key: v.DependencyRef.Key,
						LocalObjectReference: v1.LocalObjectReference{
							Name: d.Status.OutputSecretName,
						},
					},
				}

				tfVar := v1alpha1.Variable{
					Key:           v.Key,
					DependencyRef: v.DependencyRef,
					ValueFrom:     tfVarRef,
				}

				// remove the current variable from the list
				run.Spec.Variables = append(run.Spec.Variables[:index], run.Spec.Variables[index+1:]...)
				// add a new variable with the valueFrom
				run.Spec.Variables = append(run.Spec.Variables, tfVar)
			}
		}
	}
}
