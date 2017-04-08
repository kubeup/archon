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

type File struct {
	Filesystem string       `yaml:"filesystem,omitempty"`
	Path       string       `yaml:"path,omitempty"`
	Mode       int          `yaml:"mode,omitempty"`
	Contents   FileContents `yaml:"contents,omitempty"`
	User       *FileUser    `yaml:"user,omitempty"`
	Group      *FileGroup   `yaml:"group,omitempty"`
}

type FileContents struct {
	Remote *Remote `yaml:"remote,omitempty"`
	Inline string  `yaml:"inline,omitempty"`
}

type Remote struct {
	Url          string       `yaml:"url,omitempty"`
	Compression  string       `yaml:"compression,omitempty"`
	Verification Verification `yaml:"verification,omitempty"`
}

type FileUser struct {
	Id int `yaml:"id,omitempty"`
}

type FileGroup struct {
	Id int `yaml:"id,omitempty"`
}
