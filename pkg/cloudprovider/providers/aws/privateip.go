package aws

import (
	"fmt"
	aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"kubeup.com/archon/pkg/cloudprovider"
	"kubeup.com/archon/pkg/cluster"
	"kubeup.com/archon/pkg/util"
)

type PrivateIP struct {
	NetworkInterfaceID string `k8s:"private-ip-nif-id"`
}

func (p *awsCloud) PrivateIP() (cloudprovider.PrivateIPInterface, bool) {
	return p, true
}

func (p *awsCloud) EnsurePrivateIP(clusterName string, instance *cluster.Instance) (status *cluster.InstanceStatus, err error) {
	options := cluster.InstanceOptions{}
	pip := PrivateIP{}

	if instance.Labels != nil {
		err = util.MapToStruct(instance.Labels, &options, cluster.AnnotationPrefix)
		if err != nil {
			return
		}
	}

	if options.PreallocatePrivateIP && instance.Annotations != nil {
		err = util.MapToStruct(instance.Annotations, &pip, AWSAnnotationPrefix)
		if err != nil {
			return
		}
	}

	if pip.NetworkInterfaceID != "" {
		status = &instance.Status
		return
	}

	// Create a network interface
	awsnetwork := AWSNetwork{}
	if instance.Annotations != nil {
		err = util.MapToStruct(instance.Annotations, &awsnetwork, AWSAnnotationPrefix)
		if err != nil {
			err = fmt.Errorf("Network is not ready. Can't allocate private ip: %s", err.Error())
			return
		}
	}

	if awsnetwork.Subnet == "" {
		err = fmt.Errorf("Network is not ready. Can't allocate private ip")
		return
	}

	privateIP := (*string)(nil)
	if instance.Status.PrivateIP != "" {
		privateIP = aws.String(instance.Status.PrivateIP)
	}

	resp, err := p.ec2.CreateNetworkInterface(&ec2.CreateNetworkInterfaceInput{
		PrivateIpAddress: privateIP,
		SubnetId:         aws.String(awsnetwork.Subnet),
	})
	if err != nil {
		return nil, fmt.Errorf("Error creating network interface with private ip: %+v", err)
	}

	pip.NetworkInterfaceID = destring(resp.NetworkInterface.NetworkInterfaceId)

	if instance.Annotations == nil {
		instance.Annotations = make(map[string]string)
	}

	err = util.StructToMap(pip, instance.Annotations, AWSAnnotationPrefix)
	if err != nil {
		err = fmt.Errorf("Error allocating EIP: %s", err.Error())
		return
	}

	status = &instance.Status
	status.PrivateIP = destring(resp.NetworkInterface.PrivateIpAddress)

	return
}

func (p *awsCloud) EnsurePrivateIPDeleted(clusterName string, instance *cluster.Instance) (err error) {
	if instance.Annotations == nil {
		return
	}

	pip := PrivateIP{}
	err = util.MapToStruct(instance.Annotations, &pip, AWSAnnotationPrefix)
	if err != nil {
		return
	}

	if pip.NetworkInterfaceID != "" {
		_, err = p.ec2.DeleteNetworkInterface(&ec2.DeleteNetworkInterfaceInput{
			NetworkInterfaceId: aws.String(pip.NetworkInterfaceID),
		})

		if err != nil && !isNotExistError(err) {
			err = fmt.Errorf("Error releasing private ip: %s", err.Error())
			return
		}

		pip.NetworkInterfaceID = ""

		err = util.StructToMap(pip, instance.Annotations, AWSAnnotationPrefix)
		if err != nil {
			err = fmt.Errorf("Error releasing private ip: %s", err.Error())
			return
		}

		instance.Status.PrivateIP = ""
	}

	return nil
}
