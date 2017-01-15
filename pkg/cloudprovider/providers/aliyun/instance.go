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
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/util/wait"
	"kubeup.com/archon/pkg/cluster"
	"kubeup.com/archon/pkg/userdata"
	"kubeup.com/archon/pkg/util"
	"time"
)

var (
	ErrorNotFound     = fmt.Errorf("Instance is not found")
	CloudInitInterval = 5 * time.Second
	CloudInitTimeout  = 3 * time.Minute
	SSHUsername       = "root"
	SSHSystem         = "coreos"
	StateMap          = map[ecs.InstanceStatus]cluster.InstancePhase{
		ecs.Creating: cluster.InstancePending,
		ecs.Starting: cluster.InstancePending,
		ecs.Running:  cluster.InstanceRunning,
		ecs.Stopping: cluster.InstanceFailed,
		ecs.Stopped:  cluster.InstanceFailed,
	}
)

type AliyunInstanceOptions struct {
	InternetMaxBandwidthIn  int    `k8s:"internet-max-bandwidth-in"`
	InternetMaxBandwidthOut int    `k8s:"internet-max-bandwidth-out"`
	SystemDiskSize          int    `k8s:"system-disk-size"`
	SystemDiskType          string `k8s:"system-disk-type"`
}

func instanceToStatus(i ecs.InstanceAttributesType) *cluster.InstanceStatus {
	phase, ok := StateMap[i.Status]
	if !ok {
		glog.Warningf("Unknown instance state: %+v", i.Status)
		phase = cluster.InstanceUnknown
	}
	return &cluster.InstanceStatus{
		Phase:             phase,
		PrivateIP:         firstIP(i.VpcAttributes.PrivateIpAddress),
		PublicIP:          firstIP(i.PublicIpAddress),
		InstanceID:        i.InstanceId,
		CreationTimestamp: unversioned.NewTime(time.Time(i.CreationTime)),
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
	err = util.MapToStruct(instance.Annotations, &an, AliyunAnnotationPrefix)
	if err != nil || an.VSwitch == "" {
		err = fmt.Errorf("Network is not ready. Can't create instance: %v", err)
		return
	}

	if instance.Status.InstanceID != "" {
		status2 := (*cluster.InstanceStatus)(nil)
		status2, err = p.getInstance(instance.Dependency.Network.Spec.Region, instance.Status.InstanceID)

		if err != nil {
			if err == ErrorNotFound {
				return p.createInstance(clusterName, instance)
			}
			return
		}

		switch status2.Phase {
		case cluster.InstanceFailed, cluster.InstanceUnknown:
			err = p.EnsureInstanceDeleted(clusterName, instance)
			if err != nil {
				return
			}
			return p.createInstance(clusterName, instance)
		}

		status = status2
	} else {
		return p.createInstance(clusterName, instance)
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

	options := cluster.InstanceOptions{}
	if instance.Labels != nil {
		err = util.MapToStruct(instance.Labels, &options, cluster.AnnotationPrefix)
		if err != nil {
			err = fmt.Errorf("Can't get instance options: %s", err.Error())
			return
		}
	}

	networkSpec := instance.Dependency.Network.Spec

	an := AliyunNetwork{}
	err = util.MapToStruct(instance.Annotations, &an, AliyunAnnotationPrefix)
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
	if err != nil {
		err = fmt.Errorf("Unable to get aliyun instance options: %v", err)
		return
	}

	// Ignore preallocate public/private ip options. Always create instance frist, generate userdata later, then start the instance

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

	// Set password if provided in secret
	for _, s := range instance.Dependency.Secrets {
		if s.Type == api.SecretTypeBasicAuth {
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

	vpsID, err = p.ecs.CreateInstance(args)
	if err != nil {
		err = aliyunSafeError(err)
		return
	}

	// Wait until it's stopped
	p.ecs.WaitForInstance(vpsID, ecs.Stopped, 0)

	// Public IP
	publicIP, err := p.ecs.AllocatePublicIpAddress(vpsID)
	if err != nil {
		err = aliyunSafeError(err)
		return
	}

	status, err = p.getInstance(networkSpec.Region, vpsID)
	if err != nil {
		err = aliyunSafeError(err)
		return
	}

	oldStatus := instance.Status
	instance.Status = *status

	// User data
	u, err := userdata.Generate(instance)
	if err != nil {
		return nil, err
	}

	instance.Status = oldStatus

	/*
		// TODO: aliyun doesn't support modifying userdata when coreos is used as image
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
			Hostname: publicIP,
			Port:     22,
			User:     SSHUsername,
			Password: args.Password,
			Os:       SSHSystem,
			Stdout:   &LogWriter{},
			UserData: string(u),
		})
		if err != nil {
			glog.Infof("Cloudinit failed. Still waiting: %v", err)
		}
		return err == nil, nil
	})

	if err != nil {
		return nil, fmt.Errorf("Unable to ssh-cloudinit the instance: %v", err)
	}

	// Return latest status
	status, err = p.getInstance(networkSpec.Region, vpsID)
	if err != nil {
		err = aliyunSafeError(err)
	} else {
		glog.Infof("New instance created %+v", status)
	}
	return
}

func (p *aliyunCloud) EnsureInstanceDeleted(clusterName string, instance *cluster.Instance) (err error) {
	if instance.Status.InstanceID == "" {
		return nil
	}

	return p.deleteInstance(instance.Status.InstanceID)
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
