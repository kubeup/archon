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
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	archoncloudprovider "kubeup.com/archon/pkg/cloudprovider"
	"kubeup.com/archon/pkg/cluster"
	"kubeup.com/archon/pkg/userdata"
	"kubeup.com/archon/pkg/util"

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

type MatchboxInstanceOptions struct {
	Mac           string `k8s:"mac"`
	Selector      string `k8s:"selector"`
	UseProfile    string `k8s:"use-profile"`
	UseIgnition   string `k8s:"use-ignition"`
	ExtraBootArgs string `k8s:"extra-boot-args"`
}

type FakeInstance struct {
	name   string
	spec   cluster.InstanceSpec
	status *cluster.InstanceStatus
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
	if i, ok := f.FakeInstances[instance.Name]; !ok {
		return nil, ErrorNotFound
	} else {
		return i.status, f.Err
	}
	return nil, nil
}

// EnsureLoadBalancer is a test-spy implementation of LoadBalancer.EnsureLoadBalancer.
// It adds an entry "create" into the internal method call record.
func (f *MatchboxCloud) EnsureInstance(clusterName string, instance *cluster.Instance) (*cluster.InstanceStatus, error) {
	f.addCall("create")
	if f.FakeInstances == nil {
		f.FakeInstances = make(map[string]FakeInstance)
	}

	options := MatchboxInstanceOptions{}
	err := util.MapToStruct(instance.Annotations, &options, MatchboxAnnotationPrefix)
	if err != nil {
		return nil, fmt.Errorf("Can't get instance options: %s", err.Error())
	}

	if i := instance.Dependency.ReservedInstance; i.Name != "" {
		if i.Spec.InstanceID != "" {
			instance.Status.InstanceID = i.Spec.InstanceID
		}
		for _, c := range i.Spec.Configs {
			if c.Name == "spec" {
				if options.Mac == "" {
					options.Mac = c.Data["mac"]
				}
				if instance.Status.PrivateIP == "" {
					instance.Status.PrivateIP = c.Data["private-ip"]
				}
			}
		}
	}

	err = f.ensureGroup(clusterName, instance, &options)
	if err != nil {
		return nil, err
	}

	err = f.ensureProfile(clusterName, instance, &options)
	if err != nil {
		return nil, err
	}

	err = f.ensureIgnition(clusterName, instance, &options)
	if err != nil {
		return nil, err
	}

	instance.Status.Phase = cluster.InstanceRunning

	name := instance.Name
	spec := instance.Spec
	status := instance.Status

	f.FakeInstances[name] = FakeInstance{name, spec, &status}

	return &status, f.Err
}

func (f *MatchboxCloud) ensureGroup(clusterName string, instance *cluster.Instance, options *MatchboxInstanceOptions) error {
	if options.Mac == "" {
		return fmt.Errorf("No mac address on instance")
	}
	groupID := instance.Name
	if instance.Status.InstanceID != "" {
		groupID = instance.Status.InstanceID
	}
	selector := map[string]string{}
	if options.Selector != "" {
		err := json.Unmarshal([]byte(options.Selector), &selector)
		if err != nil {
			return fmt.Errorf("Unmarshal selector error: %+v", err)
		}
	}
	profileID := groupID
	if options.UseProfile != "" {
		profileID = options.UseProfile
	}
	group := &storagepb.Group{
		Id:       groupID,
		Name:     instance.Name,
		Profile:  profileID,
		Selector: selector,
	}
	// we use FF:FF:FF:FF:FF:FF as a special mac address for default group
	// in the default group, we don't set the mac selector so all servers will got a match
	if options.Mac != "FF:FF:FF:FF:FF:FF" {
		group.Selector["mac"] = options.Mac
	}
	for _, c := range instance.Spec.Configs {
		if c.Name == "meta" {
			meta, err := json.Marshal(c.Data)
			if err != nil {
				return fmt.Errorf("Marshal metadata error: %+v", err)
			}
			group.Metadata = meta
		}
	}
	req := &pb.GroupPutRequest{Group: group}
	_, err := f.client.Groups.GroupPut(context.TODO(), req)
	return err
}

func (f *MatchboxCloud) ensureProfile(clusterName string, instance *cluster.Instance, options *MatchboxInstanceOptions) error {
	if options.UseProfile != "" {
		return nil
	}
	profileID := instance.Name
	if instance.Status.InstanceID != "" {
		profileID = instance.Status.InstanceID
	}
	ignitionID := profileID
	if options.UseIgnition != "" {
		ignitionID = options.UseIgnition
	}
	profile := &storagepb.Profile{
		Id:         profileID,
		Name:       instance.Name,
		IgnitionId: ignitionID,
		Boot: &storagepb.NetBoot{
			Kernel: fmt.Sprintf("/assets/coreos/%s/coreos_production_pxe.vmlinuz", instance.Spec.Image),
			Initrd: []string{fmt.Sprintf("/assets/coreos/%s/coreos_production_pxe_image.cpio.gz", instance.Spec.Image)},
			Args: []string{
				fmt.Sprintf("coreos.config.url=http://%s/ignition?uuid=${uuid}&mac=${mac:hexhyp}", f.endpoint),
				"coreos.first_boot=yes",
				"net.ifnames=0",
				"console=tty0",
				"console=ttyS0",
			},
		},
	}
	if options.ExtraBootArgs != "" {
		profile.Boot.Args = append(profile.Boot.Args, strings.Split(options.ExtraBootArgs, ",")...)
	}
	req := &pb.ProfilePutRequest{Profile: profile}
	_, err := f.client.Profiles.ProfilePut(context.TODO(), req)
	return err
}

func (f *MatchboxCloud) ensureIgnition(clusterName string, instance *cluster.Instance, options *MatchboxInstanceOptions) error {
	if options.UseIgnition != "" {
		return nil
	}
	userdata, err := userdata.Generate(instance)
	if err != nil {
		glog.Errorf("Generate userdata error: %+v", err)
		return err
	}
	ignitionName := instance.Name
	if instance.Status.InstanceID != "" {
		ignitionName = instance.Status.InstanceID
	}

	req := &pb.IgnitionPutRequest{Name: ignitionName, Config: userdata}
	_, err = f.client.Ignition.IgnitionPut(context.TODO(), req)
	return err
}

// EnsureLoadBalancerDeleted is a test-spy implementation of LoadBalancer.EnsureLoadBalancerDeleted.
// It adds an entry "delete" into the internal method call record.
func (f *MatchboxCloud) EnsureInstanceDeleted(clusterName string, instance *cluster.Instance) error {
	f.addCall("delete")
	if f.FakeInstances == nil {
		return f.Err
	}
	if _, ok := f.FakeInstances[instance.Name]; ok {
		delete(f.FakeInstances, instance.Name)
	}
	return f.Err
}

// List is a test-spy implementation of Instances.List.
// It adds an entry "list" into the internal method call record.
func (f *MatchboxCloud) ListInstances(clusterName string, network *cluster.Network, selector map[string]string) ([]string, []*cluster.InstanceStatus, error) {
	f.addCall("list")
	result := make([]string, 0)
	instances := make([]*cluster.InstanceStatus, 0)
	for k, v := range f.FakeInstances {
		result = append(result, k)
		instances = append(instances, v.status)
	}
	return result, instances, f.Err
}

func init() {
	archoncloudprovider.RegisterCloudProvider(ProviderName, func(config io.Reader) (archoncloudprovider.Interface, error) {
		return newMatchboxCloud(config)
	})
}
