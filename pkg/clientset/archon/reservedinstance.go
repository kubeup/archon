/*
Copyright 2017 The Archon Authors.
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
	//"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	"kubeup.com/archon/pkg/cluster"
)

// ReservedInstancesGetter has a method to return a ReservedInstanceInterface.
// A group's client should implement this interface.
type ReservedInstancesGetter interface {
	ReservedInstances(namespace string) ReservedInstanceInterface
}

// ReservedInstanceInterface has methods to work with ReservedInstance resources.
type ReservedInstanceInterface interface {
	Create(*cluster.ReservedInstance) (*cluster.ReservedInstance, error)
	Update(*cluster.ReservedInstance) (*cluster.ReservedInstance, error)
	UpdateStatus(*cluster.ReservedInstance) (*cluster.ReservedInstance, error)
	Delete(name string) error
	Get(name string) (*cluster.ReservedInstance, error)
	List(metav1.ListOptions) (*cluster.ReservedInstanceList, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *cluster.ReservedInstance, err error)
	Watch(metav1.ListOptions) (watch.Interface, error)
}

// reservedInstances implements ReservedInstanceInterface
type reservedInstances struct {
	client rest.Interface
	ns     string
}

// newReservedInstances returns a ReservedInstances
func newReservedInstances(c *ArchonClient, namespace string) *reservedInstances {
	return &reservedInstances{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Create takes the representation of a instance and creates it.  Returns the server's representation of the instance, and an error, if there is any.
func (c *reservedInstances) Create(instance *cluster.ReservedInstance) (result *cluster.ReservedInstance, err error) {
	result = &cluster.ReservedInstance{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("reservedinstances").
		Body(instance).
		Do().
		Into(result)
	return
}

// Update takes the representation of a instance and updates it. Returns the server's representation of the instance, and an error, if there is any.
func (c *reservedInstances) Update(instance *cluster.ReservedInstance) (result *cluster.ReservedInstance, err error) {
	result = &cluster.ReservedInstance{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("reservedinstances").
		Name(instance.Name).
		Body(instance).
		Do().
		Into(result)
	return
}

func (c *reservedInstances) UpdateStatus(instance *cluster.ReservedInstance) (result *cluster.ReservedInstance, err error) {
	result = &cluster.ReservedInstance{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("reservedinstances").
		Name(instance.Name).
		Body(instance).
		Do().
		Into(result)
	return
}

// Delete takes name of the instance and deletes it. Returns an error if one occurs.
func (c *reservedInstances) Delete(name string) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("reservedinstances").
		Name(name).
		Do().
		Error()
}

// Get takes name of the instance, and returns the corresponding instance object, and an error if there is any.
func (c *reservedInstances) Get(name string) (result *cluster.ReservedInstance, err error) {
	result = &cluster.ReservedInstance{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("reservedinstances").
		Name(name).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ReservedInstances that match those selectors.
func (c *reservedInstances) List(options metav1.ListOptions) (result *cluster.ReservedInstanceList, err error) {
	result = &cluster.ReservedInstanceList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("reservedinstances").
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested reservedInstances.
func (c *reservedInstances) Watch(options metav1.ListOptions) (watch.Interface, error) {
	return c.client.Get().
		Prefix("watch").
		Namespace(c.ns).
		Resource("reservedinstances").
		Watch()
}

// Patch applies the patch and returns the patched replicaSet.
func (c *reservedInstances) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *cluster.ReservedInstance, err error) {
	result = &cluster.ReservedInstance{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("reservedinstances").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
