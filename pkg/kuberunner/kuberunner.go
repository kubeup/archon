package kuberunner

import (
	"net"
	"os"
	"os/signal"

	"github.com/golang/glog"
	"k8s.io/apiserver/pkg/util/feature"
	"k8s.io/kubernetes/pkg/capabilities"
	"k8s.io/kubernetes/pkg/kubelet/types"
	"k8s.io/minikube/pkg/localkube"
	"k8s.io/minikube/pkg/util"
)

func NewLocalkubeServer() *localkube.LocalkubeServer {
	_, defaultServiceClusterIPRange, _ := net.ParseCIDR("10.0.0.1/24")

	return &localkube.LocalkubeServer{
		Containerized:            false,
		DNSDomain:                "cluster.local",
		DNSIP:                    net.ParseIP("10.0.0.10"),
		LocalkubeDirectory:       "./.localkube",
		ServiceClusterIPRange:    *defaultServiceClusterIPRange,
		APIServerAddress:         net.ParseIP("127.0.0.1"),
		APIServerPort:            443,
		APIServerInsecureAddress: net.ParseIP("127.0.0.1"),
		APIServerInsecurePort:    8080,
		ShouldGenerateCerts:      true,
		ShowVersion:              false,
		RuntimeConfig:            map[string]string{"api/all": "true"},
		ExtraConfig:              util.ExtraOptionSlice{},
	}
}

func StartLocalkubeServer(s *localkube.LocalkubeServer) {
	if s.ShouldGenerateCerts {
		if err := s.GenerateCerts(); err != nil {
			glog.Errorf("Failed to create certificates!")
			panic(err)
		}
	}

	//Set feature gates
	glog.Infof("Feature gates:", s.FeatureGates)
	if s.FeatureGates != "" {
		err := feature.DefaultFeatureGate.Set(s.FeatureGates)
		if err != nil {
			glog.Errorf("Error setting feature gates: %s", err)
		}
	}

	// Setup capabilities. This can only be done once per binary.
	allSources, _ := types.GetValidatedSources([]string{types.AllSource})
	c := capabilities.Capabilities{
		AllowPrivileged: true,
		PrivilegedSources: capabilities.PrivilegedSources{
			HostNetworkSources: allSources,
			HostIPCSources:     allSources,
			HostPIDSources:     allSources,
		},
	}
	capabilities.Initialize(c)

	// setup etcd
	etcd, err := s.NewEtcd(localkube.KubeEtcdClientURLs, localkube.KubeEtcdPeerURLs, "kubeetcd", s.GetEtcdDataDirectory())
	if err != nil {
		panic(err)
	}
	etcd.Start()
	defer etcd.Stop()

	// setup access to etcd
	netIP, _ := s.GetHostIP()
	glog.Infof("localkube host ip address: %s\n", netIP.String())

	// setup apiserver
	apiserver := s.NewAPIServer()
	s.AddServer(apiserver)

	// setup controller-manager
	//controllerManager := s.NewControllerManagerServer()
	//s.AddServer(controllerManager)

	s.StartAll()
	defer s.StopAll()

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)

	<-interruptChan
	glog.Infof("Shutting down localkube...")
}
