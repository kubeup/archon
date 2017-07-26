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
	"encoding/base64"
	"fmt"
	"github.com/golang/glog"
	cloudinit "github.com/tryk8s/ssh-cloudinit/client"
	"github.com/ucloud/ucloud-sdk-go/service/uhost"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/kubernetes/pkg/api/v1"
	"kubeup.com/archon/pkg/cluster"
	"kubeup.com/archon/pkg/userdata"
	"kubeup.com/archon/pkg/util"
	"strings"
	"time"
)

var (
	ErrorNotFound     = fmt.Errorf("Instance is not found")
	CloudInitInterval = 5 * time.Second
	CloudInitTimeout  = 3 * time.Minute
	// UCloud doesn't allow tag keys starting with "aliyun"
	InitializedTagKey = "instance." + UCloudAnnotationPrefix + "initialized"
	SSHUsernameMap    = map[string]string{
		"ubuntu": "ubuntu",
		"centos": "root",
	}
	StateMap = map[string]cluster.InstancePhase{
		InstanceStateInitializing: cluster.InstancePending,
		InstanceStateStarting:     cluster.InstancePending,
		InstanceStateRunning:      cluster.InstanceRunning,
		InstanceStateStopping:     cluster.InstanceFailed,
		InstanceStateStopped:      cluster.InstanceFailed,
		InstanceStateFailed:       cluster.InstanceFailed,
		InstanceStateRebooting:    cluster.InstancePending,
	}
)

type UCloudInstanceOptions struct {
	CPU             int    `k8s:"cpu"`
	Memory          int    `k8s:"memory"`
	SystemDiskSize  int    `k8s:"system-disk-size"`
	StorageDiskSize int    `k8s:"storage-disk-size"`
	DiskType        string `k8s:"disk-type"`
	NetCapability   string `k8s:"net-capability"`
	UseSSH          bool   `k8s:"use-ssh"`
}

type UCloudInstanceInitialized struct {
	Initialized bool `k8s:"instance-initialized"`
}

func instanceToStatus(i *uhost.UHostSet) *cluster.InstanceStatus {
	if i == nil {
		return nil
	}

	phase, ok := StateMap[i.State]
	if !ok {
		glog.Warningf("Unknown instance state: %+v", i.State)
		phase = cluster.InstanceUnknown
	}

	// If the instance is not marked as initialized and its vps is in stopped state,
	// return Pending
	initialized := false
	if i.Tag == InitializedTagKey {
		initialized = true
	}
	if i.State == InstanceStateStopped && initialized == false {
		phase = cluster.InstancePending
	}

	var (
		publicIP  string
		privateIP string
	)

	for _, ip := range i.IPSet {
		if ip.Type == InstanceIPTypePrivate && privateIP == "" {
			privateIP = ip.IP
		} else if publicIP == "" {
			publicIP = ip.IP
		}
	}

	return &cluster.InstanceStatus{
		Phase:             phase,
		PrivateIP:         privateIP,
		PublicIP:          publicIP,
		InstanceID:        i.UHostId,
		CreationTimestamp: metav1.NewTime(time.Unix(int64(i.CreateTime), 0)),
	}
}

func (p *UCloud) ListInstances(clusterName string, network *cluster.Network, selector map[string]string) (names []string, statuses []*cluster.InstanceStatus, err error) {
	instances, err := p.ucloud.DescribeInstances(&uhost.DescribeUHostInstanceParams{})

	if err != nil {
		err = safeError(err)
		return
	}

	for _, instance := range instances {
		names = append(names, instance.Name)
		statuses = append(statuses, instanceToStatus(&instance))
	}

	return
}

func (p *UCloud) GetInstance(clusterName string, instance *cluster.Instance) (status *cluster.InstanceStatus, err error) {
	return p.getInstance(instance.Status.InstanceID)
}

func (p *UCloud) getInstance(instanceID string) (status *cluster.InstanceStatus, err error) {
	if instanceID == "" {
		return nil, ErrorNotFound
	}

	instance, err := p.ucloud.DescribeInstanceAttribute(instanceID)

	if err != nil {
		if isNotFound(err) {
			err = ErrorNotFound
		} else {
			err = safeError(err)
		}

		return
	}

	if instance == nil {
		err = ErrorNotFound
		return
	}

	return instanceToStatus(instance), nil
}

