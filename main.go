/*
Copyright 2021.

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
	"context"
	"flag"
	"os"
	"time"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/getsentry/sentry-go"
	"github.com/imaware/fhir-operator/api"
	fhirv1alpha1 "github.com/imaware/fhir-operator/api/v1alpha1"
	"github.com/imaware/fhir-operator/controllers"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(fhirv1alpha1.AddToScheme(scheme))
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

	config, err := api.GetConfig()
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}
	opts := zap.Options{
		Development: config.DebugEnabled,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
	if config.SentryEnabled {
		err := sentry.Init(sentry.ClientOptions{
			// Either set your DSN here or set the SENTRY_DSN environment variable.
			Dsn: config.SentryDSN,
			// Either set environment and release here or set the SENTRY_ENVIRONMENT
			// and SENTRY_RELEASE environment variables.
			Environment: config.Environment,
			Release:     config.ReleaseTag,
			// Enable printing of SDK debug messages.
			// Useful when getting started or trying to figure something out.
			Debug: config.DebugEnabled,
			// Set TracesSampleRate to 1.0 to capture 100%
			// of transactions for performance monitoring.
			// We recommend adjusting this value in production,
			TracesSampleRate: config.SentrySampleRate,
		})
		if err != nil {
			setupLog.Error(err, "sentry.Init")
		}
		// Flush buffered events before the program terminates.
		// Set the timeout to the maximum duration the program can afford to wait.
		defer sentry.Flush(2 * time.Second)
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "212829d9.imaware.com",
		Namespace:              "",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}
	ctx := context.Background()
	pubsubClient, err := pubsub.NewClient(ctx, config.GCPProject)
	if err != nil {
		setupLog.Error(err, "unable to configure pubsub client")
		os.Exit(1)
	}
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		setupLog.Error(err, "unable to configure storage client")
		os.Exit(1)
	}
	if err = (&controllers.FhirStoreReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr, config); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "FhirStore")
		os.Exit(1)
	}
	if err = (&controllers.FhirResourceReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr, config); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "FhirResource")
		os.Exit(1)
	}
	apiClientPubsub := &api.GCPPUBClient{PubsubClient: pubsubClient}

	apiClientStorage := &api.GCSClient{StorageClient: storageClient}

	if err = (&controllers.FhirGCSConnectorReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		PubSubClient:  apiClientPubsub,
		StorageClient: apiClientStorage,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "FhirGCSConnector")
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
