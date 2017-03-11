/*
Copyright 2014 The Kubernetes Authors.

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

// Package app implements a server that runs a set of active
// components.  This includes replication controllers, service endpoints and
// nodes.
//
// CAUTION: If you update code in this file, you may need to also update code
//          in contrib/mesos/pkg/controllermanager/controllermanager.go
package app

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"strconv"
	"time"

	"kubeup.com/archon/cmd/archon-controller/app/options"
	archonclientset "kubeup.com/archon/pkg/clientset"
	"kubeup.com/archon/pkg/cloudprovider"
	_ "kubeup.com/archon/pkg/cloudprovider/providers"
	"kubeup.com/archon/pkg/controller/certificate"
	"kubeup.com/archon/pkg/controller/instance"
	"kubeup.com/archon/pkg/controller/instancegroup"
	"kubeup.com/archon/pkg/controller/network"
	"kubeup.com/archon/pkg/kuberunner"

	clientv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/kubernetes/pkg/api"
	//clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/record"
	"k8s.io/kubernetes/pkg/client/leaderelection"
	"k8s.io/kubernetes/pkg/client/leaderelection/resourcelock"
	//"k8s.io/kubernetes/pkg/controller/informers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apiserver/pkg/server/healthz"
	"k8s.io/kubernetes/pkg/util/configz"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	// Jitter used when starting controller managers
	ControllerStartJitter = 1.0
)

// NewControllerManagerCommand creates a *cobra.Command object with default parameters
func NewControllerManagerCommand() *cobra.Command {
	s := options.NewCMServer()
	s.AddFlags(pflag.CommandLine)
	cmd := &cobra.Command{
		Use: "kube-controller-manager",
		Long: `The Kubernetes controller manager is a daemon that embeds
the core control loops shipped with Kubernetes. In applications of robotics and
automation, a control loop is a non-terminating loop that regulates the state of
the system. In Kubernetes, a controller is a control loop that watches the shared
state of the cluster through the apiserver and makes changes attempting to move the
current state towards the desired state. Examples of controllers that ship with
Kubernetes today are the replication controller, endpoints controller, namespace
controller, and serviceaccounts controller.`,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	return cmd
}

func ResyncPeriod(s *options.CMServer) func() time.Duration {
	return func() time.Duration {
		factor := rand.Float64() + 1
		return time.Duration(float64(s.MinResyncPeriod.Nanoseconds()) * factor)
	}
}

// Run runs the CMServer.  This should never exit.
func Run(s *options.CMServer) error {
	if s.TestRun {
		os.Exit(0)
	}

	if c, err := configz.New("componentconfig"); err == nil {
		c.Set(s.ArchonControllerManagerConfiguration)
	} else {
		glog.Errorf("unable to register configz: %s", err)
	}

	if s.EnableLocalkube {
		lk := kuberunner.NewLocalkubeServer()
		lk.LocalkubeDirectory = s.LocalkubeDirectory
		lk.APIServerAddress = s.APIServerAddress
		lk.APIServerPort = s.APIServerPort
		lk.APIServerInsecureAddress = s.APIServerInsecureAddress
		lk.APIServerInsecurePort = s.APIServerInsecurePort

		go kuberunner.StartLocalkubeServer(lk)

		s.Master = fmt.Sprintf("http://127.0.0.1:%d", lk.APIServerInsecurePort)
		if s.LeaderElection.LeaderElect == true {
			s.LeaderElection.LeaderElect = false
			glog.Warningf("Leader election is forced off when localkube is used")
		}
	}

	kubeconfig, err := clientcmd.BuildConfigFromFlags(s.Master, s.Kubeconfig)
	if err != nil {
		return err
	}

	kubeconfig.ContentConfig.ContentType = s.ContentType
	// Override kubeconfig qps/burst settings from flags
	kubeconfig.QPS = s.KubeAPIQPS
	kubeconfig.Burst = int(s.KubeAPIBurst)

	kubeClient, err := archonclientset.NewForConfig(kubeconfig)
	if err != nil {
		glog.Fatalf("Invalid API configuration: %v", err)
	}

	kubeClient.EnsureResources(30 * time.Second)

	go func() {
		mux := http.NewServeMux()
		healthz.InstallHandler(mux)
		if s.EnableProfiling {
			mux.HandleFunc("/debug/pprof/", pprof.Index)
			mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
			mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		}
		configz.InstallHandler(mux)
		mux.Handle("/metrics", prometheus.Handler())

		server := &http.Server{
			Addr:    net.JoinHostPort(s.Address, strconv.Itoa(int(s.Port))),
			Handler: mux,
		}
		glog.Fatal(server.ListenAndServe())
	}()

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: v1core.New(kubeClient.Core().RESTClient()).Events("")})
	recorder := eventBroadcaster.NewRecorder(api.Scheme, clientv1.EventSource{Component: "archon-controller"})

	run := func(stop <-chan struct{}) {
		err := StartControllers(s, kubeClient, kubeconfig, stop, recorder)
		glog.Fatalf("error running controllers: %v", err)
		panic("unreachable")
	}

	if !s.LeaderElection.LeaderElect {
		run(nil)
		panic("unreachable")
	}

	id, err := os.Hostname()
	if err != nil {
		return err
	}

	// TODO: enable other lock types
	rl := resourcelock.EndpointsLock{
		EndpointsMeta: metav1.ObjectMeta{
			Namespace: "kube-system",
			Name:      s.ControllerName,
		},
		Client: kubeClient,
		LockConfig: resourcelock.ResourceLockConfig{
			Identity:      id,
			EventRecorder: recorder,
		},
	}

	leaderelection.RunOrDie(leaderelection.LeaderElectionConfig{
		Lock:          &rl,
		LeaseDuration: s.LeaderElection.LeaseDuration.Duration,
		RenewDeadline: s.LeaderElection.RenewDeadline.Duration,
		RetryPeriod:   s.LeaderElection.RetryPeriod.Duration,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: run,
			OnStoppedLeading: func() {
				glog.Fatalf("leaderelection lost")
			},
		},
	})
	panic("unreachable")
}

func StartControllers(s *options.CMServer, kubeClient archonclientset.Interface, kubeconfig *restclient.Config, stop <-chan struct{}, recorder record.EventRecorder) error {
	// TODO: make use of sharedInformers
	//sharedInformers := informers.NewSharedInformerFactory(clientset.NewForConfigOrDie(restclient.AddUserAgent(kubeconfig, "shared-informers")), ResyncPeriod(s)())

	glog.Infof("Starting controllers in namespace: %s", s.Namespace)

	cloud, err := cloudprovider.InitCloudProvider(s.CloudProvider, s.CloudConfigFile)
	if err != nil {
		glog.Fatalf("Cloud provider could not be initialized: %v", err)
	}

	glog.Infof("Starting Instance controller")
	instanceController, err := instance.New(cloud, kubeClient, s.ClusterName, s.Namespace)
	if err != nil {
		glog.Errorf("Failed to start instance controller: %v", err)
	} else {
		instanceController.Run(5)
	}

	glog.Infof("Starting InstanceGroup controller")
	go instancegroup.NewInstanceGroupController(kubeClient, s.Namespace, 1, 10, false).Run(3, wait.NeverStop)
	time.Sleep(wait.Jitter(s.ControllerStartInterval.Duration, ControllerStartJitter))

	glog.Infof("Starting Network controller")
	networkController, err := network.New(cloud, kubeClient, s.ClusterName, s.Namespace)
	if err != nil {
		glog.Errorf("Failed to start network controller: %v", err)
	} else {
		networkController.Run(3)
	}

	glog.Infof("Starting Certificate controller")
	certificateController, err := certificate.New(kubeClient, 30*time.Second, s.ClusterSigningCertFile, s.ClusterSigningKeyFile, s.Namespace)
	if err != nil {
		glog.Errorf("Failed to start certificate controller: %v", err)
	} else {
		certificateController.Run(3, wait.NeverStop)
	}

	//sharedInformers.Start(stop)

	select {}
}
