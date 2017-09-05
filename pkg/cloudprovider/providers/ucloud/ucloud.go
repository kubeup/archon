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
	"github.com/ucloud/ucloud-sdk-go/service/uhost"
	"github.com/ucloud/ucloud-sdk-go/service/unet"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/auth"
	"io"
	cp "k8s.io/kubernetes/pkg/cloudprovider"
	"k8s.io/kubernetes/pkg/controller"
	"kubeup.com/archon/pkg/cloudprovider"
	"kubeup.com/archon/pkg/cluster"
	"os"
)

const ProviderName = "ucloud"

// Used in k8s annotations and labels
var UCloudAnnotationPrefix = "ucloud." + cluster.AnnotationPrefix

// Node name as a tag on uhost instance
var NameKey = UCloudAnnotationPrefix + "name"

func init() {
	cloudprovider.RegisterCloudProvider(ProviderName, func(config io.Reader) (cloudprovider.Interface, error) {
		return newUCloud(config)
	})
}

type UCloud struct {
	ucloud UCloudInterface
}

func newUCloud(config io.Reader) (cloudprovider.Interface, error) {
	publicKey := os.Getenv("UCLOUD_PUBLIC_KEY")
	privateKey := os.Getenv("UCLOUD_PRIVATE_KEY")
	if publicKey == "" || privateKey == "" {
		return nil, fmt.Errorf("UCLOUD_PUBLIC_KEY or UCLOUD_PRIVATE_KEY is not set")
	}

	region := os.Getenv("UCLOUD_REGION")
	if region == "" {
		return nil, fmt.Errorf("UCLOUD_REGION is not specified")
	}
	project := os.Getenv("UCLOUD_PROJECT")
	if project == "" {
		return nil, fmt.Errorf("UCLOUD_PROJECT is not specified")
	}

	hostsvc := uhost.New(&ucloud.Config{
		Credentials: &auth.KeyPair{
			PublicKey:  publicKey,
			PrivateKey: privateKey,
		},
		Region:    region,
		ProjectID: project,
	})

	netsvc := unet.New(&ucloud.Config{
		Credentials: &auth.KeyPair{
			PublicKey:  publicKey,
			PrivateKey: privateKey,
		},
		Region:    region,
		ProjectID: project,
	})

	return &UCloud{
		ucloud: &ucloudWrapper{hostsvc, netsvc},
	}, nil
}

func (p *UCloud) Initialize(clientBuilder controller.ControllerClientBuilder) {}

func (p *UCloud) ProviderName() string {
	return ProviderName
}

func (p *UCloud) Archon() (cloudprovider.ArchonInterface, bool) {
	return p, true
}

func (p *UCloud) Instances() (cp.Instances, bool) {
	return nil, false
}

func (p *UCloud) Clusters() (cp.Clusters, bool) {
	return nil, false
}

func (p *UCloud) LoadBalancer() (cp.LoadBalancer, bool) {
	return nil, false
}

func (p *UCloud) Routes() (cp.Routes, bool) {
	return nil, false
}

func (p *UCloud) Zones() (cp.Zones, bool) {
	return nil, false
}

func (p *UCloud) ScrubDNS(nameservers, searches []string) (nsOut, srchOut []string) {
	return nil, nil
}

func (p *UCloud) PrivateIP() (cloudprovider.PrivateIPInterface, bool) {
	return p, true
}

func (p *UCloud) PublicIP() (cloudprovider.PublicIPInterface, bool) {
	return p, true
}
