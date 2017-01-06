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

type UbuntuConfig struct {
	AptSources     []AptSource `yaml:"apt_sources,omitempty"`
	Packages       []string    `yaml:"packages,omitempty"`
	Runcmd         [][]string  `yaml:"runcmd,omitempty"`
	WriteFiles     []File      `yaml:"write_files,omitempty"`
	Hostname       string      `yaml:"hostname,omitempty"`
	ManageEtcHosts string      `yaml:"manage_etc_hosts,omitempty"`
}

type AptSource struct {
	Source string `yaml:"source,omitempty"`
	Key    string `yaml:"key,omitempty"`
}

func (uc UbuntuConfig) Bytes() ([]byte, error) {
	data, err := yaml.Marshal(uc)
	if err != nil {
		return nil, err
	}
	return append([]byte("#cloud-config\n"), data...), nil
}

func (uc UbuntuConfig) String() (string, error) {
	data, err := uc.Bytes()
	if err != nil {
		return "", err
	}
	return string(data), nil
}
