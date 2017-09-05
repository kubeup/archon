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
	"github.com/denverdino/aliyungo/ecs"
	"io"
	cp "k8s.io/kubernetes/pkg/cloudprovider"
	"k8s.io/kubernetes/pkg/controller"
	"kubeup.com/archon/pkg/cloudprovider"
	"kubeup.com/archon/pkg/cluster"
	"os"
)

const ProviderName = "aliyun"

// Used in k8s annotations and labels
var AliyunAnnotationPrefix = "aliyun." + cluster.AnnotationPrefix

// Node name as a tag on aliyun instance
var NameKey = AliyunAnnotationPrefix + "name"

func init() {
	cloudprovider.RegisterCloudProvider(ProviderName, func(config io.Reader) (cloudprovider.Interface, error) {
		return newAliyunCloud(config)
	})
}

type aliyunCloud struct {
	ecs ECS
}

func newAliyunCloud(config io.Reader) (cloudprovider.Interface, error) {
	accessKey := os.Getenv("ALIYUN_ACCESS_KEY")
	accessSecret := os.Getenv("ALIYUN_ACCESS_KEY_SECRET")
	if accessKey == "" || accessSecret == "" {
		return nil, fmt.Errorf("ALIYUN_ACCESS_KEY or ALIYUN_ACCESS_KEY_SECRET is not set")
	}

	debug := os.Getenv("ALIYUN_DEBUG")
	ecsClient := ecs.NewClient(accessKey, accessSecret)
	if debug == "true" {
		ecsClient.SetDebug(true)
	}

	return &aliyunCloud{
		ecs: &aliyunECS{ecsClient},
	}, nil
}

func (p *aliyunCloud) Initialize(clientBuilder controller.ControllerClientBuilder) {}

func (p *aliyunCloud) ProviderName() string {
	return ProviderName
}

func (p *aliyunCloud) Archon() (cloudprovider.ArchonInterface, bool) {
	return p, true
}

func (p *aliyunCloud) Instances() (cp.Instances, bool) {
	return nil, false
}

func (p *aliyunCloud) Clusters() (cp.Clusters, bool) {
	return nil, false
}

func (p *aliyunCloud) LoadBalancer() (cp.LoadBalancer, bool) {
	return nil, false
}

func (p *aliyunCloud) Routes() (cp.Routes, bool) {
	return nil, false
}

func (p *aliyunCloud) Zones() (cp.Zones, bool) {
	return nil, false
}

func (p *aliyunCloud) ScrubDNS(nameservers, searches []string) (nsOut, srchOut []string) {
	return nil, nil
}

func (p *aliyunCloud) PrivateIP() (cloudprovider.PrivateIPInterface, bool) {
	return p, true
}

func (p *aliyunCloud) PublicIP() (cloudprovider.PublicIPInterface, bool) {
	return p, true
}
