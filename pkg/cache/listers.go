/*
Copyright 2016 The Kubernetes Authors.
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

// This file is modified from the original kubernetes source tree

package cache

import (
	"fmt"

	"kubeup.com/archon/pkg/cluster"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	"k8s.io/kubernetes/pkg/api"
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

// GetInstanceGroup returns a list of InstanceGroups managing a instance. Returns an error only if no matching InstanceGroups are found.
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
		selector, err := metav1.LabelSelectorAsSelector(ig.Spec.Selector)
		if err != nil {
			return nil, fmt.Errorf("invalid selector: %v", err)
		}

		// If a InstanceGroup with a nil or empty selector creeps in, it should match nothing, not everything.
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

// StoreToSecretLister helps list secrets
type StoreToSecretLister struct {
	Indexer cache.Indexer
}

func (s *StoreToSecretLister) List(selector labels.Selector) (ret []*api.Secret, err error) {
	err = cache.ListAll(s.Indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*api.Secret))
	})
	return ret, err
}

func (s *StoreToSecretLister) Secrets(namespace string) storeSecretsNamespacer {
	return storeSecretsNamespacer{s.Indexer, namespace}
}

type storeSecretsNamespacer struct {
	indexer   cache.Indexer
	namespace string
}

func (s storeSecretsNamespacer) List(selector labels.Selector) (ret []*api.Secret, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*api.Secret))
	})
	return ret, err
}

func (s storeSecretsNamespacer) Get(name string) (*api.Secret, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(api.Resource("secret"), name)
	}
	return obj.(*api.Secret), nil
}
