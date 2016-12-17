package aws

import (
	"encoding/base64"
	"fmt"
	aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/glog"
	"kubeup.com/archon/pkg/cluster"
	"kubeup.com/archon/pkg/userdata"
	"kubeup.com/archon/pkg/util"
)

var ErrorNotFound = fmt.Errorf("Instance is not found")

var StateMap = map[string]cluster.InstancePhase{
	ec2.InstanceStateNamePending:      cluster.InstancePending,
	ec2.InstanceStateNameRunning:      cluster.InstanceRunning,
	ec2.InstanceStateNameShuttingDown: cluster.InstanceFailed,
	ec2.InstanceStateNameTerminated:   cluster.InstanceFailed,
	ec2.InstanceStateNameStopping:     cluster.InstanceFailed,
	ec2.InstanceStateNameStopped:      cluster.InstanceFailed,
}

type InstanceOptions struct {
	UseEIP          bool `k8s:"use-eip"`
	CustomPrivateIP bool `k8s:"custom-private-ip"`
}

type EIP struct {
	AllocationID string `k8s:"eip-allocation-id"`
	IP           string `k8s:"eip-ip"`
}

type PrivateIP struct {
	IP string `k8s:"private-ip"`
}

func instanceToStatus(i *ec2.Instance) *cluster.InstanceStatus {
	phase, ok := StateMap[destring(i.State.Name)]
	if !ok {
		glog.Infof("Unknown instance state: %+v", i.State)
		phase = cluster.InstanceUnknown
	}
	return &cluster.InstanceStatus{
		Phase:      phase,
		PrivateIP:  destring(i.PrivateIpAddress),
		PublicIP:   destring(i.PublicIpAddress),
		InstanceID: destring(i.InstanceId),
	}
}

func (p *awsCloud) ListInstances(clusterName string, network *cluster.Network, selector map[string]string) (names []string, statuses []*cluster.InstanceStatus, err error) {
	awsnetwork := AWSNetwork{}
	err = util.MapToStruct(network.Annotations, &awsnetwork, AnnotationPrefix)
	if err != nil {
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
	err = util.MapToStruct(instance.Dependency.Network.Annotations, &awsnetwork, AnnotationPrefix)
	if err != nil {
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
	err = util.MapToStruct(instance.Annotations, &awsnetwork, AnnotationPrefix)
	if err != nil {
		return
	}

	status2, err := p.getInstance(awsnetwork, instance.Name)

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
	return
}

func (p *awsCloud) createInstance(clusterName string, instance *cluster.Instance) (status *cluster.InstanceStatus, err error) {
	options := InstanceOptions{}
	if instance.Labels != nil {
		err = util.MapToStruct(instance.Labels, &options, AnnotationPrefix)
		if err != nil {
			return
		}
	}

	awsnetwork := AWSNetwork{}
	err = util.MapToStruct(instance.Annotations, &awsnetwork, AnnotationPrefix)
	if err != nil {
		return
	}

	networkSpec := cluster.NetworkSpec{}
	err = util.MapToStruct(instance.Annotations, &networkSpec, AnnotationPrefix)
	if err != nil {
		return
	}

	eip := EIP{}
	privateIP := PrivateIP{}
	awsPrivateIP := (*string)(nil)
	ifSpecs := ([]*ec2.InstanceNetworkInterfaceSpecification)(nil)

	if options.CustomPrivateIP {
		err = util.MapToStruct(instance.Annotations, &privateIP, AnnotationPrefix)
		if err != nil {
			return
		}

		if privateIP.IP == "" {
			err = fmt.Errorf("Custom private IP is not provided.")
			return
		}

		awsPrivateIP = aws.String(privateIP.IP)
	}

	if options.UseEIP {
		err = util.MapToStruct(instance.Annotations, &eip, AnnotationPrefix)
		if err != nil {
			return
		}

		if eip.AllocationID == "" || eip.IP == "" {
			err = fmt.Errorf("EIP is not provided.")
			return
		}

		// Create if
		resp, err := p.ec2.CreateNetworkInterface(&ec2.CreateNetworkInterfaceInput{
			PrivateIpAddress: awsPrivateIP,
			SubnetId:         aws.String(networkSpec.Subnet),
		})
		if err != nil {
			return nil, fmt.Errorf("Error creating network interface: %+v", err)
		}
		nif := resp.NetworkInterface

		// Associate address
		_, err = p.ec2.AssociateAddress(&ec2.AssociateAddressInput{
			AllocationId:       aws.String(eip.AllocationID),
			NetworkInterfaceId: nif.NetworkInterfaceId,
		})
		if err != nil {
			return nil, fmt.Errorf("Error associating eip with network interface: %+v", err)
		}

		// Add to ifspecs
		ifSpecs = append(ifSpecs, &ec2.InstanceNetworkInterfaceSpecification{
			DeleteOnTermination: aws.Bool(true),
			DeviceIndex:         aws.Int64(0),
			NetworkInterfaceId:  nif.NetworkInterfaceId,
		})
	} else if options.CustomPrivateIP {
		// Add to ifspecs
		ifSpecs = append(ifSpecs, &ec2.InstanceNetworkInterfaceSpecification{
			AssociatePublicIpAddress: aws.Bool(true),
			PrivateIpAddress:         awsPrivateIP,
			DeleteOnTermination:      aws.Bool(true),
			DeviceIndex:              aws.Int64(0),
		})
	}

	// TODO: Get from config
	awsInstanceType := "t2.small"

	// TODO: Get userdata
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
	params := &ec2.RunInstancesInput{
		ImageId:      aws.String(image), // Required
		MaxCount:     aws.Int64(1),      // Required
		MinCount:     aws.Int64(1),      // Required
		ClientToken:  aws.String(util.RandNano()),
		InstanceType: aws.String(awsInstanceType),
		SubnetId:     aws.String(awsnetwork.Subnet),
		UserData:     aws.String(s),
		Placement: &ec2.Placement{
			AvailabilityZone: aws.String(networkSpec.Zone),
		},
		NetworkInterfaces: ifSpecs,
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
	return
}

func (p *awsCloud) EnsureInstanceDeleted(clusterName string, instance *cluster.Instance) (err error) {
	awsnetwork := AWSNetwork{}
	err = util.MapToStruct(instance.Annotations, &awsnetwork, AnnotationPrefix)
	if err != nil {
		return
	}

	status2, err := p.getInstance(awsnetwork, instance.Name)
	if err == ErrorNotFound {
		return nil
	}

	if err != nil {
		return
	}

	vpsId := status2.InstanceID
	if vpsId == "" {
		return fmt.Errorf("Unable to get instance id: %+v", status2)
	}

	return p.deleteInstance(vpsId)
}

func (p *awsCloud) deleteInstance(vpsId string) (err error) {
	params2 := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{ // Required
			aws.String(vpsId), // Required
		},
	}
	_, err = p.ec2.TerminateInstances(params2)
	if err != nil {
		if isNotExistError(err) {
			err = nil
		}
		return
	}

	return p.ec2.WaitUntilInstanceTerminated(&ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(vpsId)},
	})
}
