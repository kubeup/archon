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
	"encoding/base64"
	"fmt"
	aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/glog"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"kubeup.com/archon/pkg/cluster"
	"kubeup.com/archon/pkg/userdata"
	"kubeup.com/archon/pkg/util"
)

type InstanceOptions struct {
	// AWS use instance profile to grant permissions to ec2 instances.
	InstanceProfile string `k8s:"instance-profile"`
}

var ErrorNotFound = fmt.Errorf("Instance is not found")

var StateMap = map[string]cluster.InstancePhase{
	ec2.InstanceStateNamePending:      cluster.InstancePending,
	ec2.InstanceStateNameRunning:      cluster.InstanceRunning,
	ec2.InstanceStateNameShuttingDown: cluster.InstanceFailed,
	ec2.InstanceStateNameTerminated:   cluster.InstanceFailed,
	ec2.InstanceStateNameStopping:     cluster.InstanceFailed,
	ec2.InstanceStateNameStopped:      cluster.InstanceFailed,
}

func instanceToStatus(i *ec2.Instance) *cluster.InstanceStatus {
	phase, ok := StateMap[destring(i.State.Name)]
	if !ok {
		glog.Infof("Unknown instance state: %+v", i.State)
		phase = cluster.InstanceUnknown
	}
	return &cluster.InstanceStatus{
		Phase:             phase,
		PrivateIP:         destring(i.PrivateIpAddress),
		PublicIP:          destring(i.PublicIpAddress),
		InstanceID:        destring(i.InstanceId),
		CreationTimestamp: unversioned.NewTime(detime(i.LaunchTime)),
	}
}

func (p *awsCloud) ListInstances(clusterName string, network *cluster.Network, selector map[string]string) (names []string, statuses []*cluster.InstanceStatus, err error) {
	awsnetwork := AWSNetwork{}
	err = util.MapToStruct(network.Annotations, &awsnetwork, AWSAnnotationPrefix)
	if err != nil {
		err = fmt.Errorf("Network is not ready. Can't list instances: %s", err.Error())
		return
	}

	filters := []*ec2.Filter{
		&ec2.Filter{
			Name:   aws.String("vpc-id"),
			Values: []*string{&awsnetwork.VPC},
		},
	}

	for k, v := range selector {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String(fmt.Sprintf("tag:%s", k)),
			Values: []*string{&v},
		})
	}

	instances, err := p.ec2.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: filters,
	})

	if err != nil {
		return
	}

	for _, instance := range instances {
		name := ""
		for _, tag := range instance.Tags {
			if *tag.Key == NameKey {
				name = *tag.Value
				break
			}
		}
		names = append(names, name)
		statuses = append(statuses, instanceToStatus(instance))
	}

	return
}

func (p *awsCloud) GetInstance(clusterName string, instance *cluster.Instance) (status *cluster.InstanceStatus, err error) {
	awsnetwork := AWSNetwork{}
	err = util.MapToStruct(instance.Dependency.Network.Annotations, &awsnetwork, AWSAnnotationPrefix)
	if err != nil {
		err = fmt.Errorf("Network is not ready. Can't get instance: %s", err.Error())
		return
	}

	return p.getInstance(awsnetwork, instance.Status.InstanceID)
}

func (p *awsCloud) getInstance(awsnetwork AWSNetwork, instanceID string) (status *cluster.InstanceStatus, err error) {
	if instanceID == "" {
		return nil, ErrorNotFound
	}

	instances, err := p.ec2.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)},
	})

	if err != nil {
		return
	}

	if len(instances) == 0 {
		return nil, ErrorNotFound
	}

	return instanceToStatus(instances[0]), nil
}

func (p *awsCloud) EnsureInstance(clusterName string, instance *cluster.Instance) (status *cluster.InstanceStatus, err error) {
	awsnetwork := AWSNetwork{}
	err = util.MapToStruct(instance.Annotations, &awsnetwork, AWSAnnotationPrefix)
	if err != nil {
		err = fmt.Errorf("Network is not ready. Can't create instance: %s", err.Error())
		return
	}

	if instance.Status.InstanceID != "" {
		status2 := (*cluster.InstanceStatus)(nil)
		status2, err = p.getInstance(awsnetwork, instance.Status.InstanceID)

		if err != nil {
			if err == ErrorNotFound {
				return p.createInstance(clusterName, instance)
			}
			return
		}

		switch status2.Phase {
		case cluster.InstanceFailed, cluster.InstanceUnknown:
			err = p.EnsureInstanceDeleted(clusterName, instance)
			if err != nil {
				return
			}
			return p.createInstance(clusterName, instance)
		}

		status = status2
	} else {
		return p.createInstance(clusterName, instance)
	}

	return
}

