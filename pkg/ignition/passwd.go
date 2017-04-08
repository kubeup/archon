// Copyright 2016 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ignition

type Passwd struct {
	Users  []User  `yaml:"users,omitempty"`
	Groups []Group `yaml:"groups,omitempty"`
}

type User struct {
	Name              string      `yaml:"name,omitempty"`
	PasswordHash      string      `yaml:"password_hash,omitempty"`
	SSHAuthorizedKeys []string    `yaml:"ssh_authorized_keys,omitempty"`
	Create            *UserCreate `yaml:"create,omitempty"`
}

type UserCreate struct {
	Uid          *uint    `yaml:"uid,omitempty"`
	GECOS        string   `yaml:"gecos,omitempty"`
	Homedir      string   `yaml:"home_dir,omitempty"`
	NoCreateHome bool     `yaml:"no_create_home,omitempty"`
	PrimaryGroup string   `yaml:"primary_group,omitempty"`
	Groups       []string `yaml:"groups,omitempty"`
	NoUserGroup  bool     `yaml:"no_user_group,omitempty"`
	System       bool     `yaml:"system,omitempty"`
	NoLogInit    bool     `yaml:"no_log_init,omitempty"`
	Shell        string   `yaml:"shell,omitempty"`
}

type Group struct {
	Name         string `yaml:"name,omitempty"`
	Gid          *uint  `yaml:"gid,omitempty"`
	PasswordHash string `yaml:"password_hash,omitempty"`
	System       bool   `yaml:"system,omitempty"`
}
