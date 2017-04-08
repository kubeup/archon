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

type Disk struct {
	Device     string      `yaml:"device,omitempty"`
	WipeTable  bool        `yaml:"wipe_table,omitempty"`
	Partitions []Partition `yaml:"partitions,omitempty"`
}

type Partition struct {
	Label    string `yaml:"label,omitempty"`
	Number   int    `yaml:"number,omitempty"`
	Size     string `yaml:"size,omitempty"`
	Start    string `yaml:"start,omitempty"`
	TypeGUID string `yaml:"type_guid,omitempty"`
}
