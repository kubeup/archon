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
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/golang/glog"
	cloudinit "github.com/tryk8s/ssh-cloudinit/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/kubernetes/pkg/api/v1"
	"kubeup.com/archon/pkg/cluster"
	"kubeup.com/archon/pkg/userdata"
	"kubeup.com/archon/pkg/util"
	"time"
)

var (
	ErrorNotFound     = fmt.Errorf("Instance is not found")
	CloudInitInterval = 5 * time.Second
	CloudInitTimeout  = 3 * time.Minute
	// Aliyun doesn't allow tag keys starting with "aliyun"
	InitializedTagKey = "instance." + AliyunAnnotationPrefix + "initialized"
	SSHUsername       = "root"
	SSHSystem         = "coreos"
	StateMap          = map[ecs.InstanceStatus]cluster.InstancePhase{
		ecs.Pending:  cluster.InstancePending,
		ecs.Starting: cluster.InstancePending,
		ecs.Running:  cluster.InstanceRunning,
		ecs.Stopping: cluster.InstanceFailed,
		ecs.Stopped:  cluster.InstanceFailed,
		ecs.Deleted:  cluster.InstanceFailed,
	}
)

type AliyunInstanceOptions struct {
	InternetMaxBandwidthIn  int    `k8s:"internet-max-bandwidth-in"`
	InternetMaxBandwidthOut int    `k8s:"internet-max-bandwidth-out"`
	SystemDiskSize          int    `k8s:"system-disk-size"`
	SystemDiskType          string `k8s:"system-disk-type"`
}

type AliyunInstanceInitialized struct {
	Initialized bool `k8s:"instance-initialized"`
}

func instanceToStatus(i ecs.InstanceAttributesType) *cluster.InstanceStatus {
	phase, ok := StateMap[i.Status]
	if !ok {
		glog.Warningf("Unknown instance state: %+v", i.Status)
		phase = cluster.InstanceUnknown
	}

	// If the instance is not marked as initialized and its vps is in stopped state,
	// return Pending
	initialized := false
	for _, tag := range i.Tags.Tag {
		if tag.TagKey == InitializedTagKey && tag.TagValue == "true" {
			initialized = true
		}
	}
	if i.Status == ecs.Stopped && initialized == false {
		phase = cluster.InstancePending
	}

	publicIP := i.EipAddress.IpAddress
	if publicIP == "" {
		publicIP = firstIP(i.PublicIpAddress)
	}

	return &cluster.InstanceStatus{
		Phase:             phase,
		PrivateIP:         firstIP(i.VpcAttributes.PrivateIpAddress),
		PublicIP:          publicIP,
		InstanceID:        i.InstanceId,
		CreationTimestamp: metav1.NewTime(time.Time(i.CreationTime)),
	}
}

func (p *aliyunCloud) ListInstances(clusterName string, network *cluster.Network, selector map[string]string) (names []string, statuses []*cluster.InstanceStatus, err error) {
	an := AliyunNetwork{}
	err = util.MapToStruct(network.Annotations, &an, AliyunAnnotationPrefix)
	if err != nil || an.VPC == "" {
		err = fmt.Errorf("Network is not ready. Can't list instances: %v", err)
		return
	}

	instances, err := p.ecs.DescribeInstances(&ecs.DescribeInstancesArgs{
		RegionId: common.Region(network.Spec.Region),
		VpcId:    an.VPC,
		Tag:      selector,
	})

	if err != nil {
		err = aliyunSafeError(err)
		return
	}

	for _, instance := range instances {
		names = append(names, instance.InstanceName)
		statuses = append(statuses, instanceToStatus(instance))
	}

	return
}

func (p *aliyunCloud) GetInstance(clusterName string, instance *cluster.Instance) (status *cluster.InstanceStatus, err error) {
	if instance.Dependency.Network.Spec.Region == "" {
		err = fmt.Errorf("Network is not ready. Can't get instance")
	}

	return p.getInstance(instance.Dependency.Network.Spec.Region, instance.Status.InstanceID)
}

func (p *aliyunCloud) getInstance(region string, instanceID string) (status *cluster.InstanceStatus, err error) {
	if instanceID == "" {
		return nil, ErrorNotFound
	}

	instance, err := p.ecs.DescribeInstanceAttribute(instanceID)

	if err != nil {
		if isNotFound(err) {
			err = ErrorNotFound
		} else {
			err = aliyunSafeError(err)
		}

		return
	}

	return instanceToStatus(*instance), nil
}

