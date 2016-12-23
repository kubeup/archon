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

type EIPController struct {
	clusterName string
	kubeClient  clientset.Interface
	archon      cloudprovider.ArchonInterface
}

func NewEIPController(cloud cloudprovider.Interface, kubeClient clientset.Interface, clusterName string) *EIPController {
	archon, _ := cloud.Archon()
	return &EIPController{
		clusterName: clusterName,
		kubeClient:  kubeClient,
		archon:      archon,
	}
}

func (ec *EIPController) SyncEIP(key string, instance *cluster.Instance, deleting bool) (err error, retryable bool) {
	if ec.archon == nil {
		return fmt.Errorf("cloudprovider doesn't support archon interface. aborting"), false
	}

	glog.V(2).Infof("Syncing EIP %s", key)

	options := cluster.InstanceOptions{}
	if instance.Labels != nil {
		err = util.MapToStruct(instance.Labels, &options, cluster.AnnotationPrefix)
		if err != nil {
			return
		}
	}

	eip, supported := ec.archon.EIP()
	if supported == false {
		if options.PreallocatePublicIP == true {
			return fmt.Errorf("Instance wants EIP but cloudprovider doesn't support it"), false
		}
		return nil, false
	}

	previousStatus := *cluster.InstanceStatusDeepCopy(&instance.Status)
	previousAnnotations := make(map[string]string)
	util.MapCopy(previousAnnotations, instance.Annotations)

	if options.PreallocatePublicIP == false || deleting {
		glog.V(2).Infof("Deleting EIP %s", key)
		err = eip.EnsureEIPDeleted(ec.clusterName, instance)
	} else if options.PreallocatePublicIP {
		glog.V(2).Infof("Ensuring EIP %s", key)
		_, err = eip.EnsureEIP(ec.clusterName, instance)
	}

	if err != nil {
		retryable = true
		return
	}

	if !deleting && (!reflect.DeepEqual(previousAnnotations, instance.Annotations) || cluster.InstanceStatusEqual(previousStatus, instance.Status)) {
		// Persist instance
		_, err = ec.kubeClient.Archon().Instances(instance.Namespace).Update(instance)
		if err != nil {
			err = fmt.Errorf("Not able to persist instance after EIP update: %s", err.Error())
			retryable = false
		}
	}

	return
}
