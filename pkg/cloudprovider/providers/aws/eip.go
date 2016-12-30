package aws

import (
	"fmt"
	aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"kubeup.com/archon/pkg/cloudprovider"
	"kubeup.com/archon/pkg/cluster"
	"kubeup.com/archon/pkg/util"
)

type EIP struct {
	AllocationID string `k8s:"eip-allocation-id"`
}

func (p *awsCloud) PublicIP() (cloudprovider.PublicIPInterface, bool) {
	return p, true
}

func (p *awsCloud) EnsurePublicIP(clusterName string, instance *cluster.Instance) (status *cluster.InstanceStatus, err error) {
	options := cluster.InstanceOptions{}
	eip := EIP{}

	if instance.Labels != nil {
		err = util.MapToStruct(instance.Labels, &options, cluster.AnnotationPrefix)
		if err != nil {
			return
		}
	}

	if options.PreallocatePublicIP && instance.Annotations != nil {
		err = util.MapToStruct(instance.Annotations, &eip, AWSAnnotationPrefix)
		if err != nil {
			return
		}
	}

	if eip.AllocationID != "" {
		status = &instance.Status
		return
	}

	resp, err := p.ec2.AllocateAddress(&ec2.AllocateAddressInput{
		Domain: aws.String("vpc"),
	})
	if err != nil {
		err = fmt.Errorf("Error allocating EIP: %s", err.Error())
		return
	}

	eip.AllocationID = destring(resp.AllocationId)
	if instance.Annotations == nil {
		instance.Annotations = make(map[string]string)
	}

	err = util.StructToMap(eip, instance.Annotations, AWSAnnotationPrefix)
	if err != nil {
		err = fmt.Errorf("Error allocating EIP: %s", err.Error())
		return
	}

	status = &instance.Status
	status.PublicIP = destring(resp.PublicIp)

	return
}

func (p *awsCloud) EnsurePublicIPDeleted(clusterName string, instance *cluster.Instance) (err error) {
	if instance.Annotations == nil {
		return
	}

	eip := EIP{}
	err = util.MapToStruct(instance.Annotations, &eip, AWSAnnotationPrefix)
	if err != nil {
		return
	}

	if eip.AllocationID != "" {
		_, err = p.ec2.ReleaseAddress(&ec2.ReleaseAddressInput{
			AllocationId: aws.String(eip.AllocationID),
		})

		if err != nil {
			err = fmt.Errorf("Error releasing EIP: %s", err.Error())
			return
		}

		eip.AllocationID = ""

		if instance.Annotations == nil {
			instance.Annotations = make(map[string]string)
		}

		err = util.StructToMap(eip, instance.Annotations, AWSAnnotationPrefix)
		if err != nil {
			err = fmt.Errorf("Error releasing EIP: %s", err.Error())
			return
		}

		instance.Status.PublicIP = ""
	}

	return nil
}
