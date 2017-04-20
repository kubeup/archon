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
	"github.com/coreos/yaml"
	"kubeup.com/archon/pkg/cloudinit"
	"kubeup.com/archon/pkg/cluster"
	"kubeup.com/archon/pkg/render"
	"strings"
)

func Generate(instance *cluster.Instance) ([]byte, error) {
	switch instance.Spec.OS {
	case "CoreOS":
		return GenerateCoreOSCloudConfig(instance)
	case "CoreOSIgnition":
		return GenerateCoreOSIgnition(instance)
	default:
		return GenerateCloudConfig(instance)
	}
}

func GenerateCoreOSCloudConfig(instance *cluster.Instance) ([]byte, error) {
	result := &cloudinit.CoreOSCloudConfig{}
	renderer, err := render.NewInstanceRenderer(instance)
	if err != nil {
		return nil, err
	}

	coreos := &cloudinit.CoreOS{}
	files := make([]cloudinit.File, 0)
	for _, t := range instance.Spec.Files {
		f := cloudinit.File{
			Encoding:           t.Encoding,
			Owner:              t.Owner,
			Path:               t.Path,
			RawFilePermissions: t.RawFilePermissions,
		}
		if t.Content == "" {
			f.Content, err = renderer.Render(t.Name, t.Template)
			if err != nil {
				return nil, err
			}
		} else {
			f.Content = t.Content
		}
		if strings.HasPrefix(f.Path, "/coreos/unit/") {
			u := cloudinit.Unit{}
			err = yaml.Unmarshal([]byte(f.Content), &u)
			if err != nil {
				return nil, err
			}
			coreos.Units = append(coreos.Units, u)
		} else if f.Path == "/coreos/update" {
			u := cloudinit.Update{}
			err = yaml.Unmarshal([]byte(f.Content), &u)
			if err != nil {
				return nil, err
			}
			coreos.Update = &u
		} else if f.Path == "/coreos/etcd2" {
			u := cloudinit.Etcd2{}
			err = yaml.Unmarshal([]byte(f.Content), &u)
			if err != nil {
				return nil, err
			}
			coreos.Etcd2 = &u
		} else if f.Path == "/coreos/flannel" {
			u := cloudinit.Flannel{}
			err = yaml.Unmarshal([]byte(f.Content), &u)
			if err != nil {
				return nil, err
			}
			coreos.Flannel = &u
		} else {
			files = append(files, f)
		}

	}
	result.WriteFiles = files
	result.CoreOS = coreos

	users := make([]cloudinit.User, 0)
	for _, t := range instance.Dependency.Users {
		u := cloudinit.User{
			Name:              t.Spec.Name,
			PasswordHash:      t.Spec.PasswordHash,
			SSHAuthorizedKeys: t.Spec.SSHAuthorizedKeys,
		}
		users = append(users, u)
	}
	result.Users = users

	if instance.Spec.Hostname != "" {
		result.Hostname, err = renderer.Render("hostname", instance.Spec.Hostname)
		if err != nil {
			return nil, err
		}
	}

	return result.Bytes()
}

func GenerateCloudConfig(instance *cluster.Instance) ([]byte, error) {
	result := &cloudinit.CloudConfig{}
	renderer, err := render.NewInstanceRenderer(instance)
	if err != nil {
		return nil, err
	}

	files := make([]cloudinit.File, 0)
	for _, t := range instance.Spec.Files {
		f := cloudinit.File{
			Encoding:           t.Encoding,
			Owner:              t.Owner,
			Path:               t.Path,
			RawFilePermissions: t.RawFilePermissions,
		}
		if t.Content == "" {
			f.Content, err = renderer.Render(t.Name, t.Template)
			if err != nil {
				return nil, err
			}
		} else {
			f.Content = t.Content
		}
		if strings.HasPrefix(f.Path, "/config/runcmd/") {
			c := []string{}
			err = yaml.Unmarshal([]byte(f.Content), &c)
			if err != nil {
				return nil, err
			}
			if len(c) == 1 {
				result.Runcmd = append(result.Runcmd, c[0])
			} else {
				result.Runcmd = append(result.Runcmd, c)
			}
		} else if f.Path == "/config/apt" {
			c := cloudinit.Apt{}
			err = yaml.Unmarshal([]byte(f.Content), &c)
			if err != nil {
				return nil, err
			}
			result.Apt = c
		} else if f.Path == "/config/yum_repos" {
			c := make(map[string]cloudinit.YumRepo)
			err = yaml.Unmarshal([]byte(f.Content), &c)
			if err != nil {
				return nil, err
			}
			result.YumRepos = c
		} else if f.Path == "/config/packages" {
			c := []string{}
			err = yaml.Unmarshal([]byte(f.Content), &c)
			if err != nil {
				return nil, err
			}
			result.Packages = c
			result.PackageUpdate = true
		} else {
			files = append(files, f)
		}

	}
	result.WriteFiles = files

	users := make([]cloudinit.User, 0)
	for _, t := range instance.Dependency.Users {
		u := cloudinit.User{
			Name:              t.Spec.Name,
			PasswordHash:      t.Spec.PasswordHash,
			SSHAuthorizedKeys: t.Spec.SSHAuthorizedKeys,
			Sudo:              t.Spec.Sudo,
			Shell:             t.Spec.Shell,
		}
		users = append(users, u)
	}
	result.Users = users

	return result.Bytes()
}
