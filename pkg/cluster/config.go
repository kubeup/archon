package cluster

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

type ConfigSpec struct {
	Name  string
	Items []*ConfigItem
}

type ConfigItem struct {
	Name  string
	Value string
}

type Config struct {
	unversioned.TypeMeta `json:",inline"`
	Metadata             api.ObjectMeta `json:"metadata"`
	Spec                 ConfigSpec
}
