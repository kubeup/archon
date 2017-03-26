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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes, addDefaultingFuncs)
	AddToScheme   = SchemeBuilder.AddToScheme
	//Scheme        = runtime.NewScheme()
)

// GroupName is the group name use in this package
const GroupName = "archon.kubeup.com"
const GroupVersion = "v1"

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: GroupVersion}

// Adds the list of known types to api.Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	RegisterDefaults(scheme)
	scheme.AddKnownTypes(SchemeGroupVersion,
		&User{},
		&UserList{},
		&Instance{},
		&InstanceList{},
		&InstanceGroup{},
		&InstanceGroupList{},
		&Network{},
		&NetworkList{},
		&ReservedInstance{},
		&ReservedInstanceList{},
	//	&api.ListOptions{},
	//		&api.DeleteOptions{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}

func (i User) Get() schema.ObjectKind {
	return &i.TypeMeta
}

func (i User) GetObjectMeta() metav1.Object {
	return &i.ObjectMeta
}

func (il UserList) Get() schema.ObjectKind {
	return &il.TypeMeta
}

func (il UserList) GetObjectMeta() metav1.List {
	return &il.ListMeta
}

func (i Network) Get() schema.ObjectKind {
	return &i.TypeMeta
}

func (i Network) GetObjectMeta() metav1.Object {
	return &i.ObjectMeta
}

func (il NetworkList) Get() schema.ObjectKind {
	return &il.TypeMeta
}

func (il NetworkList) GetObjectMeta() metav1.List {
	return &il.ListMeta
}

func (i InstanceGroup) Get() schema.ObjectKind {
	return &i.TypeMeta
}

func (i InstanceGroup) GetObjectMeta() metav1.Object {
	return &i.ObjectMeta
}

func (il InstanceGroupList) Get() schema.ObjectKind {
	return &il.TypeMeta
}

func (il InstanceGroupList) GetObjectMeta() metav1.List {
	return &il.ListMeta
}

func (i Instance) Get() schema.ObjectKind {
	return &i.TypeMeta
}

func (i Instance) GetObjectMeta() metav1.Object {
	return &i.ObjectMeta
}

func (il InstanceList) Get() schema.ObjectKind {
	return &il.TypeMeta
}

func (il InstanceList) GetObjectMeta() metav1.List {
	return &il.ListMeta
}

func (i ReservedInstance) Get() schema.ObjectKind {
	return &i.TypeMeta
}

func (i ReservedInstance) GetObjectMeta() metav1.Object {
	return &i.ObjectMeta
}

func (il ReservedInstanceList) Get() schema.ObjectKind {
	return &il.TypeMeta
}

func (il ReservedInstanceList) GetObjectMeta() metav1.List {
	return &il.ListMeta
}
