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

// Package options provides the flags used for the controller manager.
//
// CAUTION: If you update code in this file, you may need to also update code
//          in contrib/mesos/pkg/controllermanager/controllermanager.go
package options

import (
	"net"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/kubernetes/pkg/apis/componentconfig"
	"k8s.io/kubernetes/pkg/client/leaderelection"

	"github.com/spf13/pflag"
)

// Used in leader election
const ArchonControllerName = "archon-controller"

type ArchonControllerManagerConfiguration struct {
	metav1.TypeMeta

	// port is the port that the controller-manager's http service runs on.
	Port int32 `json:"port"`
	// address is the IP address to serve on (set to 0.0.0.0 for all interfaces).
	Address string `json:"address"`
	// cloudProvider is the provider for cloud services.
	CloudProvider string `json:"cloudProvider"`
	// cloudConfigFile is the path to the cloud provider configuration file.
	CloudConfigFile string `json:"cloudConfigFile"`
	// minResyncPeriod is the resync period in reflectors; will be random between
	// minResyncPeriod and 2*minResyncPeriod.
	MinResyncPeriod        metav1.Duration `json:"minResyncPeriod"`
	ClusterSigningCertFile string          `json:"clusterSigningCertFile"`
	// clusterSigningCertFile is the filename containing a PEM-encoded
	// RSA or ECDSA private key used to issue cluster-scoped certificates
	ClusterSigningKeyFile string `json:"clusterSigningKeyFile"`
	// clusterName is the instance prefix for the cluster.
	ClusterName string `json:"clusterName"`
	// leaderElection defines the configuration of leader election client.
	LeaderElection componentconfig.LeaderElectionConfiguration `json:"leaderElection"`
	// enableProfiling enables profiling via web interface host:port/debug/pprof/
	EnableProfiling bool `json:"enableProfiling"`
	// contentType is contentType of requests sent to apiserver.
	ContentType string `json:"contentType"`
	// kubeAPIQPS is the QPS to use while talking with kubernetes apiserver.
	KubeAPIQPS float32 `json:"kubeAPIQPS"`
	// kubeAPIBurst is the burst to use while talking with kubernetes apiserver.
	KubeAPIBurst int32 `json:"kubeAPIBurst"`
	// How long to wait between starting controller managers
	ControllerStartInterval metav1.Duration `json:"controllerStartInterval"`
	// enables the generic garbage collector. MUST be synced with the
	// corresponding flag of the kube-apiserver. WARNING: the generic garbage
	// collector is an alpha feature.
	EnableGarbageCollector bool `json:"enableGarbageCollector"`

	TestRun bool

	// Localkube stuff
	EnableLocalkube          bool   `json:"enableLocalkube"`
	APIServerAddress         net.IP `json:"apiServerAddress"`
	APIServerPort            int    `json:"apiServerPort"`
	APIServerInsecureAddress net.IP `json:"apiServerInsecureAddress"`
	APIServerInsecurePort    int    `json:"apiServerInsecurePort"`
	LocalkubeDirectory       string `json:"localkubeDirectory"`
}

// CMServer is the main context object for the controller manager.
type CMServer struct {
	ArchonControllerManagerConfiguration

	Master         string
	Kubeconfig     string
	Namespace      string
	ControllerName string
}

// NewCMServer creates a new CMServer with a default config.
func NewCMServer() *CMServer {
	s := CMServer{
		ArchonControllerManagerConfiguration: ArchonControllerManagerConfiguration{
			Port:                     12312,
			Address:                  "0.0.0.0",
			MinResyncPeriod:          metav1.Duration{Duration: 12 * time.Hour},
			ClusterName:              "kubernetes",
			ContentType:              "application/vnd.kubernetes.protobuf",
			KubeAPIQPS:               20.0,
			KubeAPIBurst:             30,
			LeaderElection:           leaderelection.DefaultLeaderElectionConfiguration(),
			ControllerStartInterval:  metav1.Duration{Duration: 0 * time.Second},
			EnableGarbageCollector:   true,
			ClusterSigningCertFile:   "/etc/kubernetes/ca/ca.pem",
			ClusterSigningKeyFile:    "/etc/kubernetes/ca/ca.key",
			EnableLocalkube:          false,
			LocalkubeDirectory:       "./.localkube",
			APIServerAddress:         net.ParseIP("127.0.0.1"),
			APIServerPort:            8443,
			APIServerInsecureAddress: net.ParseIP("127.0.0.1"),
			APIServerInsecurePort:    8080,
			TestRun:                  false,
		},
		Namespace:      "default",
		ControllerName: ArchonControllerName,
	}
	s.LeaderElection.LeaderElect = true
	return &s
}

