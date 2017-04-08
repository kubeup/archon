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

package userdata

import (
	"strconv"
	"strings"

	"github.com/coreos/yaml"
	"kubeup.com/archon/pkg/cluster"
	"kubeup.com/archon/pkg/ignition"
	"kubeup.com/archon/pkg/render"
)

func GenerateCoreOSIgnition(instance *cluster.Instance) ([]byte, error) {
	result := &ignition.Config{}
	renderer, err := render.NewInstanceRenderer(instance)
	if err != nil {
		return nil, err
	}

	for _, t := range instance.Spec.Files {
		f := ignition.File{
			Filesystem: t.Filesystem,
			Path:       t.Path,
		}
		if t.UserId != 0 {
			f.User = &ignition.FileUser{t.UserId}
		}
		if t.GroupId != 0 {
			f.Group = &ignition.FileGroup{t.GroupId}
		}
		if t.RawFilePermissions != "" {
			m, err := strconv.ParseUint(t.RawFilePermissions, 0, 32)
			if err != nil {
				return nil, err
			}
			f.Mode = int(m)
		}
		if t.Content == "" {
			f.Contents.Inline, err = renderer.Render(t.Name, t.Template)
			if err != nil {
				return nil, err
			}
		} else {
			f.Contents.Inline = t.Content
		}
		if strings.HasPrefix(f.Path, "/coreos/unit/") {
			u := ignition.SystemdUnit{}
			err = yaml.Unmarshal([]byte(f.Contents.Inline), &u)
			if err != nil {
				return nil, err
			}
			if result.Systemd == nil {
				result.Systemd = &ignition.Systemd{}
			}
			result.Systemd.Units = append(result.Systemd.Units, u)
		} else if f.Path == "/coreos/update" {
			u := ignition.Update{}
			err = yaml.Unmarshal([]byte(f.Contents.Inline), &u)
			if err != nil {
				return nil, err
			}
			result.Update = &u
		} else if f.Path == "/coreos/etcd" {
			u := ignition.Etcd{}
			err = yaml.Unmarshal([]byte(f.Contents.Inline), &u)
			if err != nil {
				return nil, err
			}
			result.Etcd = &u
		} else if f.Path == "/coreos/flannel" {
			u := ignition.Flannel{}
			err = yaml.Unmarshal([]byte(f.Contents.Inline), &u)
			if err != nil {
				return nil, err
			}
			result.Flannel = &u
		} else {
			if result.Storage == nil {
				result.Storage = &ignition.Storage{}
			}
			result.Storage.Files = append(result.Storage.Files, f)
		}

	}

	for _, t := range instance.Dependency.Users {
		u := ignition.User{
			Name:              t.Spec.Name,
			PasswordHash:      t.Spec.PasswordHash,
			SSHAuthorizedKeys: t.Spec.SSHAuthorizedKeys,
		}
		if result.Passwd == nil {
			result.Passwd = &ignition.Passwd{}
		}
		result.Passwd.Users = append(result.Passwd.Users, u)
	}

	data, err := yaml.Marshal(result)
	if err != nil {
		return nil, err
	}
	return append([]byte("---\n"), data...), nil
}
