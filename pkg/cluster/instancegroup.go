package cluster

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

type InstanceGroup struct {
	unversioned.TypeMeta `json:",inline"`
	Metadata             api.ObjectMeta `json:"metadata"`
	Spec                 InstanceGroupSpec
	Status               InstanceGroupStatus
}
type InstanceGroupSpec struct {
	Size     int32
	Selector map[string]string `json:"selector,omitempty" protobuf:"bytes,2,rep,name=selector"`
	Template *InstanceTemplateSpec
}

type InstanceGroupStatus struct {
}
