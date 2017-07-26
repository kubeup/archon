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

package cloudprovider

import (
	"k8s.io/kubernetes/pkg/cloudprovider"
	"kubeup.com/archon/pkg/cluster"
)

type PublicIPInterface interface {
	EnsurePublicIP(clusterName string, instance *cluster.Instance) (status *cluster.InstanceStatus, err error)
	EnsurePublicIPDeleted(clusterName string, instance *cluster.Instance) error
}

type PrivateIPInterface interface {
	EnsurePrivateIP(clusterName string, instance *cluster.Instance) (status *cluster.InstanceStatus, err error)
	EnsurePrivateIPDeleted(clusterName string, instance *cluster.Instance) error
}

type ArchonInterface interface {
	EnsureNetwork(clusterName string, network *cluster.Network) (status *cluster.NetworkStatus, err error)
	EnsureNetworkDeleted(clusterName string, network *cluster.Network) error
	//AddNetworkAnnotation(clusterName string, instance *cluster.Instance, network *cluster.Network) error

	ListInstances(clusterName string, network *cluster.Network, selector map[string]string) (names []string, instances []*cluster.InstanceStatus, err error)
	GetInstance(clusterName string, instance *cluster.Instance) (status *cluster.InstanceStatus, err error)
	EnsureInstance(clusterName string, instance *cluster.Instance) (status *cluster.InstanceStatus, err error)
	EnsureInstanceDeleted(clusterName string, instance *cluster.Instance) error

	PublicIP() (PublicIPInterface, bool)
	PrivateIP() (PrivateIPInterface, bool)
}

type Interface interface {
	cloudprovider.Interface
	Archon() (ArchonInterface, bool)
}