func (p *awsCloud) createInstance(clusterName string, instance *cluster.Instance) (status *cluster.InstanceStatus, err error) {
	nifIDs := []string{}
	defer func() {
		// Clean up nifs if create failed
		if err == nil || len(nifIDs) == 0 {
			return
		}

		err2 := p.deleteNetworkInterfaces(nifIDs)
		if err2 != nil {
			glog.Errorf("Can't delete network interface on error: %+v", err2)
		}
	}()

	options := cluster.InstanceOptions{}
	if instance.Labels != nil {
		err = util.MapToStruct(instance.Labels, &options, cluster.AnnotationPrefix)
		if err != nil {
			err = fmt.Errorf("Can't get instance options: %s", err.Error())
			return
		}
	}

	awsnetwork := AWSNetwork{}
	err = util.MapToStruct(instance.Annotations, &awsnetwork, AWSAnnotationPrefix)
	if err != nil || awsnetwork.Subnet == "" || awsnetwork.VPC == "" {
		err = fmt.Errorf("Can't get network from instance annotations: %+v", err)
		return
	}

	networkSpec := cluster.NetworkSpec{}
	err = util.MapToStruct(instance.Annotations, &networkSpec, cluster.AnnotationPrefix)
	if err != nil {
		err = fmt.Errorf("Can't get network from instance annotations: %s", err.Error())
		return
	}

	awsoptions := InstanceOptions{}
	err = util.MapToStruct(instance.Annotations, &awsoptions, AWSAnnotationPrefix)
	if err != nil {
		err = fmt.Errorf("Can't get aws instance options: %s", err.Error())
		return
	}

	eip := EIP{}
	pip := PrivateIP{}
	nif := ""
	awsPrivateIP := (*string)(nil)
	ifSpecs := ([]*ec2.InstanceNetworkInterfaceSpecification)(nil)
	subnetID := (*string)(nil)

	if options.PreallocatePrivateIP {
		err = util.MapToStruct(instance.Annotations, &pip, AWSAnnotationPrefix)
		if err != nil {
			err = fmt.Errorf("Can't get private ip interface from instance annotations: %s", err.Error())
			return
		}

		if instance.Status.PrivateIP == "" {
			err = fmt.Errorf("custom private IP is not provided.")
			return
		}

		nif = pip.NetworkInterfaceID
	}

	if options.PreallocatePublicIP {
		err = util.MapToStruct(instance.Annotations, &eip, AWSAnnotationPrefix)
		if err != nil {
			err = fmt.Errorf("Can't get eip from instance annotations: %s", err.Error())
			return
		}

		if eip.AllocationID == "" {
			err = fmt.Errorf("EIP is not provided.")
			return
		}

		if nif == "" {
			// Create if
			resp, err := p.ec2.CreateNetworkInterface(&ec2.CreateNetworkInterfaceInput{
				PrivateIpAddress: awsPrivateIP,
				SubnetId:         aws.String(awsnetwork.Subnet),
			})
			if err != nil {
				return nil, fmt.Errorf("Error creating network interface: %+v", err)
			}
			nif = destring(resp.NetworkInterface.NetworkInterfaceId)
			nifIDs = append(nifIDs, nif)
		}

		// Associate address
		_, err = p.ec2.AssociateAddress(&ec2.AssociateAddressInput{
			AllocationId:       aws.String(eip.AllocationID),
			NetworkInterfaceId: aws.String(nif),
		})
		if err != nil {
			return nil, fmt.Errorf("Error associating eip with network interface: %+v", err)
		}

		// Add to ifspecs
		ifSpecs = append(ifSpecs, &ec2.InstanceNetworkInterfaceSpecification{
			//DeleteOnTermination: aws.Bool(true),
			DeviceIndex:        aws.Int64(0),
			NetworkInterfaceId: aws.String(nif),
		})

	} else if nif != "" {
		// Add to ifspecs
		ifSpecs = append(ifSpecs, &ec2.InstanceNetworkInterfaceSpecification{
			//AssociatePublicIpAddress: aws.Bool(true),
			NetworkInterfaceId: aws.String(nif),
			DeviceIndex:        aws.Int64(0),
		})
	} else {
		subnetID = aws.String(awsnetwork.Subnet)
	}

	awsInstanceType := "t2.small"
	if instance.Spec.InstanceType != "" {
		awsInstanceType = instance.Spec.InstanceType
	}

	u, err := userdata.Generate(instance)
	if err != nil {
		return nil, err
	}
	s := base64.StdEncoding.EncodeToString(u)

	// Image and its root device
	regionProfile, ok := AWSRegions[networkSpec.Region]
	if !ok {
		return nil, fmt.Errorf("Invalid region: %s", networkSpec.Region)
	}

	image := regionProfile.HVM
	if instance.Spec.Image != "" {
		image = instance.Spec.Image
	}

	params := &ec2.RunInstancesInput{
		ImageId:      aws.String(image), // Required
		MaxCount:     aws.Int64(1),      // Required
		MinCount:     aws.Int64(1),      // Required
		ClientToken:  aws.String(util.RandNano()),
		InstanceType: aws.String(awsInstanceType),
		SubnetId:     subnetID,
		UserData:     aws.String(s),
		Placement: &ec2.Placement{
			AvailabilityZone: aws.String(networkSpec.Zone),
		},
		NetworkInterfaces: ifSpecs,
	}

	if awsoptions.InstanceProfile != "" {
		params.IamInstanceProfile = &ec2.IamInstanceProfileSpecification{
			Arn: aws.String(awsoptions.InstanceProfile),
		}
	}

	r, err := p.ec2.RunInstances(params)
	if err != nil {
		return
	}
	vps := r.Instances[0]
	vpsID := destring(vps.InstanceId)

	params2 := &ec2.ModifyInstanceAttributeInput{
		InstanceId:      aws.String(vpsID),
		SourceDestCheck: &ec2.AttributeBooleanValue{Value: aws.Bool(false)},
	}

	_, err = p.ec2.ModifyInstanceAttribute(params2)
	if err != nil {
		return
	}

	tags := []*ec2.Tag{
		{
			Key:   aws.String(NameKey),
			Value: aws.String(instance.Name),
		},
	}
	for k, v := range instance.Labels {
		tags = append(tags, &ec2.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}

	params3 := &ec2.CreateTagsInput{
		Resources: []*string{ // Required
			aws.String(vpsID),
		},
		Tags: tags}
	_, err = p.ec2.CreateTags(params3)
	if err != nil {
		return
	}

	status = instanceToStatus(vps)
	glog.Infof("New instance created %+v", status)
	nifIDs = []string{}
	return
}

func (p *awsCloud) EnsureInstanceDeleted(clusterName string, instance *cluster.Instance) (err error) {
	if instance.Status.InstanceID == "" {
		return nil
	}

	return p.deleteInstance(instance.Status.InstanceID)
}

func (p *awsCloud) deleteInstance(vpsId string) (err error) {
	resp, err := p.ec2.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(vpsId)},
	})

	if err != nil {
		if isNotExistError(err) {
			err = nil
		}
		return
	}

	if len(resp) == 0 {
		return
	}

	vps := resp[0]
	nifs := []string{}
	for _, nif := range vps.NetworkInterfaces {
		nifs = append(nifs, destring(nif.NetworkInterfaceId))
	}

	_, err = p.ec2.TerminateInstances(&ec2.TerminateInstancesInput{
		InstanceIds: []*string{aws.String(vpsId)},
	})
	if err != nil {
		if isNotExistError(err) {
			err = nil
		}
		return
	}

	err = p.ec2.WaitUntilInstanceTerminated(&ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(vpsId)},
	})
	if err != nil {
		return
	}

	return p.deleteNetworkInterfaces(nifs)
}

func (p *awsCloud) deleteNetworkInterfaces(nifs []string) (err error) {
	for _, nif := range nifs {
		_, err = p.ec2.DeleteNetworkInterface(&ec2.DeleteNetworkInterfaceInput{
			NetworkInterfaceId: aws.String(nif),
		})
		if err != nil {
			if isNotExistError(err) {
				err = nil
			}
			return
		}
	}

	return
}
