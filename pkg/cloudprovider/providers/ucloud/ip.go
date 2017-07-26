package ucloud

import (
	"fmt"
	"github.com/golang/glog"
	//"github.com/ucloud/ucloud-sdk-go/service/uhost"
	"github.com/ucloud/ucloud-sdk-go/service/unet"
	"kubeup.com/archon/pkg/cluster"
	"kubeup.com/archon/pkg/util"
)

type EIP struct {
	ID string `k8s:"eip-id"`
}

type EIPOptions struct {
	OperatorName      string `k8s:"eip-operator-name"`
	Bandwidth         int    `k8s:"eip-bandwidth"`
	ChargeType        string `k8s:"eip-charge-type"`
	PayMode           string `k8s:"eip-pay-mode"`
	SharedBandwidthID string `k8s:"eip-shared-bandwidth-id"`
}

func (p *UCloud) EnsurePublicIP(clusterName string, instance *cluster.Instance) (status *cluster.InstanceStatus, err error) {
	eip := EIP{}
	err = util.MapToStruct(instance.Annotations, &eip, UCloudAnnotationPrefix)
	if err != nil {
		return
	}

	if eip.ID != "" {
		status = &instance.Status
		return
	}

	defer func() {
		if err != nil && eip.ID != "" {
			err2 := p.ucloud.ReleaseEipAddress(eip.ID)
			if err2 != nil {
				glog.Errorf("Error rolling back eip address: %v", err2)
			}
		}
	}()

	// Default EIP options
	eo := EIPOptions{
		OperatorName: EIPOperatorBGP,
		Bandwidth:    2,
		PayMode:      EIPPayModeTraffic,
		ChargeType:   EIPChargeTypeDynamic,
	}

	err = util.MapToStruct(instance.Annotations, &eo, UCloudAnnotationPrefix)
	if err != nil {
		return
	}

	if eo.SharedBandwidthID != "" {
		eo.Bandwidth = 0
		eo.PayMode = EIPPayModeSharedBandwidth
	}

	eip0, err := p.ucloud.AllocateEipAddress(&unet.AllocateEIPParams{
		OperatorName: eo.OperatorName,
		Bandwidth:    eo.Bandwidth,
		ChargeType:   eo.ChargeType,
		// TODO Fix ucloud api
		PayMode:          eo.PayMode,
		ShareBandwidthId: eo.SharedBandwidthID,
		Name:             fmt.Sprintf("EIP-%s", instance.Name),
	})
	if err != nil {
		err = fmt.Errorf("Error allocating EIP: %s", err.Error())
		return
	}

	eip.ID = eip0.EIPId
	if instance.Annotations == nil {
		instance.Annotations = make(map[string]string)
	}

	err = util.StructToMap(eip, instance.Annotations, UCloudAnnotationPrefix)
	if err != nil {
		err = fmt.Errorf("Error allocating EIP: %s", err.Error())
		return
	}

	status = &instance.Status
	// TODO: Handle duplet ip
	status.PublicIP = (*eip0.EIPAddr)[0].IP

	return &instance.Status, nil
}

func (p *UCloud) EnsurePublicIPDeleted(clusterName string, instance *cluster.Instance) (err error) {
	if instance.Annotations == nil {
		return
	}

	eip := EIP{}
	err = util.MapToStruct(instance.Annotations, &eip, UCloudAnnotationPrefix)
	if err != nil {
		return
	}

	if eip.ID != "" {
		//p.unet.WaitForEip(common.Region(networkSpec.Region), eip.AllocationID, unet.EipStatusAvailable, 0)

		err = p.ucloud.ReleaseEipAddress(eip.ID)

		if err != nil {
			err = fmt.Errorf("Error releasing EIP: %s", err.Error())
			return
		}

		eip.ID = ""

		err = util.StructToMap(eip, instance.Annotations, UCloudAnnotationPrefix)
		if err != nil {
			err = fmt.Errorf("Error releasing EIP: %s", err.Error())
			return
		}

		instance.Status.PublicIP = ""
	}

	return nil
}

func (p *UCloud) EnsurePrivateIP(clusterName string, instance *cluster.Instance) (status *cluster.InstanceStatus, err error) {
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

func (p *UCloud) EnsurePrivateIPDeleted(clusterName string, instance *cluster.Instance) (err error) {
	if instance.Status.InstanceID == "" {
		return
	}

	err = p.deleteInstance(instance.Status.InstanceID)
	if err != nil {
		return
	}

	instance.Status.InstanceID = ""

	return
}
