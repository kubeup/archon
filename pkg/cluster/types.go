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
package cluster

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/v1"
	//"k8s.io/kubernetes/pkg/apis/meta/v1"
)

const AnnotationPrefix = "archon.kubeup.com/"

// LocalObjectReference contains enough information to let you locate the referenced object inside the same namespace.
type LocalObjectReference struct {
	//TODO: Add other useful fields.  apiVersion, kind, uid?
	Name string `json:"name,omitempty"`
}

type ConfigSpec struct {
	Name string            `json:"name,omitempty"`
	Data map[string]string `json:"data,omitempty"`
}

type UserSpec struct {
	Name              string   `json:"name,omitempty"`
	PasswordHash      string   `json:"passwordHash,omitempty"`
	SSHAuthorizedKeys []string `json:"sshAuthorizedKeys,omitempty"`
	Sudo              string   `json:"sudo,omitempty"`
	Shell             string   `json:"shell,omitempty"`
}

type User struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              UserSpec `json:"spec"`
}

type UserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []User `json:"items"`
}

type NetworkSpec struct {
	Region string `k8s:"region" json:"region,omitempty"`
	Zone   string `k8s:"zone" json:"zone,omitempty"`
	Subnet string `k8s:"subnet" json:"subnet,omitempty"`
}

type NetworkStatus struct {
	Phase NetworkPhase `json:"phase,omitempty"`
}

type NetworkPhase string

const (
	NetworkPending NetworkPhase = "Pending"
	NetworkRunning NetworkPhase = "Running"
	NetworkFailed  NetworkPhase = "Failed"
	NetworkUnknown NetworkPhase = "Unknown"
)

type Network struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              NetworkSpec   `json:"spec,omitempty"`
	Status            NetworkStatus `json:"status,omitempty"`
}

type NetworkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Network `json:"items"`
}

type InstanceGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              InstanceGroupSpec   `json:"spec,omitempty"`
	Status            InstanceGroupStatus `json:"status,omitempty"`
}

type InstanceGroupSpec struct {
	Replicas        int32                        `json:"replicas,omitempty"`
	ProvisionPolicy InstanceGroupProvisionPolicy `json:"provisionPolicy,omitempty"`

	// Minimum number of seconds for which a newly created instance should be ready
	// without any of its container crashing, for it to be considered available.
	// Defaults to 0 (instance will be considered available as soon as it is ready)
	// +optional
	MinReadySeconds          int32                 `json:"minReadySeconds,omitempty"`
	Selector                 *metav1.LabelSelector `json:"selector,omitempty"`
	ReservedInstanceSelector *metav1.LabelSelector `json:"reservedInstaceSelector,omitempty"`
	Template                 InstanceTemplateSpec  `json:"template,omitempty"`
}

type InstanceGroupConditionType string

// These are valid conditions of a replica set.
const (
	// InstanceGroupReplicaFailure is added in a replica set when one of its instances fails to be created
	// due to insufficient quota, limit ranges, instance security policy, node selectors, etc. or deleted
	// due to kubelet being down or finalizers are failing.
	InstanceGroupReplicaFailure InstanceGroupConditionType = "ReplicaFailure"
)

// InstanceGroupCondition describes the state of a replica set at a certain point.
type InstanceGroupCondition struct {
	// Type of replica set condition.
	Type InstanceGroupConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status api.ConditionStatus `json:"status"`
	// The last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// The reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`
	// A human readable message indicating details about the transition.
	// +optional
	Message string `json:"message,omitempty"`
}

