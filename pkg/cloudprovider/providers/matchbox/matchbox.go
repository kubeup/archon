/*
Copyright 2016 The Kubernetes Authors.
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

package matchbox

import (
	"context"
	"fmt"
	"io"
	"os"

	archoncloudprovider "kubeup.com/archon/pkg/cloudprovider"
	"kubeup.com/archon/pkg/cluster"
	"kubeup.com/archon/pkg/userdata"

	"github.com/coreos/matchbox/matchbox/client"
	pb "github.com/coreos/matchbox/matchbox/server/serverpb"
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
	"github.com/coreos/matchbox/matchbox/tlsutil"
	"github.com/golang/glog"
	"k8s.io/kubernetes/pkg/cloudprovider/providers/fake"
)

var (
	ErrorNotFound            = fmt.Errorf("Instance is not found")
	MatchboxAnnotationPrefix = "matchbox." + cluster.AnnotationPrefix
)

const ProviderName = "matchbox"

type FakeInstance struct {
	name     string
	spec     cluster.InstanceSpec
	userdata []byte
}

// MatchboxCloud is a test-double implementation of Interface, LoadBalancer, Instances, and Routes. It is useful for testing.
type MatchboxCloud struct {
	fake.FakeCloud
	FakeInstances map[string]FakeInstance
	client        *client.Client
	endpoint      string
}

var (
	_ archoncloudprovider.Interface = &MatchboxCloud{}
)

func newMatchboxCloud(config io.Reader) (archoncloudprovider.Interface, error) {
	httpEndpoint := os.Getenv("MATCHBOX_HTTP_ENDPOINT")
	rpcEndpoint := os.Getenv("MATCHBOX_RPC_ENDPOINT")
	caFile := os.Getenv("MATCHBOX_CA_FILE")
	certFile := os.Getenv("MATCHBOX_CERT_FILE")
	keyFile := os.Getenv("MATCHBOX_KEY_FILE")
	if httpEndpoint == "" {
		httpEndpoint = "matchbox.foo:8080"
	}
	if rpcEndpoint == "" {
		rpcEndpoint = "127.0.0.1:8081"
	}
	if caFile == "" {
		caFile = "/etc/matchbox/ca.crt"
	}
	if certFile == "" {
		certFile = "/etc/matchbox/client.crt"
	}
	if keyFile == "" {
		keyFile = "/etc/matchbox/client.key"
	}
	tlsinfo := &tlsutil.TLSInfo{
		CAFile:   caFile,
		CertFile: certFile,
		KeyFile:  keyFile,
	}
	tlscfg, err := tlsinfo.ClientConfig()
	if err != nil {
		return nil, err
	}
	cfg := &client.Config{
		Endpoints: []string{rpcEndpoint},
		TLS:       tlscfg,
	}

	// gRPC client
	client, err := client.New(cfg)
	if err != nil {
		return nil, err
	}
	return &MatchboxCloud{
		client:   client,
		endpoint: httpEndpoint,
	}, nil
}

func (f *MatchboxCloud) addCall(desc string) {
	f.FakeCloud.Calls = append(f.FakeCloud.Calls, desc)
}

func (f *MatchboxCloud) Archon() (archoncloudprovider.ArchonInterface, bool) {
	return f, true
}

func (f *MatchboxCloud) PrivateIP() (archoncloudprovider.PrivateIPInterface, bool) {
	return nil, false
}

func (f *MatchboxCloud) PublicIP() (archoncloudprovider.PublicIPInterface, bool) {
	return nil, false
}

func (f *MatchboxCloud) AddNetworkAnnotation(clustername string, instance *cluster.Instance, network *cluster.Network) error {
	return nil
}

func (f *MatchboxCloud) EnsureNetwork(clusterName string, network *cluster.Network) (status *cluster.NetworkStatus, err error) {
	return &cluster.NetworkStatus{Phase: cluster.NetworkRunning}, nil
}

func (f *MatchboxCloud) EnsureNetworkDeleted(clusterName string, network *cluster.Network) (err error) {
	return nil
}

// GetLoadBalancer is a stub implementation of LoadBalancer.GetLoadBalancer.
func (f *MatchboxCloud) GetInstance(clusterName string, instance *cluster.Instance) (*cluster.InstanceStatus, error) {
	if f.FakeInstances == nil {
		return nil, ErrorNotFound
	}
	if _, ok := f.FakeInstances[instance.Name]; !ok {
		return nil, ErrorNotFound
	}
	status := &cluster.InstanceStatus{}
	status.Phase = cluster.InstanceRunning

	return status, f.Err
}

// EnsureLoadBalancer is a test-spy implementation of LoadBalancer.EnsureLoadBalancer.
// It adds an entry "create" into the internal method call record.
func (f *MatchboxCloud) EnsureInstance(clusterName string, instance *cluster.Instance) (*cluster.InstanceStatus, error) {
	f.addCall("create")
	if f.FakeInstances == nil {
		f.FakeInstances = make(map[string]FakeInstance)
	}

	name := instance.Name
	spec := instance.Spec

	err := f.ensureGroup(clusterName, instance)
	if err != nil {
		return nil, err
	}

	err = f.ensureProfile(clusterName, instance)
	if err != nil {
		return nil, err
	}

	err = f.ensureIgnition(clusterName, instance)
	if err != nil {
		return nil, err
	}

	f.FakeInstances[name] = FakeInstance{name, spec, []byte{}}

	status := &cluster.InstanceStatus{}
	status.Phase = cluster.InstanceRunning

	return status, f.Err
}

func (f *MatchboxCloud) ensureGroup(clusterName string, instance *cluster.Instance) error {
	mac := instance.Annotations[MatchboxAnnotationPrefix+"mac"]
	if mac == "" {
		return fmt.Errorf("No mac address on instance")
	}
	group := &storagepb.Group{
		Id:      instance.Name,
		Name:    instance.Name,
		Profile: instance.Name,
		Selector: map[string]string{
			"mac": mac,
			"os":  "installed",
		},
	}
	req := &pb.GroupPutRequest{Group: group}
	_, err := f.client.Groups.GroupPut(context.TODO(), req)
	return err
}

func (f *MatchboxCloud) ensureProfile(clusterName string, instance *cluster.Instance) error {
	profile := &storagepb.Profile{
		Id:         instance.Name,
		Name:       instance.Name,
		IgnitionId: instance.Name,
		Boot: &storagepb.NetBoot{
			Kernel: fmt.Sprintf("/assets/coreos/%s/coreos_production_pxe.vmlinuz", instance.Spec.Image),
			Initrd: []string{fmt.Sprintf("/assets/coreos/%s/coreos_production_pxe_image.cpio.gz", instance.Spec.Image)},
			Args: []string{
				"root=/dev/sda1",
				fmt.Sprintf("coreos.config.url=http://%s/ignition?uuid=${uuid}&mac=${mac:hexhyp}", f.endpoint),
				"coreos.first_boot=yes",
				"net.ifnames=0",
				"console=tty0",
				"console=ttyS0",
				"coreos.autologin",
			},
		},
	}
	req := &pb.ProfilePutRequest{Profile: profile}
	_, err := f.client.Profiles.ProfilePut(context.TODO(), req)
	return err
}

func (f *MatchboxCloud) ensureIgnition(clusterName string, instance *cluster.Instance) error {
	userdata, err := userdata.Generate(instance)
	if err != nil {
		glog.Errorf("Generate userdata error: %+v", err)
		return err
	}

	glog.Infof("Create instance with userdata:\n%s", string(userdata))
	req := &pb.IgnitionPutRequest{Name: instance.Name, Config: userdata}
	_, err = f.client.Ignition.IgnitionPut(context.TODO(), req)
	return err
}

// EnsureLoadBalancerDeleted is a test-spy implementation of LoadBalancer.EnsureLoadBalancerDeleted.
// It adds an entry "delete" into the internal method call record.
func (f *MatchboxCloud) EnsureInstanceDeleted(clusterName string, instance *cluster.Instance) error {
	f.addCall("delete")
	return f.Err
}

// List is a test-spy implementation of Instances.List.
// It adds an entry "list" into the internal method call record.
func (f *MatchboxCloud) ListInstances(clusterName string, network *cluster.Network, selector map[string]string) ([]string, []*cluster.InstanceStatus, error) {
	f.addCall("list")
	result := make([]string, 0)
	instances := make([]*cluster.InstanceStatus, 0)
	return result, instances, f.Err
}

func init() {
	archoncloudprovider.RegisterCloudProvider(ProviderName, func(config io.Reader) (archoncloudprovider.Interface, error) {
		return newMatchboxCloud(config)
	})
}
