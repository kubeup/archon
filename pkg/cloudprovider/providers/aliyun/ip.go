package aliyun

import (
	"kubeup.com/archon/pkg/cluster"
)

func (p *aliyunCloud) EnsurePublicIP(clusterName string, instance *cluster.Instance) (status *cluster.InstanceStatus, err error) {
	return &instance.Status, nil
}

func (p *aliyunCloud) EnsurePublicIPDeleted(clusterName string, instance *cluster.Instance) error {
	return nil
}

func (p *aliyunCloud) EnsurePrivateIP(clusterName string, instance *cluster.Instance) (status *cluster.InstanceStatus, err error) {
	return &instance.Status, nil
}

func (p *aliyunCloud) EnsurePrivateIPDeleted(clusterName string, instance *cluster.Instance) error {
	return nil
}
