package instance

import (
	"fmt"
	"github.com/golang/glog"
	"kubeup.com/archon/pkg/clientset"
	"kubeup.com/archon/pkg/cloudprovider"
	"kubeup.com/archon/pkg/cluster"
	"kubeup.com/archon/pkg/util"
	"reflect"
)

type IPController struct {
	clusterName string
	kubeClient  clientset.Interface
	archon      cloudprovider.ArchonInterface
}

func NewIPController(cloud cloudprovider.Interface, kubeClient clientset.Interface, clusterName string) *IPController {
	c := &IPController{
		clusterName: clusterName,
		kubeClient:  kubeClient,
	}
	if cloud != nil {
		c.archon, _ = cloud.Archon()
	}
	return c
}

func (ec *IPController) SyncIP(key string, instance *cluster.Instance, deleting bool) (err error, retryable bool) {
	err, retryable = ec.syncPublicIP(key, instance, deleting)
	if err != nil {
		return
	}

	err, retryable = ec.syncPrivateIP(key, instance, deleting)
	return
}

func (ec *IPController) syncPublicIP(key string, instance *cluster.Instance, deleting bool) (err error, retryable bool) {
	if ec.archon == nil {
		return fmt.Errorf("cloudprovider doesn't support archon interface. aborting"), false
	}

	glog.V(2).Infof("Syncing Public IP %s", key)

	options := cluster.InstanceOptions{}
	if instance.Labels != nil {
		err = util.MapToStruct(instance.Labels, &options, cluster.AnnotationPrefix)
		if err != nil {
			return
		}
	}

	pip, supported := ec.archon.PublicIP()
	if supported == false {
		if options.PreallocatePublicIP == true {
			return fmt.Errorf("Instance wants preallocated Public IP but the cloudprovider doesn't support it"), false
		}
		return nil, false
	}

	previousStatus := *cluster.InstanceStatusDeepCopy(&instance.Status)
	previousAnnotations := (map[string]string)(nil)
	if instance.Annotations != nil {
		previousAnnotations = make(map[string]string)
		util.MapCopy(previousAnnotations, instance.Annotations)
	}

	if options.PreallocatePublicIP == false || deleting {
		glog.V(2).Infof("Deleting Public IP %s if needed", key)
		err = pip.EnsurePublicIPDeleted(ec.clusterName, instance)
	} else if options.PreallocatePublicIP {
		glog.V(2).Infof("Ensuring Public IP %s", key)
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

func (ec *IPController) syncPrivateIP(key string, instance *cluster.Instance, deleting bool) (err error, retryable bool) {
	if ec.archon == nil {
		return fmt.Errorf("cloudprovider doesn't support archon interface. aborting"), false
	}

	glog.V(2).Infof("Syncing Private IP %s", key)

	options := cluster.InstanceOptions{}
	if instance.Labels != nil {
		err = util.MapToStruct(instance.Labels, &options, cluster.AnnotationPrefix)
		if err != nil {
			return
		}
	}

	pip, supported := ec.archon.PrivateIP()
	if supported == false {
		if options.PreallocatePrivateIP == true {
			return fmt.Errorf("Instance wants preallocated Private IP but the cloudprovider doesn't support it"), false
		}
		return nil, false
	}

	previousStatus := *cluster.InstanceStatusDeepCopy(&instance.Status)
	previousAnnotations := (map[string]string)(nil)
	if instance.Annotations != nil {
		previousAnnotations = make(map[string]string)
		util.MapCopy(previousAnnotations, instance.Annotations)
	}

	if options.PreallocatePrivateIP == false || deleting {
		glog.V(2).Infof("Deleting Private IP %s if needed", key)
		err = pip.EnsurePrivateIPDeleted(ec.clusterName, instance)
	} else if options.PreallocatePrivateIP {
		glog.V(2).Infof("Ensuring Private IP %s", key)
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
