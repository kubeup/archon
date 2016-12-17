package cluster

import (
	"encoding/json"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/meta"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

type NetworkSpec struct {
	Region string `k8s:"region"`
	Zone   string `k8s:"zone"`
	Subnet string `k8s:"subnet"`
}

type NetworkStatus struct {
	Phase NetworkPhase
}

type NetworkPhase string

const (
	NetworkPending NetworkPhase = "Pending"
	NetworkRunning NetworkPhase = "Running"
	NetworkFailed  NetworkPhase = "Failed"
	NetworkUnknown NetworkPhase = "Unknown"
)

type Network struct {
	unversioned.TypeMeta `json:",inline"`
	api.ObjectMeta       `json:"metadata"`
	Spec                 NetworkSpec
	Status               NetworkStatus
}

type NetworkList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata"`

	Items []*Network `json:"items"`
}

func (i Network) GetObjectKind() unversioned.ObjectKind {
	return &i.TypeMeta
}

func (i Network) GetObjectMeta() meta.Object {
	return &i.ObjectMeta
}

func (il NetworkList) GetObjectKind() unversioned.ObjectKind {
	return &il.TypeMeta
}

func (il NetworkList) GetObjectMeta() unversioned.List {
	return &il.ListMeta
}

// The code below is used only to work around a known problem with third-party
// resources and ugorji. If/when these issues are resolved, the code below
// should no longer be required.

type NetworkCopy Network
type NetworkListCopy NetworkList

func (e *Network) UnmarshalJSON(data []byte) error {
	tmp := NetworkCopy{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	tmp2 := Network(tmp)
	*e = tmp2
	return nil
}

func (el *NetworkList) UnmarshalJSON(data []byte) error {
	tmp := NetworkListCopy{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	tmp2 := NetworkList(tmp)
	*el = tmp2
	return nil
}

func (e *Network) CodecEncodeSelf() {
}

func (e *Network) CodecDecodeSelf() {
}

func (el *NetworkList) CodecEncodeSelf() {
}

func (el *NetworkList) CodecDecodeSelf() {
}
