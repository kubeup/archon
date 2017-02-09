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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	k8saws "k8s.io/kubernetes/pkg/cloudprovider/providers/aws"
	"kubeup.com/archon/pkg/cluster"
	"kubeup.com/archon/pkg/util"
	"time"
)

type AWSNetworkOption struct {
	NameServers string `k8s:"name-servers"`
	DomainName  string `k8s:"domain-name"`
}

type AWSNetwork struct {
	VPC             string `k8s:"vpc-id"`
	InternetGateway string `k8s:"internet-gateway-id"`
	Subnet          string `k8s:"subnet-id"`
	RouteTable      string `k8s:"route-table-id"`
	Labels          map[string]string
	Annotations     map[string]string
}

func (p *awsCloud) EnsureNetwork(clusterName string, network *cluster.Network) (status *cluster.NetworkStatus, err error) {
	ano := &AWSNetworkOption{}
	err = util.MapToStruct(network.Annotations, ano, AWSAnnotationPrefix)
	if err != nil {
		return
	}

	an := &AWSNetwork{}
	if network.Annotations == nil {
		network.Annotations = make(map[string]string)
	}

	err = util.MapToStruct(network.Annotations, an, AWSAnnotationPrefix)
	if err != nil {
		return
	}

	an.Labels = network.Labels
	an.Annotations = network.Annotations

	if an.VPC == "" {
		an.VPC, err = p.createVPC(clusterName, network)
		if err != nil {
			return
		}
	}

	if ano.NameServers != "" || ano.DomainName != "" {
		err = p.configureDhcp(an, ano)
		if err != nil {
			return
		}
	}

	if an.InternetGateway == "" {
		an.InternetGateway, err = p.createInternetGateway(an)
		if err != nil {
			return
		}
	}

	if an.RouteTable == "" {
		an.RouteTable, err = p.createRouteTable(an)
		if err != nil {
			return
		}
	}

	if an.Subnet == "" {
		an.Subnet, err = p.createSubnet(an, network)
		if err != nil {
			return
		}
	}

	err = util.StructToMap(an, network.Annotations, AWSAnnotationPrefix)
	if err != nil {
		return
	}

	status = &cluster.NetworkStatus{
		Phase: cluster.NetworkRunning,
	}

	return
}

func (p *awsCloud) EnsureNetworkDeleted(clusterName string, network *cluster.Network) (err error) {
	an := &AWSNetwork{}
	if network.Annotations == nil {
		network.Annotations = make(map[string]string)
	}

	err = util.MapToStruct(network.Annotations, an, AWSAnnotationPrefix)
	if err != nil {
		return
	}

	if an.Subnet != "" {
		err = p.deleteSubnet(an)
		if err != nil {
			if isNotExistError(err) {
				err = nil
			} else {
				return
			}
		}
	}

	if an.InternetGateway != "" {
		err = p.deleteInternetGateway(an)
		if err != nil {
			if isNotExistError(err) {
				err = nil
			} else {
				return
			}
		}
	}

	if an.VPC != "" {
		err = p.deleteVPC(an)
		if err != nil {
			if isNotExistError(err) {
				err = nil
			} else {
				return
			}
		}
	}

	return nil
}

func (p *awsCloud) AddNetworkAnnotation(clusterName string, instance *cluster.Instance, network *cluster.Network) error {
	if instance == nil || network == nil {
		return fmt.Errorf("instance or network is nil")
	}

	if instance.Annotations == nil {
		instance.Annotations = make(map[string]string)
	}

	// Network Spec
	err := util.StructToMap(network.Spec, instance.Annotations, cluster.AnnotationPrefix)
	if err != nil {
		return err
	}

	// AWS Network
	if network.Annotations == nil {
		return fmt.Errorf("AWS network is not ready")
	}

	awsnetwork := AWSNetwork{}
	err = util.MapToStruct(network.Annotations, &awsnetwork, AWSAnnotationPrefix)
	if err != nil {
		return err
	}

	return util.StructToMap(awsnetwork, instance.Annotations, AWSAnnotationPrefix)
}

func (p *awsCloud) configureDhcp(an *AWSNetwork, ano *AWSNetworkOption) (err error) {
	var c []*ec2.NewDhcpConfiguration
	if ano.DomainName != "" {
		c = append(c, &ec2.NewDhcpConfiguration{
			Key: aws.String("domain-name"),
			Values: []*string{
				aws.String(ano.DomainName),
			},
		})
	}
	if ano.NameServers != "" {
		c = append(c, &ec2.NewDhcpConfiguration{
			Key: aws.String("domain-name-servers"),
			Values: []*string{
				aws.String(ano.NameServers),
			},
		})
	}

	params := &ec2.CreateDhcpOptionsInput{
		DhcpConfigurations: c,
	}
	resp, err := p.ec2.CreateDhcpOptions(params)
	if err != nil {
		return err
	}
	params2 := &ec2.AssociateDhcpOptionsInput{
		DhcpOptionsId: resp.DhcpOptions.DhcpOptionsId,
		VpcId:         aws.String(an.VPC),
	}
	_, err = p.ec2.AssociateDhcpOptions(params2)
	return
}

