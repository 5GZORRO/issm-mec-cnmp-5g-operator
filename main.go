/*
Copyright 2021 IBM.

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
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	fivegv1alpha1 "github.com/5GZORRO/issm-mec-cnmp-5g-operator/api/v1alpha1"
	"github.com/5GZORRO/issm-mec-cnmp-5g-operator/controllers"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(fivegv1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	k8scfg := ctrl.GetConfigOrDie()
	mgr, err := ctrl.NewManager(k8scfg, ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "9c837c66.ibm.com",
	})
	if err != nil {
		setupLog.Error(err, "unable to transition manager")
		os.Exit(1)
	}

	if _, err = controllers.NewAmfReconciler(k8scfg, mgr, ctrl.Log.WithName("controllers").WithName("Smf"),
		mgr.GetScheme()); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Smf")
		os.Exit(1)
	}

	if _, err = controllers.NewSmfReconciler(k8scfg, mgr, ctrl.Log.WithName("controllers").WithName("Smf"),
		mgr.GetScheme()); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Smf")
		os.Exit(1)
	}

	if _, err = controllers.NewNrfReconciler(k8scfg, mgr, ctrl.Log.WithName("controllers").WithName("Nrf"),
		mgr.GetScheme()); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Nrf")
		os.Exit(1)
	}

	if _, err = controllers.NewMongoReconciler(k8scfg, mgr, ctrl.Log.WithName("controllers").WithName("Mongo"),
		mgr.GetScheme()); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Mongo")
		os.Exit(1)
	}

	if _, err = controllers.NewWebconsoleReconciler(k8scfg, mgr, ctrl.Log.WithName("controllers").WithName("Webconsole"),
		mgr.GetScheme()); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Webconsole")
		os.Exit(1)
	}

	if _, err = controllers.NewAusfReconciler(k8scfg, mgr, ctrl.Log.WithName("controllers").WithName("Ausf"),
		mgr.GetScheme()); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Ausf")
		os.Exit(1)
	}

	if _, err = controllers.NewUdmReconciler(k8scfg, mgr, ctrl.Log.WithName("controllers").WithName("Udm"),
		mgr.GetScheme()); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Udm")
		os.Exit(1)
	}

	if _, err = controllers.NewUdrReconciler(k8scfg, mgr, ctrl.Log.WithName("controllers").WithName("Udr"),
		mgr.GetScheme()); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Udr")
		os.Exit(1)
	}

	if _, err = controllers.NewUpfReconciler(k8scfg, mgr, ctrl.Log.WithName("controllers").WithName("Upf"),
		mgr.GetScheme()); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Upf")
		os.Exit(1)
	}

	if _, err = controllers.NewNssfReconciler(k8scfg, mgr, ctrl.Log.WithName("controllers").WithName("Nssf"),
		mgr.GetScheme()); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Nssf")
		os.Exit(1)
	}

	if _, err = controllers.NewPcfReconciler(k8scfg, mgr, ctrl.Log.WithName("controllers").WithName("Pcf"),
		mgr.GetScheme()); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Pcf")
		os.Exit(1)
	}

	if _, err = controllers.NewN3iwfReconciler(k8scfg, mgr, ctrl.Log.WithName("controllers").WithName("N3iwf"),
		mgr.GetScheme()); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "N3iwf")
		os.Exit(1)
	}

	if _, err = controllers.NewTransitionReconciler(k8scfg, mgr, ctrl.Log.WithName("controllers").WithName("Transition"),
		mgr.GetScheme()); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Transition")
		os.Exit(1)
	}

	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