func (p *UCloud) EnsureInstance(clusterName string, instance *cluster.Instance) (status *cluster.InstanceStatus, err error) {
	var vpsID string
	defer func() {
		if err != nil && vpsID != "" {
			err2 := p.deleteInstance(vpsID)
			if err2 != nil {
				glog.Errorf("Roll back instance creation failed: %v", err2)
			} else {
				instance.Status.InstanceID = ""
			}
		}
	}()

	an := UCloudNetwork{}
	err = util.MapToStruct(instance.Dependency.Network.Annotations, &an, UCloudAnnotationPrefix)
	if err != nil {
		err = fmt.Errorf("Can't get network annotations. Can't create instance: %v", err)
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
			status, err = p.getInstance(options.UseInstanceID)
		} else {
			status, err = p.createInstance(clusterName, instance)
			if status != nil && status.InstanceID != "" {
				vpsID = status.InstanceID
			}
		}

		if err != nil {
			return
		}

		instance.Status = *status
	}

	if instance.Status.Phase == cluster.InstanceRunning {
		return &instance.Status, nil
	}

	ai := UCloudInstanceInitialized{}
	err = util.MapToStruct(instance.Annotations, &ai, UCloudAnnotationPrefix)
	if err != nil && instance.Annotations != nil {
		err = fmt.Errorf("Can't tell if the instance is initialized: %v", err)
		return
	}

	if ai.Initialized == true {
		return &instance.Status, nil
	}

	if instance.Dependency.ReservedInstance.Spec.InstanceID != "" {
		// It's initializing from a reserved instance. Reset again to modify its password
		p.resetInstance(instance)
	}

	return p.initializeInstance(clusterName, instance)
}

func (p *UCloud) resetInstance(instance *cluster.Instance) (err error) {
	options := cluster.InstanceOptions{}
	err = util.MapToStruct(instance.Labels, &options, cluster.AnnotationPrefix)
	if err != nil {
		err = fmt.Errorf("Can't get instance options: %s", err.Error())
		return
	}

	instanceID := instance.Status.InstanceID
	if instanceID == "" {
		instanceID = options.UseInstanceID
	}

	// Ignore err in case the instance is already stopped
	err = p.ucloud.StopInstance(instanceID, true)
	if err != nil {
		glog.Warningf("Unable to stop instance: %v", safeError(err))
	}

	// Remove initialized tag
	err = p.ucloud.ModifyRemark(instanceID, "")

	if err != nil {
		err = safeError(err)
		glog.Infof("Unable to untag instance uninitialized!  %+v", err)
		return
	}

	// Reinit disk
	password, err := p.getInstancePassword(instance)
	if err != nil {
		return
	}

	err = p.ucloud.ReinstallInstance(instanceID, password)
	if err != nil {
		err = safeError(err)
		glog.Warningf("Error resetting UCloud vps disk: %+v", err)
		return err
	}

	return
}

