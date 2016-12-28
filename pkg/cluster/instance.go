package cluster

import (
	"encoding/json"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/meta"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"time"
)

type InstanceTemplateSpec struct {
	api.ObjectMeta `json:"metadata"`
	Spec           InstanceSpec `json:"spec,omitempty"`
	Secrets        []api.Secret `json:"secrets,omitempty"`
}

type InstanceOptions struct {
	// A set of provider specific annotations will be set by the eip controller when an eip is automatically allocated
	PreallocatePublicIP bool `k8s:"preallocate-public-ip"`
	// If this is set
	PreallocatePrivateIP bool `k8s:"preallocate-private-ip"`
}

type InstanceDependency struct {
	Network Network      `json:"network,omitempty"`
	Secrets []api.Secret `json:"secrets,omitempty"`
	Users   []User       `json:"users,omitempty"`
}

type Instance struct {
	unversioned.TypeMeta `json:",inline"`
	api.ObjectMeta       `json:"metadata"`
	Spec                 InstanceSpec       `json:"spec,omitempty"`
	Status               InstanceStatus     `json:"status,omitempty"`
	Dependency           InstanceDependency `json:"-"`
}

type InstanceSpec struct {
	Image        string                     `json:"image,omitempty"`
	InstanceType string                     `json:"instanceType,omitempty"`
	NetworkName  string                     `json:"networkName,omitempty"`
	Files        []FileSpec                 `json:"files,omitempty"`
	Secrets      []api.LocalObjectReference `json:"secrets,omitempty"`
	Configs      []ConfigSpec               `json:"configs,omitempty"`
	Users        []api.LocalObjectReference `json:"users,omitempty"`
	Hostname     string                     `json:"hostname,omitempty"`
}

type InstanceStatus struct {
	Phase             InstancePhase       `json:"phase,omitempty"`
	Conditions        []InstanceCondition `json:"conditions,omitempty"`
	PrivateIP         string              `json:"privateIP,omitempty"`
	PublicIP          string              `json:"publicIP,omitempty"`
	ElasticIP         string              `json:"elasticIP,omitempty"`
	InstanceID        string              `json:"instanceID,omitempty"`
	CreationTimestamp unversioned.Time    `json:"creationTimestamp,omitempty" protobuf:"bytes,8,opt,name=creationTimestamp"`
}

type InstancePhase string

const (
	InstanceInitializing InstancePhase = "Initializing"
	InstancePending      InstancePhase = "Pending"
	InstanceRunning      InstancePhase = "Running"
	InstanceFailed       InstancePhase = "Failed"
	InstanceUnknown      InstancePhase = "Unknown"
)

type InstanceConditionType string

// These are valid conditions of pod.
const (
	// InstanceScheduled represents status of the scheduling process for this pod.
	InstanceScheduled InstanceConditionType = "InstanceScheduled"
	// InstanceReady means the pod is able to service requests and should be added to the
	// load balancing pools of all matching services.
	InstanceReady InstanceConditionType = "Ready"
	// InstanceInitialized means that all init containers in the pod have started successfully.
	InstanceInitialized InstanceConditionType = "Initialized"
	// InstanceReasonUnschedulable reason in InstanceScheduled InstanceCondition means that the scheduler
	// can't schedule the pod right now, for example due to insufficient resources in the cluster.
	InstanceReasonUnschedulable = "Unschedulable"
)

type InstanceCondition struct {
	Type   InstanceConditionType `json:"type"`
	Status api.ConditionStatus   `json:"status"`
	// +optional
	LastProbeTime unversioned.Time `json:"lastProbeTime,omitempty"`
	// +optional
	LastTransitionTime unversioned.Time `json:"lastTransitionTime,omitempty"`
	// +optional
	Reason string `json:"reason,omitempty"`
	// +optional
	Message string `json:"message,omitempty"`
}

type InstanceList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata"`

	Items []*Instance `json:"items"`
}

