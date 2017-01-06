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
	"encoding/json"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/meta"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

type UserSpec struct {
	Name              string   `yaml:"name"`
	PasswordHash      string   `yaml:"passwd"`
	SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys"`
}

type User struct {
	unversioned.TypeMeta `json:",inline"`
	api.ObjectMeta       `json:"metadata"`
	Spec                 UserSpec
}

type UserList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata"`

	Items []*User `json:"items"`
}

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

// The code below is used only to work around a known problem with third-party
// resources and ugorji. If/when these issues are resolved, the code below
// should no longer be required.

type UserCopy User
type UserListCopy UserList

func (e *User) UnmarshalJSON(data []byte) error {
	tmp := UserCopy{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	tmp2 := User(tmp)
	*e = tmp2
	return nil
}

func (el *UserList) UnmarshalJSON(data []byte) error {
	tmp := UserListCopy{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	tmp2 := UserList(tmp)
	*el = tmp2
	return nil
}

func (e *User) CodecEncodeSelf() {
}

func (e *User) CodecDecodeSelf() {
}

func (el *UserList) CodecEncodeSelf() {
}

func (el *UserList) CodecDecodeSelf() {
}
