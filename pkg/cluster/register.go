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
	"k8s.io/kubernetes/pkg/api/meta"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

func (i User) GetObjectKind() unversioned.ObjectKind {
	return &i.TypeMeta
}

func (i User) GetObjectMeta() meta.Object {
	return &i.ObjectMeta
}

func (il UserList) GetObjectKind() unversioned.ObjectKind {
	return &il.TypeMeta
}

func (il UserList) GetObjectMeta() unversioned.List {
	return &il.ListMeta
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
