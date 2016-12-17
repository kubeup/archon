package clientset

import (
	"kubeup.com/archon/pkg/clientset/archon"

	"github.com/golang/glog"
	kubernetes "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/client/restclient"

	_ "kubeup.com/archon/pkg/clientset/archon/install"
)

type Interface interface {
	kubernetes.Interface
	Archon() archon.ArchonInterface
}

type Clientset struct {
	kubernetes.Clientset
	*archon.ArchonClient
}

// Archon retrieves the ArchonClient
func (c *Clientset) Archon() archon.ArchonInterface {
	if c == nil {
		return nil
	}
	return c.ArchonClient
}

// NewForConfig creates a new Clientset for the given config.
func NewForConfig(c *restclient.Config) (*Clientset, error) {
	var clientset Clientset
	var err error
	cs, err := kubernetes.NewForConfig(c)
	if err != nil {
		return nil, err
	}
	clientset.Clientset = *cs
	clientset.ArchonClient, err = archon.NewForConfig(c)
	if err != nil {
		glog.Errorf("failed to create the ArchonClient: %v", err)
		return nil, err
	}
	return &clientset, nil
}

// NewForConfigOrDie creates a new Clientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *restclient.Config) *Clientset {
	var clientset Clientset
	clientset.Clientset = *kubernetes.NewForConfigOrDie(c)
	clientset.ArchonClient = archon.NewForConfigOrDie(c)
	return &clientset
}

// New creates a new Clientset for the given RESTClient.
func New(c *restclient.RESTClient) *Clientset {
	var clientset Clientset
	clientset.Clientset = *kubernetes.New(c)
	clientset.ArchonClient = archon.New(c)
	return &clientset
}