func (p *UCloud) createInstance(clusterName string, instance *cluster.Instance) (status *cluster.InstanceStatus, err error) {
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

	an := UCloudNetwork{}
	err = util.MapToStruct(instance.Dependency.Network.Annotations, &an, UCloudAnnotationPrefix)
	if err != nil {
		err = fmt.Errorf("Can't get network annotations: %+v", err)
		return
	}

	// UCloud defaults
	aio := UCloudInstanceOptions{
		CPU:    1,
		Memory: 2048,
		UseSSH: true,
	}
	err = util.MapToStruct(instance.Annotations, &aio, UCloudAnnotationPrefix)
	if err != nil && instance.Annotations != nil {
		err = fmt.Errorf("Unable to get aliyun instance options: %v", err)
		return
	}

	// Instance type
	instanceType := instance.Spec.InstanceType
	if instanceType == "" {
		return nil, fmt.Errorf("Instance type must be specified")
	}

	// Image
	image := instance.Spec.Image
	if image == "" {
		return nil, fmt.Errorf("Instance image must be specified")
	}

	// TODO: Fix ucloud api
	args := &uhost.CreateUHostInstanceParams{
		Zone:            networkSpec.Zone,
		ImageId:         image,
		SecurityGroupId: an.SecurityGroup,
		LoginMode:       InstanceLoginModePassword,
		ChargeType:      InstanceChargeTypeDynamic,
		UHostType:       instanceType,
		Name:            instance.Name,
		CPU:             aio.CPU,
		Memory:          aio.Memory,
		DiskSpace:       aio.StorageDiskSize,
		BootDiskSpace:   aio.SystemDiskSize,
		//StorageType:     aio.StorageType,
		NetworkId:     an.Subnet,
		NetCapability: aio.NetCapability,
	}

	// Set password if provided in secret
	args.Password, err = p.getInstancePassword(instance)
	if err != nil {
		return
	}

	vpsID, err = p.ucloud.CreateInstance(args)
	if err != nil {
		err = safeError(err)
		return
	}

	// Wait until it's stopped
	err = p.ucloud.WaitForInstance(vpsID, InstanceStateRunning, 0)
	if err != nil {
		err = safeError(err)
		return
	}

	err = p.ucloud.StopInstance(vpsID, false)
	if err != nil {
		err = safeError(err)
		return
	}

	p.ucloud.WaitForInstance(vpsID, InstanceStateStopped, 0)
	if err != nil {
		err = safeError(err)
		return
	}

	// EIP. Have to bind here to ensure we get the correct ips when syncing instance
	// status from the cloudprovider
	eip := EIP{}
	err = util.MapToStruct(instance.Annotations, &eip, UCloudAnnotationPrefix)
	if err != nil {
		return
	}
	if eip.ID != "" {
		resp, err2 := p.ucloud.DescribeEipAddress(eip.ID)
		if err2 != nil {
			return nil, err2
		}

		glog.Warning("EIP: %+v", resp)
		eipInstance := ""
		if resp.Resource != nil {
			resource := (*resp.Resource)
			if resource.ResourceType == ResourceTypeUHost {
				eipInstance = resource.ResourceId
			}
		}

		if resp.Status != EIPStatusFree && eipInstance != vpsID {
			return nil, fmt.Errorf("EIP address is in use by another instance: %+v", resp)
		} else if resp.Status == EIPStatusFree {
			err = p.ucloud.AssociateEipAddress(vpsID, eip.ID)
			if err != nil {
				return
			}
		}

		// Already associated
	}

	//
	status, err = p.getInstance(vpsID)
	if err != nil {
		err = safeError(err)
		return
	}

	glog.Infof("Instance is created. %v", instance.Name)

	return
}

func (p *UCloud) getInstancePassword(instance *cluster.Instance) (string, error) {
	var (
		password []byte
		ok       bool
	)
	for _, s := range instance.Dependency.Secrets {
		if s.Type == v1.SecretTypeBasicAuth {
			_, ok = s.Data["username"]
			if ok {
				glog.V(4).Infof("Username in secret %s is ignored", s.Name)
			}

			password, ok = s.Data["password"]
			if ok {
				return base64.StdEncoding.EncodeToString(password), nil
			}
		}
	}

	if len(password) == 0 {
		// We will have to ssh into it and start cloudinit manually. To do this, we need
		// a password. If it's not provided by the user, we generate one
		return "", fmt.Errorf("UCloud instance requires a password but none is specified.")
	}

	return string(password), nil
}