func (p *aliyunCloud) EnsureInstance(clusterName string, instance *cluster.Instance) (status *cluster.InstanceStatus, err error) {
	an := AliyunNetwork{}
	err = util.MapToStruct(instance.Dependency.Network.Annotations, &an, AliyunAnnotationPrefix)
	if err != nil || an.VSwitch == "" {
		err = fmt.Errorf("Network is not ready. Can't create instance: %v", err)
		return
	}

	options := cluster.InstanceOptions{}
	err = util.MapToStruct(instance.Labels, &options, cluster.AnnotationPrefix)
	if err != nil {
		err = fmt.Errorf("Can't get instance options: %s", err.Error())
		return
	}

	if instance.Status.InstanceID == "" {
		if options.UseInstanceID != "" {
			networkSpec := instance.Dependency.Network.Spec
			status, err = p.getInstance(networkSpec.Region, options.UseInstanceID)
		} else {
			status, err = p.createInstance(clusterName, instance)
		}

		if err != nil {
			return
		}

		instance.Status = *status
	}

	ai := AliyunInstanceInitialized{}
	err = util.MapToStruct(instance.Annotations, &ai, AliyunAnnotationPrefix)
	if err != nil && instance.Annotations != nil {
		err = fmt.Errorf("Can't tell if the instance is initialized: %v", err)
		return
	}

	if ai.Initialized == true {
		return &instance.Status, nil
	}

	return p.initializeInstance(clusterName, instance)
}

func (p *aliyunCloud) resetInstance(instance *cluster.Instance) (err error) {
	options := cluster.InstanceOptions{}
	err = util.MapToStruct(instance.Labels, &options, cluster.AnnotationPrefix)
	if err != nil {
		err = fmt.Errorf("Can't get instance options: %s", err.Error())
		return
	}

	networkSpec := instance.Dependency.Network.Spec

	instanceID := instance.Status.InstanceID
	if instanceID == "" {
		instanceID = options.UseInstanceID
	}

	// Ignore err in case the instance is already stopped
	err = p.ecs.StopInstance(instanceID, false)
	if err != nil {
		err = fmt.Errorf("Unable to stop instance: %v", aliyunSafeError(err))
		return
	}

	// Remove initialized tag
	tags := map[string]string{
		InitializedTagKey: "false",
	}
	// RemoveTags doesn't work. Not sure why
	err = p.ecs.AddTags(&ecs.AddTagsArgs{
		RegionId:     common.Region(networkSpec.Region),
		ResourceId:   instanceID,
		ResourceType: ecs.TagResourceInstance,
		Tag:          tags,
	})

	if err != nil {
		err = aliyunSafeError(err)
		glog.Infof("Unable to untag instance uninitialized!  %+v", err)
		return
	}

	// Find the system disk
	disks, err := p.ecs.DescribeDisks(&ecs.DescribeDisksArgs{
		RegionId:   common.Region(networkSpec.Region),
		InstanceId: instanceID,
	})
	if err != nil {
		err = aliyunSafeError(err)
		glog.Warningf("Error getting Aliyun vps disks: %+v", err)
		return
	}
	sdisk := (*ecs.DiskItemType)(nil)
	for _, d := range disks {
		if d.Type == ecs.DiskTypeAllSystem || d.Type == ecs.DiskTypeAll {
			sdisk = &d
		}
	}
	if sdisk == nil {
		err = fmt.Errorf("Unable to find the system disk of this vps %s", instanceID)
		return
	}

	err = p.ecs.WaitForInstance(instanceID, ecs.Stopped, 0)
	if err != nil {
		err = aliyunSafeError(err)
		glog.Warningf("Error waiting Aliyun vps: %+v", err)
		return
	}

	// Reinit disk
	err = p.ecs.ReInitDisk(sdisk.DiskId)
	if err != nil {
		err = aliyunSafeError(err)
		glog.Warningf("Error resetting Aliyun vps disk: %+v", err)
		return err
	}

	return
}

