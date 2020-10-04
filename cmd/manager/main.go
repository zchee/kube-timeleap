// Copyright 2020 The kube-timeleap Authors.
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"flag"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp" // for gcp auth provider
	ctrlconfig "sigs.k8s.io/controller-runtime/pkg/client/config"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	zapf "sigs.k8s.io/controller-runtime/pkg/log/zap"
	ctrlmanager "sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	injectorv1alpha1 "github.com/zchee/kube-timeleap/apis/injector/v1alpha1"
	timeleapv1alpha1 "github.com/zchee/kube-timeleap/apis/timeleap/v1alpha1"
	timeleapcontrollers "github.com/zchee/kube-timeleap/controllers/timeleap"
	"github.com/zchee/kube-timeleap/pkg/signalctx"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = logf.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(timeleapv1alpha1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

var (
	flagMetricsAddr          string
	flagEnableLeaderElection bool
)

const (
	flagMetricsAddrName  = "metrics-addr"
	flagMetricsAddrUsage = "The address the metric endpoint binds to."

	flagEnableLeaderElectionName = "enable-leader-election"
	flagEnableLeaderElectioUsage = "Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager."
)

func main() {
	flag.StringVar(&flagMetricsAddr, flagMetricsAddrName, ":8080", flagMetricsAddrUsage)
	flag.BoolVar(&flagEnableLeaderElection, flagEnableLeaderElectionName, false, flagEnableLeaderElectioUsage)
	flag.Parse()

	logf.SetLogger(zapf.New(zapf.UseDevMode(true)))

	mgr, err := ctrlmanager.New(ctrlconfig.GetConfigOrDie(), ctrlmanager.Options{
		Scheme:             scheme,
		MetricsBindAddress: flagMetricsAddr,
		Port:               9443,
		LeaderElection:     flagEnableLeaderElection,
		LeaderElectionID:   "timeleap.x-k8s.io",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&timeleapcontrollers.TimeLeapReconciler{
		Client: mgr.GetClient(),
		Log:    logf.Log.WithName("controllers").WithName("timeleap").WithName("TimeLeap"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "TimeLeap")
		os.Exit(1)
	}
	if err = (&timeleapv1alpha1.TimeLeap{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "TimeLeap")
		os.Exit(1)
	}
	podInjectorWebhook := &admission.Webhook{
		Handler: &injectorv1alpha1.Pod{
			Client: mgr.GetClient(),
		},
	}
	podInjectorWebhook.InjectLogger(logf.Log.WithName("injector").WithName("Pod"))
	webhookServer := mgr.GetWebhookServer()
	webhookServer.Register(injectorv1alpha1.Path, podInjectorWebhook)
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")

	ctx := signalctx.NewContext()
	if err := mgr.Start(ctx); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
