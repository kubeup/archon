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

package render

import (
	"bytes"
	"github.com/Masterminds/sprig"
	"k8s.io/kubernetes/pkg/api"
	"kubeup.com/archon/pkg/cluster"
	"text/template"
)

const (
	SecretAliasKey = "archon.kubeup.com/alias"
)

type Renderer interface {
	Render(name, tpl string) (string, error)
}

type InstanceRenderer struct {
	Configs map[string]map[string]string
	Secrets map[string]map[string][]byte
	Network cluster.Network
	Status  cluster.InstanceStatus
	Meta    api.ObjectMeta
}

func (r *InstanceRenderer) Render(name, tpl string) (string, error) {
	t, err := template.New(name).Funcs(sprig.TxtFuncMap()).Parse(tpl)
	s := bytes.NewBufferString("")
	err = t.Execute(s, r)
	if err != nil {
		return "", err
	}
	return s.String(), nil
}

func NewInstanceRenderer(instance *cluster.Instance) (r *InstanceRenderer, err error) {
	r = &InstanceRenderer{
		Configs: make(map[string]map[string]string),
		Secrets: make(map[string]map[string][]byte),
		Network: instance.Dependency.Network,
		Status:  instance.Status,
		Meta:    instance.ObjectMeta,
	}

	for _, c := range instance.Spec.Configs {
		r.Configs[c.Name] = c.Data
	}

	for _, c := range instance.Dependency.Secrets {
		if alias := c.Annotations[SecretAliasKey]; alias != "" {
			r.Secrets[alias] = c.Data
		} else {
			r.Secrets[c.Name] = c.Data
		}
	}
	return
}
