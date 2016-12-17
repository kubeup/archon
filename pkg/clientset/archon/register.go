package archon

import (
	"kubeup.com/archon/pkg/cluster"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/runtime"
)

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
	Scheme        = runtime.NewScheme()
)

// GroupName is the group name use in this package
const GroupName = "archon.kubeup.com"
const GroupVersion = "v1"

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = unversioned.GroupVersion{Group: GroupName, Version: GroupVersion}

// Adds the list of known types to api.Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&cluster.Instance{},
		&cluster.InstanceList{},
		&cluster.Network{},
		&cluster.NetworkList{},
		&api.ListOptions{},
		&api.DeleteOptions{},
	)
	return nil
}