type InstanceGroupStatus struct {
	// Replicas is the number of actual replicas.
	Replicas int32 `json:"replicas"`

	// The number of instances that have labels matching the labels of the instance template of the instancegroup.
	// +optional
	FullyLabeledReplicas int32 `json:"fullyLabeledReplicas,omitempty"`

	// The number of ready replicas for this replica set.
	// +optional
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`

	// The number of available replicas (ready for at least minReadySeconds) for this replica set.
	// +optional
	AvailableReplicas int32 `json:"availableReplicas,omitempty"`

	// ObservedGeneration is the most recent generation observed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	Conditions []InstanceGroupCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

type InstanceGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []InstanceGroup `json:"items"`
}

type InstanceGroupProvisionPolicy string

const (
	InstanceGroupProvisionReservedOnly  InstanceGroupProvisionPolicy = "ReservedOnly"
	InstanceGroupProvisionDynamicOnly   InstanceGroupProvisionPolicy = "DynamicOnly"
	InstanceGroupProvisionReservedFirst InstanceGroupProvisionPolicy = "ReservedFirst"
)

type FileSpec struct {
	Name               string `json:"name,omitempty" yaml:"name,omitempty"`
	Encoding           string `json:"encoding,omitempty" yaml:"encoding,omitempty" valid:"^(base64|b64|gz|gzip|gz\\+base64|gzip\\+base64|gz\\+b64|gzip\\+b64)$"`
	Content            string `json:"content,omitempty" yaml:"content,omitempty"`
	Template           string `json:"template,omitempty" yaml:"template,omitempty"`
	Owner              string `json:"owner,omitempty" yaml:"owner,omitempty"`
	UserID             int    `json:"userID,omitempty" yaml:"userID,omitempty"`
	GroupID            int    `json:"groupID,omitempty" yaml:"groupID,omitempty"`
	Filesystem         string `json:"filesystem,omitempty" yaml:"filesystem,omitempty"`
	Path               string `json:"path,omitempty" yaml:"path,omitempty"`
	RawFilePermissions string `json:"permissions,omitempty" yaml:"permissions,omitempty" valid:"^0?[0-7]{3,4}$"`
}

type InstanceTemplateSpec struct {
	metav1.ObjectMeta `json:"metadata"`
	Spec              InstanceSpec `json:"spec,omitempty"`
	Secrets           []v1.Secret  `json:"secrets,omitempty"`
}

type InstanceOptions struct {
	UseInstanceID string `k8s:"use-instance-id"`
	UsePrivateIP  string `k8s:"use-private-ip"`
}

type InstanceDependency struct {
	Network          Network          `json:"network,omitempty"`
	Secrets          []v1.Secret      `json:"secrets,omitempty"`
	Users            []User           `json:"users,omitempty"`
	ReservedInstance ReservedInstance `json:"reservedInstance,omitempty"`
}

type Instance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              InstanceSpec       `json:"spec,omitempty"`
	Status            InstanceStatus     `json:"status,omitempty"`
	Dependency        InstanceDependency `json:"-"`
}

var _ metav1.Object = &Instance{}

type InstanceSpec struct {
	OS                  string                 `json:"os,omitempty"`
	Image               string                 `json:"image,omitempty"`
	InstanceType        string                 `json:"instanceType,omitempty"`
	NetworkName         string                 `json:"networkName,omitempty"`
	ReclaimPolicy       InstanceReclaimPolicy  `json:"reclaimPolicy,omitempty"`
	Files               []FileSpec             `json:"files,omitempty"`
	Secrets             []LocalObjectReference `json:"secrets,omitempty"`
	Configs             []ConfigSpec           `json:"configs,omitempty"`
	Users               []LocalObjectReference `json:"users,omitempty"`
	Hostname            string                 `json:"hostname,omitempty"`
	ReservedInstanceRef *LocalObjectReference  `json:"reservedInstanceRef,omitempty"`
}

type InstanceStatus struct {
	Phase      InstancePhase       `json:"phase,omitempty"`
	Conditions []InstanceCondition `json:"conditions,omitempty"`
	// TODO: allow multiple ips
	PrivateIP         string      `json:"privateIP,omitempty"`
	PublicIP          string      `json:"publicIP,omitempty"`
	InstanceID        string      `json:"instanceID,omitempty"`
	CreationTimestamp metav1.Time `json:"creationTimestamp,omitempty" protobuf:"bytes,8,opt,name=creationTimestamp"`
}

type InstancePhase string
type InstanceReclaimPolicy string

const (
	// TODO: Populate defaults when a new resource is created. So we don't have use
	// "" as a meaningful and the default state
	InstancePending      InstancePhase = "Pending"
	InstanceInitializing InstancePhase = "Initializing"
	InstanceRunning      InstancePhase = "Running"
	InstanceFailed       InstancePhase = "Failed"
	InstanceUnknown      InstancePhase = "Unknown"

	InstanceReclaimRecycle InstanceReclaimPolicy = "Recycle"
	InstanceReclaimDelete  InstanceReclaimPolicy = "Delete"

	InstanceOSCoreOS string = "CoreOS"
	InstanceOSUbuntu string = "Ubuntu"
	InstanceOSCentOS string = "CentOS"

	// Instance with this annotation will be removed first if InstanceGroup controller
	// is scaling down an InstanceGroup
	InstanceToBeRemovedAnnotation string = AnnotationPrefix + "to-be-removed"
)

type InstanceConditionType string

// These are valid conditions of instance.
const (
	// InstanceScheduled represents status of the scheduling process for this instance.
	InstanceScheduled InstanceConditionType = "InstanceScheduled"
	// InstanceReady means the instance is able to service requests and should be added to the
	// load balancing pools of all matching services.
	InstanceReady InstanceConditionType = "Ready"
	// InstanceInitialized means that all init containers in the instance have started successfully.
	InstanceInitialized InstanceConditionType = "Initialized"
	// InstanceReasonUnschedulable reason in InstanceScheduled InstanceCondition means that the scheduler
	// can't schedule the instance right now, for example due to insufficient resources in the cluster.
	InstanceReasonUnschedulable = "Unschedulable"
)

type InstanceCondition struct {
	Type   InstanceConditionType `json:"type"`
	Status api.ConditionStatus   `json:"status"`
	// +optional
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty"`
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// +optional
	Reason string `json:"reason,omitempty"`
	// +optional
	Message string `json:"message,omitempty"`
}

type InstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Instance `json:"items"`
}

type ReservedInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              ReservedInstanceSpec   `json:"spec,omitempty"`
	Status            ReservedInstanceStatus `json:"status,omitempty"`
}

type ReservedInstanceSpec struct {
	OS           string       `json:"os,omitempty"`
	Image        string       `json:"image,omitempty"`
	InstanceType string       `json:"instanceType,omitempty"`
	NetworkName  string       `json:"networkName,omitempty"`
	Hostname     string       `json:"hostname,omitempty"`
	InstanceID   string       `json:"instanceID,omitempty"`
	Configs      []ConfigSpec `json:"configs,omitempty"`
}

type ReservedInstanceStatus struct {
	Phase ReservedInstancePhase `json:"phase,omitempty"`
	// Binding reference to the instance
	InstanceName string `json:"instanceName,omitempty"`
}

type ReservedInstancePhase string

const (
	ReservedInstanceBound     ReservedInstancePhase = "Bound"
	ReservedInstanceAvailable ReservedInstancePhase = "Available"
)

type ReservedInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []ReservedInstance `json:"items"`
}
