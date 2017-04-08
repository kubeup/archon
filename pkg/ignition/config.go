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

type Config struct {
	Ignition  *Ignition  `yaml:"ignition,omitempty"`
	Storage   *Storage   `yaml:"storage,omitempty"`
	Systemd   *Systemd   `yaml:"systemd,omitempty"`
	Networkd  *Networkd  `yaml:"networkd,omitempty"`
	Passwd    *Passwd    `yaml:"passwd,omitempty"`
	Etcd      *Etcd      `yaml:"etcd,omitempty"`
	Flannel   *Flannel   `yaml:"flannel,omitempty"`
	Update    *Update    `yaml:"update,omitempty"`
	Docker    *Docker    `yaml:"docker,omitempty"`
	Locksmith *Locksmith `yaml:"locksmith,omitempty"`
}

type Ignition struct {
	Config IgnitionConfig `yaml:"config,omitempty"`
}

type IgnitionConfig struct {
	Append  []ConfigReference `yaml:"append,omitempty"`
	Replace *ConfigReference  `yaml:"replace,omitempty"`
}

type ConfigReference struct {
	Source       string       `yaml:"source,omitempty"`
	Verification Verification `yaml:"verification,omitempty"`
}