func (i Instance) GetObjectKind() unversioned.ObjectKind {
	return &i.TypeMeta
}

func (i Instance) GetObjectMeta() meta.Object {
	return &i.ObjectMeta
}

func (il InstanceList) GetObjectKind() unversioned.ObjectKind {
	return &il.TypeMeta
}

func (il InstanceList) GetObjectMeta() unversioned.List {
	return &il.ListMeta
}

// The code below is used only to work around a known problem with third-party
// resources and ugorji. If/when these issues are resolved, the code below
// should no longer be required.

type InstanceCopy Instance
type InstanceListCopy InstanceList

func (e *Instance) UnmarshalJSON(data []byte) error {
	tmp := InstanceCopy{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	tmp2 := Instance(tmp)
	*e = tmp2
	return nil
}

func (el *InstanceList) UnmarshalJSON(data []byte) error {
	tmp := InstanceListCopy{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	tmp2 := InstanceList(tmp)
	*el = tmp2
	return nil
}

func (e *Instance) CodecEncodeSelf() {
}

func (e *Instance) CodecDecodeSelf() {
}

func (el *InstanceList) CodecEncodeSelf() {
}

func (el *InstanceList) CodecDecodeSelf() {
}

// IsInstanceAvailable returns true if a instance is available; false otherwise.
// Precondition for an available instance is that it must be ready. On top
// of that, there are two cases when a instance can be considered available:
// 1. minReadySeconds == 0, or
// 2. LastTransitionTime (is set) + minReadySeconds < current time
func IsInstanceAvailable(instance *Instance, minReadySeconds int32, now unversioned.Time) bool {
	if !IsInstanceReady(instance) {
		return false
	}

	c := GetInstanceReadyCondition(instance.Status)
	minReadySecondsDuration := time.Duration(minReadySeconds) * time.Second
	if minReadySeconds == 0 || !c.LastTransitionTime.IsZero() && c.LastTransitionTime.Add(minReadySecondsDuration).Before(now.Time) {
		return true
	}
	return false
}

// IsInstanceReady returns true if a instance is ready; false otherwise.
func IsInstanceReady(instance *Instance) bool {
	return IsInstanceReadyConditionTrue(instance.Status)
}

// IsInstanceReady retruns true if a instance is ready; false otherwise.
func IsInstanceReadyConditionTrue(status InstanceStatus) bool {
	condition := GetInstanceReadyCondition(status)
	return condition != nil && condition.Status == api.ConditionTrue
}

// Extracts the instance ready condition from the given status and returns that.
// Returns nil if the condition is not present.
func GetInstanceReadyCondition(status InstanceStatus) *InstanceCondition {
	_, condition := GetInstanceCondition(&status, InstanceReady)
	return condition
}

// GetInstanceCondition extracts the provided condition from the given status and returns that.
// Returns nil and -1 if the condition is not present, and the index of the located condition.
func GetInstanceCondition(status *InstanceStatus, conditionType InstanceConditionType) (int, *InstanceCondition) {
	if status == nil {
		return -1, nil
	}
	for i := range status.Conditions {
		if status.Conditions[i].Type == conditionType {
			return i, &status.Conditions[i]
		}
	}
	return -1, nil
}

func InstanceStatusEqual(l, r InstanceStatus) bool {
	return l.Phase == r.Phase && l.PrivateIP == r.PrivateIP && l.PublicIP == r.PublicIP && l.InstanceID == r.InstanceID
}

func InstanceStatusDeepCopy(s *InstanceStatus) *InstanceStatus {
	ret := &InstanceStatus{
		Phase:             s.Phase,
		PrivateIP:         s.PrivateIP,
		PublicIP:          s.PublicIP,
		InstanceID:        s.InstanceID,
		CreationTimestamp: s.CreationTimestamp,
	}
	for _, c := range s.Conditions {
		s.Conditions = append(s.Conditions, c)
	}

	return ret
}
