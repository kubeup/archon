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

type Systemd struct {
	Units []SystemdUnit `yaml:"units,omitempty"`
}

type SystemdUnit struct {
	Name     string              `yaml:"name,omitempty"`
	Enable   bool                `yaml:"enable,omitempty"`
	Mask     bool                `yaml:"mask,omitempty"`
	Contents string              `yaml:"contents,omitempty"`
	DropIns  []SystemdUnitDropIn `yaml:"dropins,omitempty"`
}

type SystemdUnitDropIn struct {
	Name     string `yaml:"name,omitempty"`
	Contents string `yaml:"contents,omitempty"`
}
