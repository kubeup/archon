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

package aliyun

import (
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
)

type ECS interface {
	CreateInstance(*ecs.CreateInstanceArgs) (string, error)
	StopInstance(string, bool) error
	StartInstance(string) error
	DeleteInstance(string) error
	DescribeInstances(args *ecs.DescribeInstancesArgs) (instances []ecs.InstanceAttributesType, err error)
	DescribeInstanceAttribute(string) (*ecs.InstanceAttributesType, error)
	WaitForInstance(string, ecs.InstanceStatus, int) error
	ModifyInstanceAttribute(*ecs.ModifyInstanceAttributeArgs) error

	AllocatePublicIpAddress(string) (string, error)
	AllocateEipAddress(*ecs.AllocateEipAddressArgs) (string, string, error)
	WaitForEip(common.Region, string, ecs.EipStatus, int) error
	AssociateEipAddress(string, string) error
	UnassociateEipAddress(string, string) error
	ReleaseEipAddress(string) error

	CreateSecurityGroup(*ecs.CreateSecurityGroupArgs) (string, error)
	DeleteSecurityGroup(common.Region, string) error
	AuthorizeSecurityGroup(*ecs.AuthorizeSecurityGroupArgs) error
	AuthorizeSecurityGroupEgress(*ecs.AuthorizeSecurityGroupEgressArgs) error
	CreateVSwitch(*ecs.CreateVSwitchArgs) (string, error)
	DeleteVSwitch(string) error
	CreateVpc(*ecs.CreateVpcArgs) (*ecs.CreateVpcResponse, error)
	DeleteVpc(string) error
	WaitForVpcAvailable(common.Region, string, int) error

	AddTags(*ecs.AddTagsArgs) error
}

type aliyunECS struct {
	ecs *ecs.Client
}

var _ ECS = &aliyunECS{}

func (p *aliyunECS) CreateInstance(input *ecs.CreateInstanceArgs) (string, error) {
	return p.ecs.CreateInstance(input)
}

func (p *aliyunECS) StopInstance(id string, force bool) error {
	return p.ecs.StopInstance(id, force)
}

func (p *aliyunECS) StartInstance(id string) error {
	return p.ecs.StartInstance(id)
}

func (p *aliyunECS) DeleteInstance(id string) error {
	return p.ecs.DeleteInstance(id)
}

func (p *aliyunECS) DescribeInstances(args *ecs.DescribeInstancesArgs) (instances []ecs.InstanceAttributesType, err error) {
	pg := &common.Pagination{
		PageNumber: 0,
		PageSize:   50,
	}

	for {
		ret, pgresp, err := p.ecs.DescribeInstances(args)
		if err != nil {
			return nil, err
		}
		instances = append(instances, ret...)
		pg = pgresp.NextPage()
		if pg == nil {
			break
		}
	}

	return
}

func (p *aliyunECS) DescribeInstanceAttribute(id string) (*ecs.InstanceAttributesType, error) {
	return p.ecs.DescribeInstanceAttribute(id)
}

func (p *aliyunECS) WaitForInstance(id string, status ecs.InstanceStatus, timeout int) error {
	return p.ecs.WaitForInstance(id, status, timeout)
}
func (p *aliyunECS) ModifyInstanceAttribute(input *ecs.ModifyInstanceAttributeArgs) error {
	return p.ecs.ModifyInstanceAttribute(input)
}

func (p *aliyunECS) AllocatePublicIpAddress(id string) (string, error) {
	return p.ecs.AllocatePublicIpAddress(id)
}

func (p *aliyunECS) AllocateEipAddress(input *ecs.AllocateEipAddressArgs) (string, string, error) {
	return p.ecs.AllocateEipAddress(input)
}

func (p *aliyunECS) WaitForEip(region common.Region, id string, status ecs.EipStatus, timeout int) error {
	return p.ecs.WaitForEip(region, id, status, timeout)
}

func (p *aliyunECS) AssociateEipAddress(instanceId string, eipId string) error {
	return p.ecs.AssociateEipAddress(instanceId, eipId)
}

func (p *aliyunECS) UnassociateEipAddress(instanceId string, eipId string) error {
	return p.ecs.UnassociateEipAddress(instanceId, eipId)
}

func (p *aliyunECS) ReleaseEipAddress(id string) error {
	return p.ecs.ReleaseEipAddress(id)
}

func (p *aliyunECS) CreateSecurityGroup(input *ecs.CreateSecurityGroupArgs) (string, error) {
	return p.ecs.CreateSecurityGroup(input)
}

func (p *aliyunECS) DeleteSecurityGroup(region common.Region, id string) error {
	return p.ecs.DeleteSecurityGroup(region, id)
}

func (p *aliyunECS) AuthorizeSecurityGroup(input *ecs.AuthorizeSecurityGroupArgs) error {
	return p.ecs.AuthorizeSecurityGroup(input)
}

func (p *aliyunECS) AuthorizeSecurityGroupEgress(input *ecs.AuthorizeSecurityGroupEgressArgs) error {
	return p.ecs.AuthorizeSecurityGroupEgress(input)
}

func (p *aliyunECS) CreateVSwitch(input *ecs.CreateVSwitchArgs) (string, error) {
	return p.ecs.CreateVSwitch(input)
}

func (p *aliyunECS) DeleteVSwitch(id string) error {
	return p.ecs.DeleteVSwitch(id)
}

func (p *aliyunECS) CreateVpc(input *ecs.CreateVpcArgs) (*ecs.CreateVpcResponse, error) {
	return p.ecs.CreateVpc(input)
}

func (p *aliyunECS) DeleteVpc(id string) error {
	return p.ecs.DeleteVpc(id)
}

func (p *aliyunECS) WaitForVpcAvailable(region common.Region, id string, timeout int) error {
	return p.ecs.WaitForVpcAvailable(region, id, timeout)
}

func (p *aliyunECS) AddTags(args *ecs.AddTagsArgs) error {
	return p.ecs.AddTags(args)
}
