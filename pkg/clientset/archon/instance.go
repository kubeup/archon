package archon

import (
	//"github.com/golang/glog"
	"k8s.io/kubernetes/pkg/api"
	rest "k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/watch"
	"kubeup.com/archon/pkg/cluster"
)

// InstancesGetter has a method to return a InstanceInterface.
// A group's client should implement this interface.
type InstancesGetter interface {
	Instances(namespace string) InstanceInterface
}

// InstanceInterface has methods to work with Instance resources.
type InstanceInterface interface {
	Create(*cluster.Instance) (*cluster.Instance, error)
	Update(*cluster.Instance) (*cluster.Instance, error)
	UpdateStatus(*cluster.Instance) (*cluster.Instance, error)
	Delete(name string) error
	Get(name string) (*cluster.Instance, error)
	List(api.ListOptions) (*cluster.InstanceList, error)
	Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *cluster.Instance, err error)
	Watch(api.ListOptions) (watch.Interface, error)
}

// instances implements InstanceInterface
type instances struct {
	client rest.Interface
	ns     string
}

// newInstances returns a Instances
func newInstances(c *ArchonClient, namespace string) *instances {
	return &instances{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Create takes the representation of a instance and creates it.  Returns the server's representation of the instance, and an error, if there is any.
func (c *instances) Create(instance *cluster.Instance) (result *cluster.Instance, err error) {
	result = &cluster.Instance{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("instances").
		Body(instance).
		Do().
		Into(result)
	return
}

// Update takes the representation of a instance and updates it. Returns the server's representation of the instance, and an error, if there is any.
func (c *instances) Update(instance *cluster.Instance) (result *cluster.Instance, err error) {
	result = &cluster.Instance{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("instances").
		Name(instance.Name).
		Body(instance).
		Do().
		Into(result)
	return
}

func (c *instances) UpdateStatus(instance *cluster.Instance) (result *cluster.Instance, err error) {
	result = &cluster.Instance{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("instances").
		Name(instance.Name).
		Body(instance).
		Do().
		Into(result)
	return
}

// Delete takes name of the instance and deletes it. Returns an error if one occurs.
func (c *instances) Delete(name string) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("instances").
		Name(name).
		Do().
		Error()
}

// Get takes name of the instance, and returns the corresponding instance object, and an error if there is any.
func (c *instances) Get(name string) (result *cluster.Instance, err error) {
	result = &cluster.Instance{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("instances").
		Name(name).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Instances that match those selectors.
func (c *instances) List(options api.ListOptions) (result *cluster.InstanceList, err error) {
	result = &cluster.InstanceList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("instances").
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested instances.
func (c *instances) Watch(options api.ListOptions) (watch.Interface, error) {
	return c.client.Get().
		Prefix("watch").
		Namespace(c.ns).
		Resource("instances").
		Watch()
}

// Patch applies the patch and returns the patched replicaSet.
func (c *instances) Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *cluster.Instance, err error) {
	result = &cluster.Instance{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("instances").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
