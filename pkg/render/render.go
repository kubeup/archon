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
