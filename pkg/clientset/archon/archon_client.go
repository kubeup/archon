/*
Copyright 2016 The Archon Authors.
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

package archon

import (
	"k8s.io/apimachinery/pkg/runtime"
	restclient "k8s.io/client-go/rest"
	"k8s.io/kubernetes/pkg/api"
	"kubeup.com/archon/pkg/cluster"
)

type ArchonInterface interface {
	RESTClient() restclient.Interface
	InstancesGetter
	InstanceGroupsGetter
	NetworksGetter
	UsersGetter
	ReservedInstancesGetter
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

func (c *ArchonClient) ReservedInstances(namespace string) ReservedInstanceInterface {
	return newReservedInstances(c, namespace)
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
	config.GroupVersion = &cluster.SchemeGroupVersion
	config.ContentType = runtime.ContentTypeJSON
	// DirectCodec doesn't do defaulting. Since we are using TPR, there's no defaulting on
	// the server side. We will have to do it by ourselves on the clientside when decoding
	config.NegotiatedSerializer = DirectDefaultingCodecFactory{api.Codecs}
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