func (p *aliyunCloud) createInstance(clusterName string, instance *cluster.Instance) (status *cluster.InstanceStatus, err error) {
	var vpsID string
	defer func() {
		if err != nil && vpsID != "" {
			err2 := p.deleteInstance(vpsID)
			if err2 != nil {
				glog.Errorf("Roll back instance creation failed: %v", err2)
			}
		}
	}()

	glog.Infof("Creating instance: %v", instance.Name)
	options := cluster.InstanceOptions{}
	err = util.MapToStruct(instance.Labels, &options, cluster.AnnotationPrefix)
	if err != nil {
		err = fmt.Errorf("Can't get instance options: %s", err.Error())
		return
	}

	networkSpec := instance.Dependency.Network.Spec

	an := AliyunNetwork{}
	err = util.MapToStruct(instance.Dependency.Network.Annotations, &an, AliyunAnnotationPrefix)
	if err != nil || an.VSwitch == "" || an.VPC == "" {
		err = fmt.Errorf("Can't get network from instance annotations: %+v", err)
		return
	}

	// Aliyun defaults
	aio := AliyunInstanceOptions{
		InternetMaxBandwidthOut: 100,
		InternetMaxBandwidthIn:  200,
	}
	err = util.MapToStruct(instance.Annotations, &aio, AliyunAnnotationPrefix)
	if err != nil && instance.Annotations != nil {
		err = fmt.Errorf("Unable to get aliyun instance options: %v", err)
		return
	}

	// Instance type
	instanceType := instance.Spec.InstanceType
	if instanceType == "" {
		return nil, fmt.Errorf("Instance type must be specified")
	}

	// Image and its root device
	image := instance.Spec.Image
	if image == "" {
		return nil, fmt.Errorf("Instance image must be specified")
	}

	args := &ecs.CreateInstanceArgs{
		RegionId:                common.Region(networkSpec.Region),
		ZoneId:                  networkSpec.Zone,
		ImageId:                 image,
		InstanceType:            instanceType,
		SecurityGroupId:         an.SecurityGroup,
		InstanceName:            instance.Name,
		Description:             "Archon managed instance",
		HostName:                instance.Spec.Hostname,
		PrivateIpAddress:        options.UsePrivateIP,
		IoOptimized:             ecs.IoOptimizedOptimized,
		InternetChargeType:      common.PayByTraffic,
		InternetMaxBandwidthOut: aio.InternetMaxBandwidthOut,
		InternetMaxBandwidthIn:  aio.InternetMaxBandwidthIn,
		VSwitchId:               an.VSwitch,
		ClientToken:             util.RandNano(),
	}

	if aio.SystemDiskSize > 0 || aio.SystemDiskType != "" {
		args.SystemDisk = ecs.SystemDiskType{
			Size:     aio.SystemDiskSize,
			Category: ecs.DiskCategory(aio.SystemDiskType),
		}
	}

	vpsID, err = p.ecs.CreateInstance(args)
	if err != nil {
		err = aliyunSafeError(err)
		return
	}

	// Wait until it's stopped
	p.ecs.WaitForInstance(vpsID, ecs.Stopped, 0)

	// EIP. Have to bind here to ensure we get the correct ips when syncing instance
	// status from the cloudprovider
	eip := EIP{}
	err = util.MapToStruct(instance.Annotations, &eip, AliyunAnnotationPrefix)
	if err != nil {
		return
	}
	if eip.AllocationID != "" {
		resp, err2 := p.ecs.DescribeEipAddress(networkSpec.Region, eip.AllocationID)
		if err2 != nil {
			return nil, err2
		}

		if resp.InstanceId != vpsID && resp.Status != ecs.EipStatusAvailable {
			return nil, fmt.Errorf("EIP address is in use by another instance: %+v", resp)
		} else if resp.Status == ecs.EipStatusAvailable {
			err = p.ecs.AssociateEipAddress(eip.AllocationID, vpsID)
			if err != nil {
				return
			}
		}

		// Already associated
	}

	status, err = p.getInstance(networkSpec.Region, vpsID)
	if err != nil {
		err = aliyunSafeError(err)
		return
	}

	glog.Infof("Instance is created. %v", instance.Name)

	return
}

