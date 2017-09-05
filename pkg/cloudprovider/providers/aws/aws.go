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
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/golang/glog"
	"io"
	cp "k8s.io/kubernetes/pkg/cloudprovider"
	"k8s.io/kubernetes/pkg/controller"
	"kubeup.com/archon/pkg/cloudprovider"
	"kubeup.com/archon/pkg/cluster"
)

const ProviderName = "aws"

// Used in k8s annotations and labels
var AWSAnnotationPrefix = "aws." + cluster.AnnotationPrefix

// Node name as a tag on aws instance
var NameKey = AWSAnnotationPrefix + "name"

func init() {
	cloudprovider.RegisterCloudProvider(ProviderName, func(config io.Reader) (cloudprovider.Interface, error) {
		creds := credentials.NewChainCredentials(
			[]credentials.Provider{
				&credentials.EnvProvider{},
				&credentials.SharedCredentialsProvider{},
			})
		aws := newAWSSDKProvider(creds)
		return newAWSCloud(config, aws)
	})
}

type awsCloud struct {
	ec2 EC2
	s3  S3
	iam IAM

	region string
}

func newAWSCloud(config io.Reader, service *awsSDKProvider) (cloudprovider.Interface, error) {
	metadata, err := service.Metadata()
	if err != nil {
		return nil, fmt.Errorf("error creating AWS metadata client: %v", err)
	}

	cfg, err := readAWSCloudConfig(config, metadata)
	if err != nil {
		return nil, fmt.Errorf("unable to read AWS cloud provider config file: %v", err)
	}

	zone := cfg.Global.Zone
	if len(zone) <= 1 {
		return nil, fmt.Errorf("invalid AWS zone in config file: %s", zone)
	}

	regionName, err := azToRegion(zone)
	if err != nil {
		return nil, err
	}

	if !cfg.Global.DisableStrictZoneCheck {
		valid := isRegionValid(regionName)
		if !valid {
			return nil, fmt.Errorf("not a valid AWS zone (unknown region): %s", zone)
		}
	} else {
		glog.Warningf("Strict AWS zone checking is disabled.  Proceeding with zone: %s", zone)
	}

	ec2, err := service.Compute(regionName)
	if err != nil {
		return nil, fmt.Errorf("error creating AWS EC2 client: %v", err)
	}

	iam, err := service.IAM(regionName)
	if err != nil {
		return nil, fmt.Errorf("error creating AWS IAM client: %v", err)
	}

	s3, err := service.S3(regionName)
	if err != nil {
		return nil, fmt.Errorf("error creating AWS S3 client: %v", err)
	}

	return &awsCloud{
		ec2:    ec2,
		iam:    iam,
		s3:     s3,
		region: regionName,
	}, nil
}

func (p *awsCloud) Initialize(clientBuilder controller.ControllerClientBuilder) {}

func (p *awsCloud) ProviderName() string {
	return ProviderName
}

func (p *awsCloud) Archon() (cloudprovider.ArchonInterface, bool) {
	return p, true
}

func (p *awsCloud) Instances() (cp.Instances, bool) {
	return nil, false
}

func (p *awsCloud) Clusters() (cp.Clusters, bool) {
	return nil, false
}

func (p *awsCloud) LoadBalancer() (cp.LoadBalancer, bool) {
	return nil, false
}

func (p *awsCloud) Routes() (cp.Routes, bool) {
	return nil, false
}

func (p *awsCloud) Zones() (cp.Zones, bool) {
	return nil, false
}

func (p *awsCloud) ScrubDNS(nameservers, searches []string) (nsOut, srchOut []string) {
	return nil, nil
}
