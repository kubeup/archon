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

package instance

import (
	"github.com/golang/glog"
	"k8s.io/kubernetes/pkg/runtime"
	"kubeup.com/archon/pkg/clientset"
	"kubeup.com/archon/pkg/cloudprovider"
	"kubeup.com/archon/pkg/cluster"
	"kubeup.com/archon/pkg/initializer"
	"kubeup.com/archon/pkg/util"

	"fmt"
	"reflect"
)

type IPInitializer struct {
	clusterName string
	kubeClient  clientset.Interface
	archon      cloudprovider.ArchonInterface
}

var _ initializer.Initializer = &IPInitializer{}

func IPInitializerFactory(kubeClient clientset.Interface, cloud cloudprovider.Interface, clusterName string) (initializer.Initializer, error) {
	c := &IPInitializer{
		clusterName: clusterName,
		kubeClient:  kubeClient,
	}
	if cloud != nil {
		c.archon, _ = cloud.Archon()
	}
	return c, nil
}

func (ec *IPInitializer) Token() string {
	return "ip"
}

func (ec *IPInitializer) Initialize(obj runtime.Object) (updatedObj runtime.Object, err error, retryable bool) {
	return ec.sync(obj, false)
}

func (ec *IPInitializer) Finalize(obj runtime.Object) (updatedObj runtime.Object, err error, retryable bool) {
	return ec.sync(obj, true)
}

func (ec *IPInitializer) sync(obj runtime.Object, finalizing bool) (updatedObj runtime.Object, err error, retryable bool) {
	instance, _ := obj.(*cluster.Instance)
	err, retryable = ec.syncPublicIP(instance, finalizing)
	if err != nil {
		return
	}

	err, retryable = ec.syncPrivateIP(instance, finalizing)
	if err != nil {
		return
	}

	updatedObj = instance
	return
}

func (ec *IPInitializer) syncPublicIP(instance *cluster.Instance, deleting bool) (err error, retryable bool) {
	if ec.archon == nil {
		return fmt.Errorf("cloudprovider doesn't support archon interface. aborting"), false
	}

	glog.V(2).Infof("Syncing Public IP %s", instance.Name)

	pip, supported := ec.archon.PublicIP()
	if supported == false {
		return fmt.Errorf("Instance wants preallocated Public IP but the cloudprovider doesn't support it"), false
	}

	previousStatus := *cluster.InstanceStatusDeepCopy(&instance.Status)
	previousAnnotations := (map[string]string)(nil)
	if instance.Annotations != nil {
		previousAnnotations = make(map[string]string)
		util.MapCopy(previousAnnotations, instance.Annotations)
	}

	if deleting {
		glog.V(2).Infof("Deleting Public IP %s if needed", instance.Name)
		err = pip.EnsurePublicIPDeleted(ec.clusterName, instance)
	} else {
		glog.V(2).Infof("Ensuring Public IP %s", instance.Name)
		_, err = pip.EnsurePublicIP(ec.clusterName, instance)
	}

	if err != nil {
		retryable = true
		return
	}

	if !deleting && (!reflect.DeepEqual(previousAnnotations, instance.Annotations) || !cluster.InstanceStatusEqual(previousStatus, instance.Status)) {
		// Persist instance
		ret, err2 := ec.kubeClient.Archon().Instances(instance.Namespace).Update(instance)
		if err2 != nil {
			err = fmt.Errorf("Not able to persist instance after Public IP update: %s", err.Error())
			retryable = false
		} else {
			*instance = *ret
		}
	}

	return
}

func (ec *IPInitializer) syncPrivateIP(instance *cluster.Instance, deleting bool) (err error, retryable bool) {
	if ec.archon == nil {
		return fmt.Errorf("cloudprovider doesn't support archon interface. aborting"), false
	}

	glog.V(2).Infof("Syncing Private IP %s", instance.Name)

	pip, supported := ec.archon.PrivateIP()
	if supported == false {
		return fmt.Errorf("Instance wants preallocated Private IP but the cloudprovider doesn't support it"), false
	}

	previousStatus := *cluster.InstanceStatusDeepCopy(&instance.Status)
	previousAnnotations := (map[string]string)(nil)
	if instance.Annotations != nil {
		previousAnnotations = make(map[string]string)
		util.MapCopy(previousAnnotations, instance.Annotations)
	}

	if deleting {
		glog.V(2).Infof("Deleting Private IP %s if needed", instance.Name)
		err = pip.EnsurePrivateIPDeleted(ec.clusterName, instance)
	} else {
		glog.V(2).Infof("Ensuring Private IP %s", instance.Name)
		_, err = pip.EnsurePrivateIP(ec.clusterName, instance)
	}

	if err != nil {
		retryable = true
		return
	}

	if !deleting && (!reflect.DeepEqual(previousAnnotations, instance.Annotations) || !cluster.InstanceStatusEqual(previousStatus, instance.Status)) {
		// Persist instance
		ret, err2 := ec.kubeClient.Archon().Instances(instance.Namespace).Update(instance)
		if err2 != nil {
			err = fmt.Errorf("Not able to persist instance after Private IP update: %s", err.Error())
			retryable = false
		} else {
			*instance = *ret
		}
	}

	return
}
