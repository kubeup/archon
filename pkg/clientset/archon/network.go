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
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	"kubeup.com/archon/pkg/cluster"
)

// NetworksGetter has a method to return a NetworkInterface.
// A group's client should implement this interface.
type NetworksGetter interface {
	Networks(namespace string) NetworkInterface
}

// NetworkInterface has methods to work with Network resources.
type NetworkInterface interface {
	Create(*cluster.Network) (*cluster.Network, error)
	Update(*cluster.Network) (*cluster.Network, error)
	Delete(name string) error
	Get(name string) (*cluster.Network, error)
	List() (*cluster.NetworkList, error)
	Watch() (watch.Interface, error)
}

// networks implements NetworkInterface
type networks struct {
	client rest.Interface
	ns     string
}

// newNetworks returns a Networks
func newNetworks(c *ArchonClient, namespace string) *networks {
	return &networks{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Create takes the representation of a network and creates it.  Returns the server's representation of the network, and an error, if there is any.
func (c *networks) Create(network *cluster.Network) (result *cluster.Network, err error) {
	result = &cluster.Network{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("networks").
		Body(network).
		Do().
		Into(result)
	return
}

// Update takes the representation of a network and updates it. Returns the server's representation of the network, and an error, if there is any.
func (c *networks) Update(network *cluster.Network) (result *cluster.Network, err error) {
	result = &cluster.Network{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("networks").
		Name(network.Name).
		Body(network).
		Do().
		Into(result)
	return
}

// Delete takes name of the network and deletes it. Returns an error if one occurs.
func (c *networks) Delete(name string) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("networks").
		Name(name).
		Do().
		Error()
}

// Get takes name of the network, and returns the corresponding network object, and an error if there is any.
func (c *networks) Get(name string) (result *cluster.Network, err error) {
	result = &cluster.Network{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("networks").
		Name(name).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Networks that match those selectors.
func (c *networks) List() (result *cluster.NetworkList, err error) {
	result = &cluster.NetworkList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("networks").
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested networks.
func (c *networks) Watch() (watch.Interface, error) {
	return c.client.Get().
		Prefix("watch").
		Namespace(c.ns).
		Resource("networks").
		Watch()
}
