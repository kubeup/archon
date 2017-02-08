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

package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/s3"
)

type EC2 interface {
	DescribeImages(input *ec2.DescribeImagesInput) (*ec2.DescribeImagesOutput, error)
	RunInstances(input *ec2.RunInstancesInput) (*ec2.Reservation, error)
	WaitUntilInstanceTerminated(input *ec2.DescribeInstancesInput) error
	ModifyInstanceAttribute(input *ec2.ModifyInstanceAttributeInput) (*ec2.ModifyInstanceAttributeOutput, error)
	DescribeInstances(request *ec2.DescribeInstancesInput) ([]*ec2.Instance, error)
	TerminateInstances(input *ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error)

	CreateVpc(input *ec2.CreateVpcInput) (*ec2.CreateVpcOutput, error)
	ModifyVpcAttribute(input *ec2.ModifyVpcAttributeInput) (*ec2.ModifyVpcAttributeOutput, error)
	DescribeSecurityGroups(input *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error)
	AuthorizeSecurityGroupIngress(input *ec2.AuthorizeSecurityGroupIngressInput) (*ec2.AuthorizeSecurityGroupIngressOutput, error)
	CreateInternetGateway(input *ec2.CreateInternetGatewayInput) (*ec2.CreateInternetGatewayOutput, error)
	AttachInternetGateway(input *ec2.AttachInternetGatewayInput) (*ec2.AttachInternetGatewayOutput, error)
	DescribeRouteTables(input *ec2.DescribeRouteTablesInput) (*ec2.DescribeRouteTablesOutput, error)
	CreateRoute(input *ec2.CreateRouteInput) (*ec2.CreateRouteOutput, error)
	CreateSubnet(input *ec2.CreateSubnetInput) (*ec2.CreateSubnetOutput, error)
	ModifySubnetAttribute(input *ec2.ModifySubnetAttributeInput) (*ec2.ModifySubnetAttributeOutput, error)
	DeleteVpc(input *ec2.DeleteVpcInput) (*ec2.DeleteVpcOutput, error)
	DeleteInternetGateway(input *ec2.DeleteInternetGatewayInput) (*ec2.DeleteInternetGatewayOutput, error)
	DetachInternetGateway(input *ec2.DetachInternetGatewayInput) (*ec2.DetachInternetGatewayOutput, error)
	DeleteSubnet(input *ec2.DeleteSubnetInput) (*ec2.DeleteSubnetOutput, error)

	CreateNetworkInterface(input *ec2.CreateNetworkInterfaceInput) (*ec2.CreateNetworkInterfaceOutput, error)
	DeleteNetworkInterface(input *ec2.DeleteNetworkInterfaceInput) (*ec2.DeleteNetworkInterfaceOutput, error)
	AssociateAddress(input *ec2.AssociateAddressInput) (*ec2.AssociateAddressOutput, error)
	AllocateAddress(input *ec2.AllocateAddressInput) (*ec2.AllocateAddressOutput, error)
	ReleaseAddress(input *ec2.ReleaseAddressInput) (*ec2.ReleaseAddressOutput, error)

	CreateTags(input *ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error)
}

type S3 interface {
	PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
}

type IAM interface {
	GetInstanceProfile(input *iam.GetInstanceProfileInput) (*iam.GetInstanceProfileOutput, error)
}

type awsSdkEC2 struct {
	ec2 *ec2.EC2
}

var _ EC2 = &awsSdkEC2{}

func (p *awsSdkEC2) DescribeImages(input *ec2.DescribeImagesInput) (*ec2.DescribeImagesOutput, error) {
	return p.ec2.DescribeImages(input)
}

func (p *awsSdkEC2) RunInstances(input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
	return p.ec2.RunInstances(input)
}

func (p *awsSdkEC2) CreateNetworkInterface(input *ec2.CreateNetworkInterfaceInput) (*ec2.CreateNetworkInterfaceOutput, error) {
	return p.ec2.CreateNetworkInterface(input)
}

func (p *awsSdkEC2) DeleteNetworkInterface(input *ec2.DeleteNetworkInterfaceInput) (*ec2.DeleteNetworkInterfaceOutput, error) {
	return p.ec2.DeleteNetworkInterface(input)
}

func (p *awsSdkEC2) AllocateAddress(input *ec2.AllocateAddressInput) (*ec2.AllocateAddressOutput, error) {
	return p.ec2.AllocateAddress(input)
}

func (p *awsSdkEC2) ReleaseAddress(input *ec2.ReleaseAddressInput) (*ec2.ReleaseAddressOutput, error) {
	return p.ec2.ReleaseAddress(input)
}

func (p *awsSdkEC2) AssociateAddress(input *ec2.AssociateAddressInput) (*ec2.AssociateAddressOutput, error) {
	return p.ec2.AssociateAddress(input)
}