// AddFlags adds flags for a specific CMServer to the specified FlagSet
func (s *CMServer) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&s.EnableLocalkube, "local", false, "Enable localkube")
	fs.StringVar(&s.LocalkubeDirectory, "localkube-directory", s.LocalkubeDirectory, "The directory localkube will store files in")
	fs.IPVar(&s.APIServerAddress, "apiserver-address", s.APIServerAddress, "The address the localkube apiserver will listen securely on")
	fs.IntVar(&s.APIServerPort, "apiserver-port", s.APIServerPort, "The port the localkube apiserver will listen securely on")
	fs.IPVar(&s.APIServerInsecureAddress, "apiserver-insecure-address", s.APIServerInsecureAddress, "The address the localkube apiserver will listen insecurely on")
	fs.IntVar(&s.APIServerInsecurePort, "apiserver-insecure-port", s.APIServerInsecurePort, "The port the localkube apiserver will listen insecurely on")
	fs.Int32Var(&s.Port, "port", s.Port, "The port that the controller-manager's http service runs on")
	fs.Var(componentconfig.IPVar{Val: &s.Address}, "address", "The IP address to serve on (set to 0.0.0.0 for all interfaces)")
	fs.StringVar(&s.CloudProvider, "cloud-provider", s.CloudProvider, "The provider for cloud services.  Empty string for no provider.")
	fs.StringVar(&s.CloudConfigFile, "cloud-config", s.CloudConfigFile, "The path to the cloud provider configuration file.  Empty string for no configuration file.")
	fs.StringVar(&s.Namespace, "namespace", s.Namespace, "The namespace that the controller will work in. Empty for default namespace")
	fs.DurationVar(&s.MinResyncPeriod.Duration, "min-resync-period", s.MinResyncPeriod.Duration, "The resync period in reflectors will be random between MinResyncPeriod and 2*MinResyncPeriod")
	fs.StringVar(&s.ClusterSigningCertFile, "cluster-signing-cert-file", s.ClusterSigningCertFile, "Filename containing a PEM-encoded X509 CA certificate used to issue cluster-scoped certificates")
	fs.StringVar(&s.ClusterSigningKeyFile, "cluster-signing-key-file", s.ClusterSigningKeyFile, "Filename containing a PEM-encoded RSA or ECDSA private key used to sign cluster-scoped certificates")
	fs.BoolVar(&s.EnableProfiling, "profiling", true, "Enable profiling via web interface host:port/debug/pprof/")
	fs.StringVar(&s.ClusterName, "cluster-name", s.ClusterName, "The instance prefix for the cluster")
	fs.StringVar(&s.ControllerName, "controller-name", s.ControllerName, "The name of the controller")
	fs.StringVar(&s.Master, "master", s.Master, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	fs.StringVar(&s.Kubeconfig, "kubeconfig", s.Kubeconfig, "Path to kubeconfig file with authorization and master location information.")
	fs.StringVar(&s.ContentType, "kube-api-content-type", s.ContentType, "Content type of requests sent to apiserver.")
	fs.Float32Var(&s.KubeAPIQPS, "kube-api-qps", s.KubeAPIQPS, "QPS to use while talking with kubernetes apiserver")
	fs.Int32Var(&s.KubeAPIBurst, "kube-api-burst", s.KubeAPIBurst, "Burst to use while talking with kubernetes apiserver")
	fs.DurationVar(&s.ControllerStartInterval.Duration, "controller-start-interval", s.ControllerStartInterval.Duration, "Interval between starting controller managers.")
	fs.BoolVar(&s.EnableGarbageCollector, "enable-garbage-collector", s.EnableGarbageCollector, "Enables the generic garbage collector. MUST be synced with the corresponding flag of the kube-apiserver.")
	fs.BoolVar(&s.TestRun, "test-run", s.TestRun, "Test if the binary is working correctly")

	leaderelection.BindFlags(&s.LeaderElection, fs)
	utilfeature.DefaultFeatureGate.AddFlag(fs)
}
