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
	"kubeup.com/archon/pkg/cluster"
)

type UCloudNetwork struct {
	VPC           string `k8s:"vpc-id"`
	Subnet        string `k8s:"subnet-id"`
	SecurityGroup string `k8s:"security-group-id"`
}

func (p *UCloud) EnsureNetwork(clusterName string, network *cluster.Network) (status *cluster.NetworkStatus, err error) {
	if network.Status.Phase != cluster.NetworkRunning {
		err = fmt.Errorf("UCloud doesn't support VPC apis. You need to create a running Network manually with VPC, Subnet and Security Group annotations")
	} else {
		status = &network.Status
	}

	return
}

func (p *UCloud) EnsureNetworkDeleted(clusterName string, network *cluster.Network) (err error) {
	glog.Warningf("UCloud doesn't support VPC apis. So you need to delete the network manually")
	return
}
