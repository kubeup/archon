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

package cloudinit

import (
	"github.com/coreos/yaml"
)

type CloudConfig struct {
	Apt           Apt      `yaml:"apt,omitempty"`
	Packages      []string `yaml:"packages,omitempty"`
	PackageUpdate bool     `yaml:"package_update,omitempty"`
	// Runcmd could be either []string or string
	Runcmd         []interface{}      `yaml:"runcmd,omitempty"`
	WriteFiles     []File             `yaml:"write_files,omitempty"`
	Hostname       string             `yaml:"hostname,omitempty"`
	Users          []User             `yaml:"users,omitempty"`
	ManageEtcHosts string             `yaml:"manage_etc_hosts,omitempty"`
	YumRepos       map[string]YumRepo `yaml:"yum_repos,omitempty"`
}

type Apt struct {
	Sources             map[string]AptSource `yaml:"sources,omitempty"`
	Primary             []AptMirror          `yaml:"primary,omitempty"`
	PreserveSourcesList bool                 `yaml:"preserve_sources_list,omitempty"`
}

type AptMirror struct {
	Arches []string `yaml:"arches,omitempty"`
	URI    string   `yaml:"uri,omitempty"`
}

type AptSource struct {
	Source string `yaml:"source,omitempty"`
	Key    string `yaml:"key,omitempty"`
}

type YumRepo struct {
	Name     string `yaml:"name,omitempty"`
	BaseUrl  string `yaml:"baseurl,omitempty"`
	Enabled  bool   `yaml:"enabled,omitempty"`
	GPGCheck bool   `yaml:"gpgcheck,omitempty"`
	GPGKey   string `yaml:"gpgkey,omitempty"`
}

func (uc CloudConfig) Bytes() ([]byte, error) {
	data, err := yaml.Marshal(uc)
	if err != nil {
		return nil, err
	}
	return append([]byte("#cloud-config\n"), data...), nil
}

func (uc CloudConfig) String() (string, error) {
	data, err := uc.Bytes()
	if err != nil {
		return "", err
	}
	return string(data), nil
}
