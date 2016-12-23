package cluster

import (
	"encoding/json"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/meta"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

type InstanceGroup struct {
	unversioned.TypeMeta `json:",inline"`
	api.ObjectMeta       `json:"metadata"`
	Spec                 InstanceGroupSpec   `json:"spec,omitempty"`
	Status               InstanceGroupStatus `json:"status,omitempty"`
}
type InstanceGroupSpec struct {
	Replicas int32

	// Minimum number of seconds for which a newly created pod should be ready
	// without any of its container crashing, for it to be considered available.
	// Defaults to 0 (pod will be considered available as soon as it is ready)
	// +optional
	MinReadySeconds int32                      `json:"minReadySeconds,omitempty"`
	Selector        *unversioned.LabelSelector `json:"selector,omitempty"`
	Template        InstanceTemplateSpec       `json:"template,omitempty"`
}

type InstanceGroupConditionType string

// These are valid conditions of a replica set.
const (
	// InstanceGroupReplicaFailure is added in a replica set when one of its pods fails to be created
	// due to insufficient quota, limit ranges, pod security policy, node selectors, etc. or deleted
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
	LastTransitionTime unversioned.Time `json:"lastTransitionTime,omitempty"`
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

	// The number of pods that have labels matching the labels of the pod template of the replicaset.
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
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata"`

	Items []*InstanceGroup `json:"items"`
}

func (i InstanceGroup) GetObjectKind() unversioned.ObjectKind {
	return &i.TypeMeta
}

func (i InstanceGroup) GetObjectMeta() meta.Object {
	return &i.ObjectMeta
}

func (il InstanceGroupList) GetObjectKind() unversioned.ObjectKind {
	return &il.TypeMeta
}

func (il InstanceGroupList) GetObjectMeta() unversioned.List {
	return &il.ListMeta
}

// The code below is used only to work around a known problem with third-party
// resources and ugorji. If/when these issues are resolved, the code below
// should no longer be required.

type InstanceGroupCopy InstanceGroup
type InstanceGroupListCopy InstanceGroupList

func (e *InstanceGroup) UnmarshalJSON(data []byte) error {
	tmp := InstanceGroupCopy{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	tmp2 := InstanceGroup(tmp)
	*e = tmp2
	return nil
}

func (el *InstanceGroupList) UnmarshalJSON(data []byte) error {
	tmp := InstanceGroupListCopy{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	tmp2 := InstanceGroupList(tmp)
	*el = tmp2
	return nil
}

func (e *InstanceGroup) CodecEncodeSelf() {
}

func (e *InstanceGroup) CodecDecodeSelf() {
}

func (el *InstanceGroupList) CodecEncodeSelf() {
}

func (el *InstanceGroupList) CodecDecodeSelf() {
}