func (p *UCloud) initializeInstance(clusterName string, instance *cluster.Instance) (status *cluster.InstanceStatus, err error) {
	vpsID := instance.Status.InstanceID

	defer func() {
		if err != nil {
			err2 := p.ucloud.StopInstance(vpsID, false)
			if err2 != nil {
				glog.Errorf("Unable to rollback failed initialization: %v", err2)
			} else {
				p.ucloud.WaitForInstance(vpsID, InstanceStateStopped, 0)
			}
		}
	}()

	glog.Infof("Initializing instance: %s", instance.Name)
	aio := UCloudInstanceOptions{
		UseSSH: true,
	}
	err = util.MapToStruct(instance.Annotations, &aio, UCloudAnnotationPrefix)
	if err != nil && instance.Annotations != nil {
		err = fmt.Errorf("Unable to get aliyun instance options: %v", err)
		return
	}

	// User data
	u, err := userdata.Generate(instance)
	if err != nil {
		return nil, fmt.Errorf("Failed to generated user-data: %v", err)
	}

	// Password
	password, err := p.getInstancePassword(instance)
	if err != nil {
		return
	}

	passwordData, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		err = fmt.Errorf("Unable to decode password: %v", err)
		return
	}
	password = string(passwordData)

	// Start instance
	err = p.ucloud.StartInstance(vpsID)
	if err != nil {
		err = safeError(err)
		return nil, fmt.Errorf("Error starting instance: %v", err)
	}

	// Wait until it's running
	p.ucloud.WaitForInstance(vpsID, InstanceStateRunning, 0)

	if aio.UseSSH {
		// Userdata is not supported, wait until the instance is started
		// and ssh-cloudinit with userdata

		if instance.Status.PublicIP == "" {
			err = fmt.Errorf("Unable to run ssh-cloudinit because instance %v has no public ip.", instance.Name)
			return
		}

		// Try cloudinit
		// Implement ssh timeout
		// TODO: If cloudinit is already there, don't upgrade
		os := strings.ToLower(instance.Spec.OS)
		username := "root"
		if _username, ok := SSHUsernameMap[os]; ok {
			username = _username
		}
		glog.V(2).Infof("Run ssh-cloudinit on the instance %v as %v", instance.Name, username)
		wait.PollImmediate(CloudInitInterval, CloudInitTimeout, func() (bool, error) {
			err = cloudinit.Run(&cloudinit.Config{
				Hostname:           instance.Status.PublicIP,
				Port:               22,
				User:               username,
				UseCloudDataSource: false,
				Password:           password,
				Os:                 os,
				Stdout:             &LogWriter{},
				UserData:           string(u),
			})
			if _, ok := err.(*cloudinit.ConnectionError); err != nil && ok {
				glog.V(2).Infof("Cloudinit failed. Still waiting: %v", err)
				return false, nil
			}

			return err == nil, err
		})

		if err != nil {
			return nil, fmt.Errorf("Unable to ssh-cloudinit the instance: %v", err, err)
		}
	}

	// Put initialized annotation
	ai := UCloudInstanceInitialized{
		Initialized: true,
	}
	if instance.Annotations == nil {
		instance.Annotations = make(map[string]string)
	}
	err = util.StructToMap(ai, instance.Annotations, UCloudAnnotationPrefix)
	if err != nil {
		return nil, fmt.Errorf("Unable to set initialized flag: %v", err)
	}

	// Put initialized tag
	err = p.ucloud.ModifyRemark(vpsID, InitializedTagKey)
	if err != nil {
		err = fmt.Errorf("Unable to tag instance: %v", safeError(err))
		glog.Infof("Unable to tag instance as initialized!  %+v", err)
		return
	}

	// Return latest status
	status, err = p.getInstance(vpsID)
	if err != nil {
		err = safeError(err)
		err = fmt.Errorf("Error get instance status: %v", err)
	} else {
		glog.Infof("Instance is initialized %+v", status)
	}
	return
}

func (p *UCloud) EnsureInstanceDeleted(clusterName string, instance *cluster.Instance) (err error) {
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

func (p *UCloud) deleteInstance(vpsID string) (err error) {
	p.ucloud.StopInstance(vpsID, false)

	err = p.ucloud.WaitForInstance(vpsID, InstanceStateStopped, 0)
	gone := isNotFound(err)
	if err != nil && !gone {
		glog.Warningf("Error stopping UCloud vps: %+v.", safeError(err))
		p.ucloud.StopInstance(vpsID, true)
	}

	err = p.ucloud.DeleteInstance(vpsID)
	gone = isNotFound(err)
	if err != nil && !gone {
		err = fmt.Errorf("Unable to delete instance %s: %v", vpsID, err)
		return safeError(err)
	}

	return nil
}
