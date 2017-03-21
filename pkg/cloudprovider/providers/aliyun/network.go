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
	"fmt"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/golang/glog"
	"kubeup.com/archon/pkg/cluster"
	"kubeup.com/archon/pkg/util"
	"net"
)

type AliyunNetwork struct {
	VPC           string `k8s:"vpc-id"`
	VSwitch       string `k8s:"vswitch-id"`
	SecurityGroup string `k8s:"security-group-id"`
	RouteTable    string `k8s:"route-table-id"`
	Router        string `k8s:"router-id"`
	VPCIPRange    string `k8s:"vpc-ip-range"`
}

func (p *aliyunCloud) EnsureNetwork(clusterName string, network *cluster.Network) (status *cluster.NetworkStatus, err error) {
	an := &AliyunNetwork{}
	defer func() {
		if err != nil {
			err2 := p.deleteNetwork(network, an)
			if err2 != nil {
				glog.Errorf("Unable to rollback network creation: %v", err2)
			}
		}
	}()
	if network.Annotations == nil {
		network.Annotations = make(map[string]string)
	}

	err = util.MapToStruct(network.Annotations, an, AliyunAnnotationPrefix)
	if err != nil {
		return
	}

	if an.VPC == "" {
		err = p.createVPC(clusterName, network, an)
		if err != nil {
			return
		}
	}

	if an.VSwitch == "" {
		err = p.createVSwitch(clusterName, network, an)
		if err != nil {
			return
		}
	}

	err = util.StructToMap(an, network.Annotations, AliyunAnnotationPrefix)
	if err != nil {
		return
	}

	status = &cluster.NetworkStatus{
		Phase: cluster.NetworkRunning,
	}

	return
}

func (p *aliyunCloud) EnsureNetworkDeleted(clusterName string, network *cluster.Network) (err error) {
	an := &AliyunNetwork{}
	if network.Annotations == nil {
		return nil
	}

	err = util.MapToStruct(network.Annotations, an, AliyunAnnotationPrefix)
	if err != nil {
		return
	}

	return p.deleteNetwork(network, an)
}

func (p *aliyunCloud) deleteNetwork(network *cluster.Network, an *AliyunNetwork) (err error) {
	if an.VSwitch != "" {
		err = p.deleteVSwitch(an)
		if err != nil {
			if isNotFound(err) {
				err = nil
			} else {
				return
			}
		}
	}

	if an.VPC != "" {
		err = p.deleteVPC(network, an)
		if err != nil {
			if isNotFound(err) {
				err = nil
			} else {
				return
			}
		}
	}

	return nil
}

func (p *aliyunCloud) AddNetworkAnnotation(clusterName string, instance *cluster.Instance, network *cluster.Network) error {
	if instance == nil || network == nil {
		return fmt.Errorf("instance or network is nil")
	}

	if network.Annotations == nil {
		return fmt.Errorf("Aliyun network is not ready")
	}

	if instance.Annotations == nil {
		instance.Annotations = make(map[string]string)
	}

	// Network Spec
	err := util.StructToMap(network.Spec, instance.Annotations, cluster.AnnotationPrefix)
	if err != nil {
		return err
	}

	aliyunNetwork := AliyunNetwork{}
	err = util.MapToStruct(network.Annotations, &aliyunNetwork, AliyunAnnotationPrefix)
	if err != nil {
		return err
	}

	return util.StructToMap(aliyunNetwork, instance.Annotations, AliyunAnnotationPrefix)
}

func (p *aliyunCloud) createVPC(clusterName string, network *cluster.Network, an *AliyunNetwork) (err error) {
	if an.VPCIPRange == "" {
		return fmt.Errorf("vpn-ip-range should not be empty")
	}
	region := common.Region(network.Spec.Region)

	glog.V(4).Infof("Creating VPC in aliyun for %s (%v)", clusterName, network.Spec)
	resp, err := p.ecs.CreateVpc(&ecs.CreateVpcArgs{
		RegionId:    region,
		CidrBlock:   an.VPCIPRange,
		VpcName:     clusterName + "-vpc",
		Description: "Archon managed vpc",
		ClientToken: util.RandNano(),
	})
	if err != nil {
		err = aliyunSafeError(err)
		return
	}

	an.VPC = resp.VpcId
	an.Router = resp.VRouterId
	an.RouteTable = resp.RouteTableId

	err2 := p.ecs.WaitForVpcAvailable(region, an.VPC, 0)
	if err2 != nil {
		err2 = aliyunSafeError(err)
		glog.V(4).Infof("Warning: waiting for vpc to be available time out: %v. Ignored", err2)
	}

	// Create sg
	sg, err := p.ecs.CreateSecurityGroup(&ecs.CreateSecurityGroupArgs{
		RegionId:          region,
		SecurityGroupName: clusterName + "-sg",
		VpcId:             an.VPC,
		ClientToken:       util.RandNano(),
	})
	if err != nil {
		return aliyunSafeError(err)
	}
	an.SecurityGroup = sg

	err = p.ecs.AuthorizeSecurityGroup(&ecs.AuthorizeSecurityGroupArgs{
		SecurityGroupId: an.SecurityGroup,
		RegionId:        region,
		NicType:         ecs.NicTypeIntranet,
		IpProtocol:      ecs.IpProtocolTCP,
		Policy:          ecs.PermissionPolicyAccept,
		PortRange:       "1/65535",
		SourceCidrIp:    "0.0.0.0/0",
	})
	if err != nil {
		return aliyunSafeError(err)
	}

	return
}

func (p *aliyunCloud) createVSwitch(clusterName string, network *cluster.Network, an *AliyunNetwork) (err error) {
	_, sub, _ := net.ParseCIDR(network.Spec.Subnet)
	sub2 := util.FromIPNet(sub)

	args := &ecs.CreateVSwitchArgs{
		ZoneId:      network.Spec.Zone,
		CidrBlock:   sub2.String(),
		VpcId:       an.VPC,
		VSwitchName: clusterName + "-vswitch",
	}
	vswitch, err := p.ecs.CreateVSwitch(args)
	if err != nil {
		return aliyunSafeError(err)
	}
	an.VSwitch = vswitch
	return
}

func (p *aliyunCloud) deleteVSwitch(an *AliyunNetwork) (err error) {
	err = p.ecs.DeleteVSwitch(an.VSwitch)
	if err != nil && isNotFound(err) {
		err = nil
	}

	return err
}

func (p *aliyunCloud) deleteVPC(network *cluster.Network, an *AliyunNetwork) (err error) {
	region := common.Region(network.Spec.Region)
	err = p.ecs.DeleteSecurityGroup(region, an.SecurityGroup)
	if err != nil && !isNotFound(err) {
		return aliyunSafeError(err)
	}

	err = p.ecs.DeleteVpc(an.VPC)
	if err != nil && !isNotFound(err) {
		return aliyunSafeError(err)
	}

	return nil
}
