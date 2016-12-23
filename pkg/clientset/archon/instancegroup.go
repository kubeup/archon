package archon

import (
	"k8s.io/kubernetes/pkg/api"
	rest "k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/watch"
	"kubeup.com/archon/pkg/cluster"
)

// InstanceGroupsGetter has a method to return a InstanceGroupInterface.
// A group's client should implement this interface.
type InstanceGroupsGetter interface {
	InstanceGroups(namespace string) InstanceGroupInterface
}

// InstanceGroupInterface has methods to work with InstanceGroup resources.
type InstanceGroupInterface interface {
	Create(*cluster.InstanceGroup) (*cluster.InstanceGroup, error)
	Update(*cluster.InstanceGroup) (*cluster.InstanceGroup, error)
	UpdateStatus(*cluster.InstanceGroup) (*cluster.InstanceGroup, error)
	Delete(name string) error
	Get(name string) (*cluster.InstanceGroup, error)
	List(api.ListOptions) (*cluster.InstanceGroupList, error)
	Watch(api.ListOptions) (watch.Interface, error)
}

// instancegroups implements InstanceGroupInterface
type instancegroups struct {
	client rest.Interface
	ns     string
}

// newInstanceGroups returns a InstanceGroups
func newInstanceGroups(c *ArchonClient, namespace string) *instancegroups {
	return &instancegroups{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Create takes the representation of a instance and creates it.  Returns the server's representation of the instance, and an error, if there is any.
func (c *instancegroups) Create(instance *cluster.InstanceGroup) (result *cluster.InstanceGroup, err error) {
	result = &cluster.InstanceGroup{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("instancegroups").
		Body(instance).
		Do().
		Into(result)
	return
}

// Update takes the representation of a instance and updates it. Returns the server's representation of the instance, and an error, if there is any.
func (c *instancegroups) Update(instance *cluster.InstanceGroup) (result *cluster.InstanceGroup, err error) {
	result = &cluster.InstanceGroup{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("instancegroups").
		Name(instance.Name).
		Body(instance).
		Do().
		Into(result)
	return
}

func (c *instancegroups) UpdateStatus(instance *cluster.InstanceGroup) (result *cluster.InstanceGroup, err error) {
	result = &cluster.InstanceGroup{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("instancegroups").
		Name(instance.Name).
		Body(instance).
		Do().
		Into(result)
	return
}

// Delete takes name of the instance and deletes it. Returns an error if one occurs.
func (c *instancegroups) Delete(name string) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("instancegroups").
		Name(name).
		Do().
		Error()
}

// Get takes name of the instance, and returns the corresponding instance object, and an error if there is any.
func (c *instancegroups) Get(name string) (result *cluster.InstanceGroup, err error) {
	result = &cluster.InstanceGroup{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("instancegroups").
		Name(name).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of InstanceGroups that match those selectors.
func (c *instancegroups) List(options api.ListOptions) (result *cluster.InstanceGroupList, err error) {
	result = &cluster.InstanceGroupList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("instancegroups").
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested instancegroups.
func (c *instancegroups) Watch(options api.ListOptions) (watch.Interface, error) {
	return c.client.Get().
		Prefix("watch").
		Namespace(c.ns).
		Resource("instancegroups").
		Watch()
}
