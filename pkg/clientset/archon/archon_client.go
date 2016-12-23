package archon

import (
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/runtime/serializer"
)

type ArchonInterface interface {
	RESTClient() restclient.Interface
	InstancesGetter
	InstanceGroupsGetter
	NetworksGetter
	UsersGetter
}

// StorageV1beta1Client is used to interact with features provided by the k8s.io/kubernetes/pkg/apimachinery/registered.Group group.
type ArchonClient struct {
	restClient restclient.Interface
}

func (c *ArchonClient) Instances(namespace string) InstanceInterface {
	return newInstances(c, namespace)
}

func (c *ArchonClient) InstanceGroups(namespace string) InstanceGroupInterface {
	return newInstanceGroups(c, namespace)
}

func (c *ArchonClient) Networks(namespace string) NetworkInterface {
	return newNetworks(c, namespace)
}

func (c *ArchonClient) Users(namespace string) UserInterface {
	return newUsers(c, namespace)
}

// NewForConfig creates a new StorageV1beta1Client for the given config.
func NewForConfig(c *restclient.Config) (*ArchonClient, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := restclient.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &ArchonClient{client}, nil
}

// NewForConfigOrDie creates a new StorageV1beta1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *restclient.Config) *ArchonClient {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new StorageV1beta1Client for the given RESTClient.
func New(c restclient.Interface) *ArchonClient {
	return &ArchonClient{c}
}

func setConfigDefaults(config *restclient.Config) error {
	config.APIPath = "/apis"
	config.GroupVersion = &SchemeGroupVersion
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: serializer.NewCodecFactory(Scheme)}
	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *ArchonClient) RESTClient() restclient.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
