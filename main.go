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

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	crtlmetrics "sigs.k8s.io/controller-runtime/pkg/metrics"

	"github.com/kuptan/terraform-operator/api/v1alpha1"
	"github.com/kuptan/terraform-operator/controllers"
	"github.com/kuptan/terraform-operator/internal/kube"
	"github.com/kuptan/terraform-operator/internal/metrics"
	"github.com/kuptan/terraform-operator/internal/utils"
	//+kubebuilder:scaffold:imports
)

var (
	scheme            = runtime.NewScheme()
	setupLog          = ctrl.Log.WithName("setup")
	requeueDependency time.Duration
	requeueJobWatch   time.Duration
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(v1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.DurationVar(&requeueJobWatch, "requeue-job-watch", 10*time.Second, "The interval at which job status is reevaluated after a workflow is submitted.")
	flag.DurationVar(&requeueDependency, "requeue-dependency", 20*time.Second, "The interval at which dependencies are reevaluated.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: false,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	metricsRecorder := metrics.NewRecorder()
	crtlmetrics.Registry.MustRegister(metricsRecorder.Collectors()...)

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "d5cf1615.terraform-operator.io",
	})

	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	setupLog.Info(fmt.Sprintf("requeue dependency interval: %s", requeueDependency))
	setupLog.Info(fmt.Sprintf("requeue job watch interval: %s", requeueJobWatch))

	if err = (&controllers.TerraformReconciler{
		Client:          mgr.GetClient(),
		Scheme:          mgr.GetScheme(),
		Recorder:        mgr.GetEventRecorderFor("terraform-controller"),
		MetricsRecorder: metricsRecorder,
		Log:             ctrl.Log.WithName("controllers").WithName("TerraformController"),
	}).SetupWithManager(mgr, controllers.TerraformReconcilerOptions{
		RequeueDependencyInterval: requeueDependency,
		RequeueJobWatchInterval:   requeueJobWatch,
	}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Terraform")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err = mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}

	if err = mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	// Registering clientset
	_, err = kube.CreateK8SConfig()

	if err != nil {
		setupLog.Error(err, "could not create Kubernetes REST config")
		os.Exit(1)
	}

	utils.LoadEnv()

	setupLog.Info("starting manager")

	if err = mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