func (p *awsCloud) createVPC(clusterName string, network *cluster.Network) (vpcID string, err error) {
	r, err := p.ec2.CreateVpc(&ec2.CreateVpcInput{
		CidrBlock: aws.String(network.Spec.Subnet),
	})

	if err != nil {
		return
	}

	vpcID = *r.Vpc.VpcId
	time.Sleep(10 * time.Second)

	_, err = p.ec2.ModifyVpcAttribute(&ec2.ModifyVpcAttributeInput{
		EnableDnsHostnames: &ec2.AttributeBooleanValue{
			Value: aws.Bool(true),
		},
		VpcId: r.Vpc.VpcId,
	})

	if err != nil {
		return
	}

	params2 := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("vpc-id"),
				Values: []*string{
					r.Vpc.VpcId, // Required
				},
			},
		},
	}
	r2, err := p.ec2.DescribeSecurityGroups(params2)
	if err != nil {
		return
	}

	awssg := *r2.SecurityGroups[0].GroupId
	params3 := &ec2.AuthorizeSecurityGroupIngressInput{
		CidrIp:     aws.String("0.0.0.0/0"),
		FromPort:   aws.Int64(0),
		GroupId:    aws.String(awssg),
		IpProtocol: aws.String("-1"),
		ToPort:     aws.Int64(65535),
	}

	_, err = p.ec2.AuthorizeSecurityGroupIngress(params3)
	if err != nil {
		return
	}

	return
}

func (p *awsCloud) createInternetGateway(an *AWSNetwork) (igID string, err error) {
	r, err := p.ec2.CreateInternetGateway(&ec2.CreateInternetGatewayInput{})
	if err != nil {
		return
	}

	igID = *r.InternetGateway.InternetGatewayId
	time.Sleep(10 * time.Second)

	params2 := &ec2.AttachInternetGatewayInput{
		InternetGatewayId: r.InternetGateway.InternetGatewayId, // Required
		VpcId:             aws.String(an.VPC),                  // Required
	}
	_, err = p.ec2.AttachInternetGateway(params2)
	if err != nil {
		return
	}

	return
}

func (p *awsCloud) createRouteTable(an *AWSNetwork) (rtID string, err error) {
	params3 := &ec2.DescribeRouteTablesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("vpc-id"),
				Values: []*string{
					aws.String(an.VPC), // Required
				},
			},
		},
	}

	resp, err := p.ec2.DescribeRouteTables(params3)
	if err != nil {
		return
	}

	rtID = *resp.RouteTables[0].RouteTableId
	params4 := &ec2.CreateRouteInput{
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
		RouteTableId:         aws.String(rtID),
		GatewayId:            aws.String(an.InternetGateway),
	}
	_, err = p.ec2.CreateRoute(params4)
	if err != nil {
		return
	}

	if cn, ok := an.Labels[k8saws.TagNameKubernetesCluster]; ok {
		params5 := &ec2.CreateTagsInput{
			Resources: []*string{
				resp.RouteTables[0].RouteTableId,
			},
			Tags: []*ec2.Tag{
				{
					Key:   aws.String(k8saws.TagNameKubernetesCluster),
					Value: aws.String(cn),
				},
			},
		}
		_, err = p.ec2.CreateTags(params5)
		if err != nil {
			return
		}

	}
	return
}

func (p *awsCloud) createSubnet(an *AWSNetwork, network *cluster.Network) (subnetID string, err error) {
	params := &ec2.CreateSubnetInput{
		VpcId:            aws.String(an.VPC),
		CidrBlock:        aws.String(network.Spec.Subnet),
		AvailabilityZone: aws.String(network.Spec.Zone),
	}

	r, err := p.ec2.CreateSubnet(params)
	if err != nil {
		return
	}

	subnetID = *r.Subnet.SubnetId
	time.Sleep(10 * time.Second)

	params2 := &ec2.ModifySubnetAttributeInput{
		SubnetId: r.Subnet.SubnetId,
		MapPublicIpOnLaunch: &ec2.AttributeBooleanValue{
			Value: aws.Bool(true),
		},
	}
	_, err = p.ec2.ModifySubnetAttribute(params2)
	if err != nil {
		return
	}

	return
}

func (p *awsCloud) deleteSubnet(an *AWSNetwork) (err error) {
	params := &ec2.DeleteSubnetInput{
		SubnetId: aws.String(an.Subnet), // Required
	}
	_, err = p.ec2.DeleteSubnet(params)
	return err
}

func (p *awsCloud) deleteInternetGateway(an *AWSNetwork) (err error) {
	params := &ec2.DetachInternetGatewayInput{
		InternetGatewayId: aws.String(an.InternetGateway), // Required
		VpcId:             aws.String(an.VPC),             // Required
	}
	_, err = p.ec2.DetachInternetGateway(params)
	if err != nil && !isNotExistError(err) {
		return
	}

	_, err = p.ec2.DeleteInternetGateway(&ec2.DeleteInternetGatewayInput{
		InternetGatewayId: aws.String(an.InternetGateway),
	})
	return err
}

func (p *awsCloud) deleteVPC(an *AWSNetwork) (err error) {
	_, err = p.ec2.DeleteVpc(&ec2.DeleteVpcInput{
		VpcId: aws.String(an.VPC),
	})
	return err
}
