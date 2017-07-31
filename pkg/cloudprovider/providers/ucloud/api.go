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

package ucloud

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/ucloud/ucloud-sdk-go/service/uhost"
	"github.com/ucloud/ucloud-sdk-go/service/unet"
	"k8s.io/apimachinery/pkg/util/wait"
	"time"
)

type UCloudInterface interface {
	CreateInstance(*uhost.CreateUHostInstanceParams) (string, error)
	StopInstance(string, bool) error
	StartInstance(string) error
	DeleteInstance(string) error
	DescribeInstances(args *uhost.DescribeUHostInstanceParams) (instances []uhost.UHostSet, err error)
	DescribeInstanceAttribute(string) (*uhost.UHostSet, error)

	ReinstallInstance(string, string) error

	//
	//WaitForInstance(string, uhost.InstanceStatus, int) error
	//ModifyInstanceAttribute(*uhost.ModifyInstanceAttributeArgs) error
	//DescribeDisks(args *uhost.DescribeDisksArgs) (disks []uhost.DiskItemType, err error)
	//ReInitDisk(string) error

	//AllocatePublicIpAddress(string) (string, error)
	AllocateEipAddress(*unet.AllocateEIPParams) (*unet.EIPSet, error)
	//WaitForEip(common.Region, string, uhost.EipStatus, int) error
	AssociateEipAddress(string, string) error
	DescribeEipAddress(string) (*unet.EIPSet, error)
	UnassociateEipAddress(string, string) error
	ReleaseEipAddress(string) error

	CreateSecurityGroup(*unet.CreateSecurityGroupParams) error
	DeleteSecurityGroup(int) error

	/*CreateVSwitch(*uhost.CreateVSwitchArgs) (string, error)
	DeleteVSwitch(string) error
	CreateVpc(*uhost.CreateVpcArgs) (*uhost.CreateVpcResponse, error)
	DeleteVpc(string) error
	WaitForVpcAvailable(common.Region, string, int) error
	*/

	ModifyRemark(string, string) error
	WaitForInstance(id string, status string, timeout int) error
}

type ucloudWrapper struct {
	uhost *uhost.UHost
	unet  *unet.UNet
}

var _ UCloudInterface = &ucloudWrapper{}

func (p *ucloudWrapper) CreateInstance(input *uhost.CreateUHostInstanceParams) (id string, err error) {
	if input.Region == "" {
		input.Region = p.uhost.Config.Region
	}

	resp, err := p.uhost.CreateUHostInstance(input)
	if resp != nil && len(resp.UHostIds) > 0 {
		id = resp.UHostIds[0]
	}
	return id, err
}

func (p *ucloudWrapper) StopInstance(id string, force bool) (err error) {
	// TODO: add force mode
	if force == false {
		_, err = p.uhost.StopUHostInstance(&uhost.StopUHostInstanceParams{
			Region:  p.uhost.Config.Region,
			UHostId: id,
		})
	} else {
		_, err = p.uhost.PoweroffUHostInstance(&uhost.PoweroffUHostInstanceParams{
			Region:  p.uhost.Config.Region,
			UHostId: id,
		})
	}
	return err
}

func (p *ucloudWrapper) StartInstance(id string) error {
	_, err := p.uhost.StartUHostInstance(&uhost.StartUHostInstanceParams{
		Region:  p.uhost.Config.Region,
		UHostId: id,
	})
	return err
}

func (p *ucloudWrapper) DeleteInstance(id string) error {
	_, err := p.uhost.TerminateUHostInstance(&uhost.TerminateUHostInstanceParams{
		Region:  p.uhost.Config.Region,
		UHostId: id,
	})
	return err
}

func (p *ucloudWrapper) DescribeInstances(args *uhost.DescribeUHostInstanceParams) (instances []uhost.UHostSet, err error) {
	offset := 0
	limit := 20
	total := 9999

	if args.Region == "" {
		args.Region = p.uhost.Config.Region
	}

	for len(instances) < total {
		args.Offset = offset
		args.Limit = limit
		resp, err := p.uhost.DescribeUHostInstance(args)
		if err != nil {
			return nil, err
		}
		instances = append(instances, resp.UHostSet...)
		total = resp.TotalCount
	}

	return
}

