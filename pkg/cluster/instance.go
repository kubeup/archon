package cluster

import (
	"encoding/json"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/meta"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

type InstanceTemplateSpec struct {
	api.ObjectMeta `json:"metadata"`
	Spec           InstanceSpec
}

type InstanceDependency struct {
	Network Network
	Users   []User
}

type Instance struct {
	unversioned.TypeMeta `json:",inline"`
	api.ObjectMeta       `json:"metadata"`
	Spec                 InstanceSpec
	Status               InstanceStatus
	Dependency           InstanceDependency `json:"-"`
}

type InstanceSpec struct {
	Image        string
	InstanceType string
	NetworkName  string
	Files        []FileSpec
	// Secrets      []Secret
	Configs  []ConfigSpec
	Users    []api.LocalObjectReference
	Hostname string
}

type InstanceStatus struct {
	Phase             InstancePhase
	PrivateIP         string
	PublicIP          string
	InstanceID        string
	CreationTimestamp unversioned.Time `json:"creationTimestamp,omitempty" protobuf:"bytes,8,opt,name=creationTimestamp"`
}

type InstancePhase string

const (
	InstancePending InstancePhase = "Pending"
	InstanceRunning InstancePhase = "Running"
	InstanceFailed  InstancePhase = "Failed"
	InstanceUnknown InstancePhase = "Unknown"
)

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
