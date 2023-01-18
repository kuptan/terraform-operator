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
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"

	"github.com/go-logr/logr"
	"github.com/kuptan/terraform-operator/api/v1alpha1"
	"github.com/kuptan/terraform-operator/internal/metrics"
)

// TerraformReconciler reconciles a Terraform object
type TerraformReconciler struct {
	client.Client
	Scheme            *runtime.Scheme
	Recorder          record.EventRecorder
	MetricsRecorder   metrics.RecorderInterface
	Log               logr.Logger
	requeueDependency time.Duration
	requeueJobWatch   time.Duration
}

// TerraformReconcilerOptions holds additional options
type TerraformReconcilerOptions struct {
	RequeueDependencyInterval time.Duration
	RequeueJobWatchInterval   time.Duration
}

//+kubebuilder:rbac:groups=run.terraform-operator.io,resources=terraforms,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=run.terraform-operator.io,resources=terraforms/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=run.terraform-operator.io,resources=terraforms/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Terraform object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *TerraformReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	run := &v1alpha1.Terraform{}
	start := time.Now()
	durationMsg := fmt.Sprintf("reconcilation finished in %s", time.Since(start).String())

	if err := r.Get(ctx, req.NamespacedName, run); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if !controllerutil.ContainsFinalizer(run, v1alpha1.TerraformFinalizer) {
		controllerutil.AddFinalizer(run, v1alpha1.TerraformFinalizer)
		if err := r.Update(ctx, run); err != nil {
			r.Log.Error(err, "unable to register finalizer")

			return ctrl.Result{}, err
		}

		r.Recorder.Event(run, corev1.EventTypeNormal, "Added finalizer", "Object finalizer is added")

		return ctrl.Result{}, nil
	}

	// Examine if the object is under deletion
	if !run.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.handleRunDelete(ctx, run)
	}

	if run.IsSubmitted() || run.IsWaiting() {
		result, err := r.handleRunCreate(ctx, run, req.NamespacedName)

		if err != nil {
			return ctrl.Result{}, err
		}

		r.Recorder.Event(run, "Normal", "Created", fmt.Sprintf("Run(%s) submitted", run.Status.RunID))
		r.MetricsRecorder.RecordTotal(run.Name, run.Namespace)

		if result.RequeueAfter > 0 {
			r.Log.Info(fmt.Sprintf("%s, next run in %s", durationMsg, result.RequeueAfter.String()))

			return result, nil
		}

		return result, nil
	}

	if run.IsStarted() {
		result, err := r.handleRunJobWatch(ctx, run)

		if err != nil {
			return ctrl.Result{}, err
		}

		if result.RequeueAfter > 0 {
			r.Log.Info(fmt.Sprintf("%s, next run in %s", durationMsg, result.RequeueAfter.String()))

			return result, nil
		}

		return result, nil
	}

	if run.IsUpdated() {
		r.Log.Info("updating a terraform run")

		result, err := r.handleRunUpdate(ctx, run, req.NamespacedName)

		if err != nil {
			return ctrl.Result{}, err
		}

		if result.RequeueAfter > 0 {
			r.Log.Info(fmt.Sprintf("%s, next run in %s", durationMsg, result.RequeueAfter.String()))

			return result, nil
		}

		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TerraformReconciler) SetupWithManager(mgr ctrl.Manager, opts TerraformReconcilerOptions) error {
	r.requeueDependency = opts.RequeueDependencyInterval
	r.requeueJobWatch = opts.RequeueJobWatchInterval

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Terraform{}).
		Owns(&batchv1.Job{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}
