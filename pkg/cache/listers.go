package cache

import (
	"kubeup.com/archon/pkg/cluster"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/errors"
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
