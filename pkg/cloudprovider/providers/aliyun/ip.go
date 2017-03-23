package aliyun

import (
	"fmt"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/golang/glog"
	"kubeup.com/archon/pkg/cluster"
	"kubeup.com/archon/pkg/util"
)

type EIP struct {
	AllocationID string `k8s:"eip-allocation-id"`
}

type EIPOptions struct {
	Bandwidth          int    `k8s:"eip-bandwidth"`
	InternetChargeType string `k8s:"eip-internet-charge-type"`
}

func (p *aliyunCloud) EnsurePublicIP(clusterName string, instance *cluster.Instance) (status *cluster.InstanceStatus, err error) {
	eip := EIP{}
	err = util.MapToStruct(instance.Annotations, &eip, AliyunAnnotationPrefix)
	if err != nil {
		return
	}

	if eip.AllocationID != "" {
		status = &instance.Status
		return
	}

	defer func() {
		if err != nil && eip.AllocationID != "" {
			err2 := p.ecs.ReleaseEipAddress(eip.AllocationID)
			if err2 != nil {
				glog.Errorf("Error rolling back eip address: %v", err2)
			}
		}
	}()

	eipOptions := EIPOptions{
		Bandwidth:          5,
		InternetChargeType: string(common.PayByTraffic),
	}
	err = util.MapToStruct(instance.Annotations, &eipOptions, AliyunAnnotationPrefix)
	if err != nil {
		return
	}

	networkSpec := instance.Dependency.Network.Spec
	address, allocationID, err := p.ecs.AllocateEipAddress(&ecs.AllocateEipAddressArgs{
		RegionId:           common.Region(networkSpec.Region),
		Bandwidth:          eipOptions.Bandwidth,
		InternetChargeType: common.InternetChargeType(eipOptions.InternetChargeType),
		ClientToken:        util.RandNano(),
	})
	if err != nil {
		err = fmt.Errorf("Error allocating EIP: %s", err.Error())
		return
	}

	eip.AllocationID = allocationID
	if instance.Annotations == nil {
		instance.Annotations = make(map[string]string)
	}

	err = util.StructToMap(eip, instance.Annotations, AliyunAnnotationPrefix)
	if err != nil {
		err = fmt.Errorf("Error allocating EIP: %s", err.Error())
		return
	}

	status = &instance.Status
	status.PublicIP = address

	return &instance.Status, nil
}

func (p *aliyunCloud) EnsurePublicIPDeleted(clusterName string, instance *cluster.Instance) (err error) {
	if instance.Annotations == nil {
		return
	}

	eip := EIP{}
	err = util.MapToStruct(instance.Annotations, &eip, AliyunAnnotationPrefix)
	if err != nil {
		return
	}

	networkSpec := instance.Dependency.Network.Spec

	if eip.AllocationID != "" {
		p.ecs.WaitForEip(common.Region(networkSpec.Region), eip.AllocationID, ecs.EipStatusAvailable, 0)

		err = p.ecs.ReleaseEipAddress(eip.AllocationID)

		if err != nil {
			err = fmt.Errorf("Error releasing EIP: %s", err.Error())
			return
		}

		eip.AllocationID = ""

		err = util.StructToMap(eip, instance.Annotations, AliyunAnnotationPrefix)
		if err != nil {
			err = fmt.Errorf("Error releasing EIP: %s", err.Error())
			return
		}

		instance.Status.PublicIP = ""
	}

	return nil
}

func (p *aliyunCloud) EnsurePrivateIP(clusterName string, instance *cluster.Instance) (status *cluster.InstanceStatus, err error) {
	if instance.Status.InstanceID != "" {
		status = &instance.Status
		return
	}

	status, err = p.createInstance(clusterName, instance)
	if err == nil {
		if status.PublicIP == "" {
			status.PublicIP = instance.Status.PublicIP
		}
		instance.Status = *status
	}
	return
}

func (p *aliyunCloud) EnsurePrivateIPDeleted(clusterName string, instance *cluster.Instance) (err error) {
	if instance.Status.InstanceID == "" {
		return
	}

	err = p.deleteInstance(instance.Status.InstanceID)
	if err != nil {
		return
	}

	p.ecs.WaitForInstance(instance.Status.InstanceID, ecs.Deleted, 0)
	instance.Status.InstanceID = ""

	return
}