func (p *aliyunCloud) initializeInstance(clusterName string, instance *cluster.Instance) (status *cluster.InstanceStatus, err error) {
	vpsID := instance.Status.InstanceID

	defer func() {
		if err != nil {
			err2 := p.ecs.StopInstance(vpsID, true)
			if err2 != nil {
				glog.Errorf("Unable to rollback failed initialization: %v", err2)
			}
		}
	}()

	glog.Infof("Initializing instance: %s", instance.Name)
	// User data
	u, err := userdata.Generate(instance)
	if err != nil {
		return nil, err
	}

	args := &ecs.ModifyInstanceAttributeArgs{
		InstanceId: vpsID,
	}

	// Set password if provided in secret
	for _, s := range instance.Dependency.Secrets {
		if s.Type == v1.SecretTypeBasicAuth {
			_, ok := s.Data["username"]
			if ok {
				glog.V(4).Infof("Username in secret %s is ignored", s.Name)
			}

			password, ok := s.Data["password"]
			if ok {
				args.Password = string(password)
				break
			}
		}
	}

	// Since Aliyun doesn't support allocating IP beforehand, or userdata in coreos,
	// We will have to ssh into it and start cloudinit manually. To do this, we need
	// a password. If it's not provided by the user, we generate a difficult one
	if args.Password == "" {
		args.Password = util.RandPassword(30)
	}

	err = p.ecs.ModifyInstanceAttribute(args)
	if err != nil {
		return nil, err
	}

	/*
		// Modify userdata, aliyun api will encode it
		err = p.ecs.ModifyInstanceAttribute(&ecs.ModifyInstanceAttributeArgs{
				InstanceId: vpsID,
				UserData:   u,
				})
		if err != nil {
			return nil, err
		}
	*/

	// Start instance
	err = p.ecs.StartInstance(vpsID)
	if err != nil {
		err = aliyunSafeError(err)
		return nil, err
	}

	// Userdata is not supported, wait until the instance is started
	// and ssh-cloudinit with userdata

	// Wait until it's running
	p.ecs.WaitForInstance(vpsID, ecs.Running, 0)

	// Try cloudinit
	wait.PollImmediate(CloudInitInterval, CloudInitTimeout, func() (bool, error) {
		err = cloudinit.Run(&cloudinit.Config{
			Hostname: instance.Status.PublicIP,
			Port:     22,
			User:     SSHUsername,
			Password: args.Password,
			Os:       SSHSystem,
			Stdout:   &LogWriter{},
			UserData: string(u),
		})
		if err != nil {
			glog.V(4).Infof("Cloudinit failed. Still waiting: %v", err)
		}
		return err == nil, nil
	})

	if err != nil {
		return nil, fmt.Errorf("Unable to ssh-cloudinit the instance: %v", err)
	}

	// Put initialized annotation
	ai := AliyunInstanceInitialized{
		Initialized: true,
	}
	if instance.Annotations == nil {
		instance.Annotations = make(map[string]string)
	}
	err = util.StructToMap(ai, instance.Annotations, AliyunAnnotationPrefix)
	if err != nil {
		return nil, fmt.Errorf("Unable to set initialized flag: %v", err)
	}

	// Put initialized tag
	tags := map[string]string{
		InitializedTagKey: "true",
	}
	network := instance.Dependency.Network
	err = p.ecs.AddTags(&ecs.AddTagsArgs{
		RegionId:     common.Region(network.Spec.Region),
		ResourceId:   vpsID,
		ResourceType: ecs.TagResourceInstance,
		Tag:          tags,
	})
	if err != nil {
		err = aliyunSafeError(err)
		glog.Infof("Unable to tag instance as initialized!  %+v", err)
		return
	}

	// Return latest status
	status, err = p.getInstance(network.Spec.Region, vpsID)
	if err != nil {
		err = aliyunSafeError(err)
	} else {
		glog.Infof("Instance is initialized %+v", status)
	}
	return
}

func (p *aliyunCloud) EnsureInstanceDeleted(clusterName string, instance *cluster.Instance) (err error) {
	if instance.Status.InstanceID == "" {
		return nil
	}

	policy := instance.Spec.ReclaimPolicy
	switch policy {
	case cluster.InstanceReclaimDelete:
		err = p.deleteInstance(instance.Status.InstanceID)
		if err != nil {
			return
		}
		p.ecs.WaitForInstance(instance.Status.InstanceID, ecs.Deleted, 0)
	case cluster.InstanceReclaimRecycle:
		err = p.resetInstance(instance)
		if err != nil {
			return
		}
	default:
		err = fmt.Errorf("Unsupported instance reclaim policy for instance %v: %v", instance.Name, policy)
		return
	}

	instance.Status.InstanceID = ""

	return
}

func (p *aliyunCloud) deleteInstance(vpsID string) (err error) {
	p.ecs.StopInstance(vpsID, false)

	err = p.ecs.WaitForInstance(vpsID, ecs.Stopped, 0)
	gone := isNotFound(err)
	if err != nil && !gone {
		glog.Warningf("Error stopping Aliyun vps: %+v. Will try killing it", aliyunSafeError(err))
		p.ecs.StopInstance(vpsID, true)
	}

	err = p.ecs.DeleteInstance(vpsID)
	if err != nil {
		return aliyunSafeError(err)
	}

	return
}
