package cache

import (
	"fmt"

	"kubeup.com/archon/pkg/cluster"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/errors"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/client/cache"
	"k8s.io/kubernetes/pkg/labels"
)

// StoreToInstanceLister helps list instances
type StoreToInstanceLister struct {
	Indexer cache.Indexer
}

func (s *StoreToInstanceLister) List(selector labels.Selector) (ret []*cluster.Instance, err error) {
	err = cache.ListAll(s.Indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*cluster.Instance))
	})
	return ret, err
}

func (s *StoreToInstanceLister) Instances(namespace string) storeInstancesNamespacer {
	return storeInstancesNamespacer{s.Indexer, namespace}
}

type storeInstancesNamespacer struct {
	indexer   cache.Indexer
	namespace string
}

func (s storeInstancesNamespacer) List(selector labels.Selector) (ret []*cluster.Instance, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*cluster.Instance))
	})
	return ret, err
}

func (s storeInstancesNamespacer) Get(name string) (*cluster.Instance, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(api.Resource("instance"), name)
	}
	return obj.(*cluster.Instance), nil
}

// StoreToInstanceGroupLister helps list instances
type StoreToInstanceGroupLister struct {
	Indexer cache.Indexer
}

func (s *StoreToInstanceGroupLister) List(selector labels.Selector) (ret []*cluster.InstanceGroup, err error) {
	err = cache.ListAll(s.Indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*cluster.InstanceGroup))
	})
	return ret, err
}

func (s *StoreToInstanceGroupLister) InstanceGroups(namespace string) storeInstanceGroupsNamespacer {
	return storeInstanceGroupsNamespacer{s.Indexer, namespace}
}

type storeInstanceGroupsNamespacer struct {
	indexer   cache.Indexer
	namespace string
}

func (s storeInstanceGroupsNamespacer) List(selector labels.Selector) (ret []*cluster.InstanceGroup, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*cluster.InstanceGroup))
	})
	return ret, err
}

func (s storeInstanceGroupsNamespacer) Get(name string) (*cluster.InstanceGroup, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(api.Resource("instancegroup"), name)
	}
	return obj.(*cluster.InstanceGroup), nil
}

// GetPodReplicaSets returns a list of ReplicaSets managing a pod. Returns an error only if no matching ReplicaSets are found.
func (s *StoreToInstanceGroupLister) GetInstanceGroup(instance *cluster.Instance) (igs []*cluster.InstanceGroup, err error) {
	if len(instance.Labels) == 0 {
		err = fmt.Errorf("no InstanceGroups found for instance %v because it has no labels", instance.Name)
		return
	}

	list, err := s.InstanceGroups(instance.Namespace).List(labels.Everything())
	if err != nil {
		return
	}
	for _, ig := range list {
		if ig.Namespace != instance.Namespace {
			continue
		}
		selector, err := unversioned.LabelSelectorAsSelector(ig.Spec.Selector)
		if err != nil {
			return nil, fmt.Errorf("invalid selector: %v", err)
		}

		// If a ReplicaSet with a nil or empty selector creeps in, it should match nothing, not everything.
		if selector.Empty() || !selector.Matches(labels.Set(instance.Labels)) {
			continue
		}
		igs = append(igs, ig)
	}
	if len(igs) == 0 {
		err = fmt.Errorf("could not find InstanceGroup for instance %s in namespace %s with labels: %v", instance.Name, instance.Namespace, instance.Labels)
	}
	return
}