func (p *ucloudWrapper) DescribeInstanceAttribute(id string) (*uhost.UHostSet, error) {
	ret, err := p.DescribeInstances(&uhost.DescribeUHostInstanceParams{
		Region:   p.uhost.Config.Region,
		UHostIds: []string{id},
	})

	if len(ret) > 0 {
		return &ret[0], err
	}

	return nil, err
}

func (p *ucloudWrapper) ReinstallInstance(id string, password string) error {
	_, err := p.uhost.ReinstallUHostInstance(&uhost.ReinstallUHostInstanceParams{
		Region:   p.uhost.Config.Region,
		UHostId:  id,
		Password: password,
	})
	return err
}

func (p *ucloudWrapper) AllocateEipAddress(input *unet.AllocateEIPParams) (*unet.EIPSet, error) {
	if input.Region == "" {
		input.Region = p.uhost.Config.Region
	}

	resp, err := p.unet.AllocateEIP(input)
	if err != nil {
		return nil, err
	}
	if resp.EIPSet == nil || len(*resp.EIPSet) == 0 {
		return nil, fmt.Errorf("UCloud returned an empty eip set.")
	}

	return &(*resp.EIPSet)[0], nil
}

func (p *ucloudWrapper) AssociateEipAddress(instanceId string, eipId string) error {
	_, err := p.unet.BindEIP(&unet.BindEIPParams{
		Region:       p.uhost.Config.Region,
		EIPId:        eipId,
		ResourceId:   instanceId,
		ResourceType: ResourceTypeUHost,
	})
	return err
}

func (p *ucloudWrapper) UnassociateEipAddress(instanceId string, eipId string) error {
	_, err := p.unet.UnBindEIP(&unet.UnBindEIPParams{
		Region:       p.uhost.Config.Region,
		EIPId:        eipId,
		ResourceId:   instanceId,
		ResourceType: ResourceTypeUHost,
	})
	return err
}

func (p *ucloudWrapper) ReleaseEipAddress(id string) error {
	_, err := p.unet.ReleaseEIP(&unet.ReleaseEIPParams{
		Region: p.uhost.Config.Region,
		EIPId:  id,
	})
	return err
}

func (p *ucloudWrapper) DescribeEipAddress(id string) (eip *unet.EIPSet, err error) {
	resp, err := p.unet.DescribeEIP(&unet.DescribeEIPParams{
		Region: p.uhost.Config.Region,
		EIPIds: []string{id},
	})
	if resp != nil && resp.EIPSet != nil && len(*(resp.EIPSet)) > 0 {
		eip = &(*(resp.EIPSet))[0]
	}
	return
}

func (p *ucloudWrapper) CreateSecurityGroup(input *unet.CreateSecurityGroupParams) error {
	if input.Region == "" {
		input.Region = p.uhost.Config.Region
	}
	_, err := p.unet.CreateSecurityGroup(input)
	return err
}

func (p *ucloudWrapper) DeleteSecurityGroup(id int) error {
	_, err := p.unet.DeleteSecurityGroup(&unet.DeleteSecurityGroupParams{
		Region:  p.uhost.Config.Region,
		GroupId: id,
	})
	return err
}

func (p *ucloudWrapper) ModifyRemark(id string, remark string) error {
	_, err := p.uhost.ModifyUHostInstanceRemark(&uhost.ModifyUHostInstanceRemarkParams{
		Region:  p.uhost.Config.Region,
		UHostId: id,
		Remark:  remark,
	})
	return err
}

func (p *ucloudWrapper) WaitForInstance(id string, status string, timeout int) error {
	if timeout == 0 {
		timeout = DefaultWaitTimeout
	}
	glog.V(2).Infof("Start waiting for instance status: %s %s", id, status)
	return wait.PollImmediate(1*time.Second, time.Duration(timeout)*time.Second, func() (bool, error) {
		instance, err := p.DescribeInstanceAttribute(id)
		if err != nil {
			glog.V(2).Infof("Waiting error:", err)
			return false, err
		}

		if instance == nil {
			glog.V(2).Infof("Waiting error: Instance is not found")
			return false, ErrorNotFound
		}

		glog.V(2).Infof("Waiting %s %s", instance.State, status)
		return instance.State == status, err
	})
}