func (p *awsSdkEC2) WaitUntilInstanceTerminated(input *ec2.DescribeInstancesInput) error {
	return p.ec2.WaitUntilInstanceTerminated(input)
}

func (p *awsSdkEC2) CreateTags(input *ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
	return p.ec2.CreateTags(input)
}

func (p *awsSdkEC2) ModifyInstanceAttribute(input *ec2.ModifyInstanceAttributeInput) (*ec2.ModifyInstanceAttributeOutput, error) {
	return p.ec2.ModifyInstanceAttribute(input)
}

func (p *awsSdkEC2) DescribeInstances(request *ec2.DescribeInstancesInput) ([]*ec2.Instance, error) {
	// Instances are paged
	results := []*ec2.Instance{}
	var nextToken *string

	for {
		response, err := p.ec2.DescribeInstances(request)
		if err != nil {
			return nil, fmt.Errorf("error listing AWS instances: %v", err)
		}

		for _, reservation := range response.Reservations {
			results = append(results, reservation.Instances...)
		}

		nextToken = response.NextToken
		if nextToken == nil || *nextToken == "" {
			break
		}
		request.NextToken = nextToken
	}

	return results, nil
}

func (p *awsSdkEC2) TerminateInstances(input *ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error) {
	return p.ec2.TerminateInstances(input)
}

func (p *awsSdkEC2) CreateVpc(input *ec2.CreateVpcInput) (*ec2.CreateVpcOutput, error) {
	return p.ec2.CreateVpc(input)
}

func (p *awsSdkEC2) ModifyVpcAttribute(input *ec2.ModifyVpcAttributeInput) (*ec2.ModifyVpcAttributeOutput, error) {
	return p.ec2.ModifyVpcAttribute(input)
}

func (p *awsSdkEC2) DescribeSecurityGroups(input *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
	return p.ec2.DescribeSecurityGroups(input)
}

func (p *awsSdkEC2) AuthorizeSecurityGroupIngress(input *ec2.AuthorizeSecurityGroupIngressInput) (*ec2.AuthorizeSecurityGroupIngressOutput, error) {
	return p.ec2.AuthorizeSecurityGroupIngress(input)
}

func (p *awsSdkEC2) CreateInternetGateway(input *ec2.CreateInternetGatewayInput) (*ec2.CreateInternetGatewayOutput, error) {
	return p.ec2.CreateInternetGateway(input)
}

func (p *awsSdkEC2) AttachInternetGateway(input *ec2.AttachInternetGatewayInput) (*ec2.AttachInternetGatewayOutput, error) {
	return p.ec2.AttachInternetGateway(input)
}

func (p *awsSdkEC2) DescribeRouteTables(input *ec2.DescribeRouteTablesInput) (*ec2.DescribeRouteTablesOutput, error) {
	return p.ec2.DescribeRouteTables(input)
}

func (p *awsSdkEC2) CreateRoute(input *ec2.CreateRouteInput) (*ec2.CreateRouteOutput, error) {
	return p.ec2.CreateRoute(input)
}

func (p *awsSdkEC2) CreateSubnet(input *ec2.CreateSubnetInput) (*ec2.CreateSubnetOutput, error) {
	return p.ec2.CreateSubnet(input)
}

func (p *awsSdkEC2) ModifySubnetAttribute(input *ec2.ModifySubnetAttributeInput) (*ec2.ModifySubnetAttributeOutput, error) {
	return p.ec2.ModifySubnetAttribute(input)
}

func (p *awsSdkEC2) DeleteVpc(input *ec2.DeleteVpcInput) (*ec2.DeleteVpcOutput, error) {
	return p.ec2.DeleteVpc(input)
}

func (p *awsSdkEC2) DetachInternetGateway(input *ec2.DetachInternetGatewayInput) (*ec2.DetachInternetGatewayOutput, error) {
	return p.ec2.DetachInternetGateway(input)
}

func (p *awsSdkEC2) DeleteInternetGateway(input *ec2.DeleteInternetGatewayInput) (*ec2.DeleteInternetGatewayOutput, error) {
	return p.ec2.DeleteInternetGateway(input)
}

func (p *awsSdkEC2) DeleteSubnet(input *ec2.DeleteSubnetInput) (*ec2.DeleteSubnetOutput, error) {
	return p.ec2.DeleteSubnet(input)
}

type awsSdkIAM struct {
	iam *iam.IAM
}

func (c *awsSdkIAM) GetInstanceProfile(input *iam.GetInstanceProfileInput) (*iam.GetInstanceProfileOutput, error) {
	return c.iam.GetInstanceProfile(input)
}

type awsSdkS3 struct {
	s3 *s3.S3
}

func (s *awsSdkS3) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return s.s3.PutObject(input)
}
