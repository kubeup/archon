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
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubeup.com/archon/pkg/clientset"
	"kubeup.com/archon/pkg/cloudprovider"
	"kubeup.com/archon/pkg/cluster"
	"kubeup.com/archon/pkg/initializer"
	"kubeup.com/archon/pkg/util"

	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"reflect"
)

var (
	PublicIPToken  = cluster.AnnotationPrefix + "public-ip"
	PrivateIPToken = cluster.AnnotationPrefix + "private-ip"
)

type IPInitializer struct {
	clusterName string
	kubeClient  clientset.Interface
	archon      cloudprovider.ArchonInterface
	token       string
	sync        func(obj initializer.Object, finalizing bool) (initializer.Object, error, bool)
}

var _ initializer.Initializer = &IPInitializer{}

func (ec *IPInitializer) Token() string {
	return ec.token
}

func (ec *IPInitializer) Initialize(obj initializer.Object) (updatedObj initializer.Object, err error, retryable bool) {
	return ec.sync(obj, false)
}

func (ec *IPInitializer) Finalize(obj initializer.Object) (updatedObj initializer.Object, err error, retryable bool) {
	return ec.sync(obj, true)
}

func NewPrivateIPInitializer(kubeClient clientset.Interface, cloud cloudprovider.Interface, clusterName string) (initializer.Initializer, error) {
	c := &IPInitializer{
		clusterName: clusterName,
		kubeClient:  kubeClient,
		token:       PrivateIPToken,
	}
	c.sync = c.syncPrivateIP
	if cloud != nil {
		c.archon, _ = cloud.Archon()
	}
	return c, nil
}

func NewPublicIPInitializer(kubeClient clientset.Interface, cloud cloudprovider.Interface, clusterName string) (initializer.Initializer, error) {
	c := &IPInitializer{
		clusterName: clusterName,
		kubeClient:  kubeClient,
		token:       PublicIPToken,
	}
	c.sync = c.syncPublicIP
	if cloud != nil {
		c.archon, _ = cloud.Archon()
	}
	return c, nil
}

func (ec *IPInitializer) syncPublicIP(obj initializer.Object, deleting bool) (updatedObj initializer.Object, err error, retryable bool) {
	instance, _ := obj.(*cluster.Instance)
	if ec.archon == nil {
		err = fmt.Errorf("cloudprovider doesn't support archon interface. aborting")
		return
	}

	glog.V(2).Infof("Syncing Public IP %s", instance.Name)

	pip, supported := ec.archon.PublicIP()
	if supported == false {
		err = fmt.Errorf("Instance wants preallocated Public IP but the cloudprovider doesn't support it")
		return
	}

	options := cluster.InstanceOptions{}
	err = util.MapToStruct(instance.Labels, &options, cluster.AnnotationPrefix)
	if err != nil {
		err = fmt.Errorf("Can't get instance options: %s", err.Error())
		return
	}

	if options.UseInstanceID != "" {
		err = fmt.Errorf("IP initializer is not supported on preallocated instances")
		return
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
		initializer.RemoveInitializer(instance, PublicIPToken)
		initializer.AddFinalizer(instance, PublicIPToken)
		ret, err2 := ec.kubeClient.Archon().Instances(instance.Namespace).Update(instance)
		if err2 != nil {
			err = fmt.Errorf("Not able to persist instance after Public IP update: %s", err2.Error())
			retryable = false
		} else {
			updatedObj = ret
		}
	} else if deleting {
		initializer.RemoveFinalizer(instance, PublicIPToken)
		_, err2 := ec.kubeClient.Archon().Instances(instance.Namespace).Get(instance.Name)
		if err2 != nil {
			if errors.IsNotFound(err2) {
				updatedObj = instance
			} else {
				err = fmt.Errorf("Not able to get instance status: %s", err2.Error())
				retryable = false
			}
		} else {
			ret, err3 := ec.kubeClient.Archon().Instances(instance.Namespace).Update(instance)
			if err3 != nil {
				err = fmt.Errorf("Not able to persist instance after Public IP update: %s", err3.Error())
				retryable = false
			} else {
				updatedObj = ret
			}
		}
	}

	return
}

func (ec *IPInitializer) syncPrivateIP(obj initializer.Object, deleting bool) (updatedObj initializer.Object, err error, retryable bool) {
	// Aliyun provider requires public ip to be provisioned first.
	if initializer.HasInitializer(obj, PublicIPToken) {
		err = initializer.ErrSkip
		return
	}

	instance, _ := obj.(*cluster.Instance)
	if ec.archon == nil {
		err = fmt.Errorf("cloudprovider doesn't support archon interface. aborting")
		return
	}

	glog.V(2).Infof("Syncing Private IP %s", instance.Name)

	pip, supported := ec.archon.PrivateIP()
	if supported == false {
		err = fmt.Errorf("Instance wants preallocated Private IP but the cloudprovider doesn't support it")
		return
	}

	options := cluster.InstanceOptions{}
	err = util.MapToStruct(instance.Labels, &options, cluster.AnnotationPrefix)
	if err != nil {
		err = fmt.Errorf("Can't get instance options: %s", err.Error())
		return
	}

	if options.UseInstanceID != "" {
		err = fmt.Errorf("IP initializer is not supported on preallocated instances")
		return
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
		glog.Infof("Ensured private ip: %+v", instance.Status)
	}

	if err != nil {
		retryable = true
		return
	}

	if !deleting && (!reflect.DeepEqual(previousAnnotations, instance.Annotations) || !cluster.InstanceStatusEqual(previousStatus, instance.Status)) {
		// Persist instance
		initializer.RemoveInitializer(instance, PrivateIPToken)
		initializer.AddFinalizer(instance, PrivateIPToken)
		ret, err2 := ec.kubeClient.Archon().Instances(instance.Namespace).Update(instance)
		if err2 != nil {
			err = fmt.Errorf("Not able to persist instance after Private IP update: %s", err2.Error())
			retryable = false
		} else {
			ret.Dependency = instance.Dependency
			updatedObj = ret
		}
	} else if deleting {
		initializer.RemoveFinalizer(instance, PrivateIPToken)
		_, err2 := ec.kubeClient.Archon().Instances(instance.Namespace).Get(instance.Name)
		if err2 != nil {
			if errors.IsNotFound(err2) {
				updatedObj = instance
			} else {
				err = fmt.Errorf("Not able to get instance status: %s", err2.Error())
				retryable = false
			}
		} else {
			ret, err3 := ec.kubeClient.Archon().Instances(instance.Namespace).Update(instance)
			if err3 != nil {
				err = fmt.Errorf("Not able to persist instance after Public IP update: %s", err3.Error())
				retryable = false
			} else {
				updatedObj = ret
			}
		}
	}

	return
}
