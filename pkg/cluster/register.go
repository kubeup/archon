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
	"k8s.io/apimachinery/pkg/runtime/schema"
)

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
