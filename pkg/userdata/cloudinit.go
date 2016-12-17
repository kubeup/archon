package userdata

import (
	"bytes"
	"github.com/Masterminds/sprig"
	"github.com/coreos/yaml"
	"kubeup.com/archon/cloudinit"
	"kubeup.com/archon/pkg/cluster"
	"strings"
	"text/template"
)

type Renderer interface {
	Render(file *cluster.FileSpec) (string, error)
}

type tmplRender struct {
	Configs map[string]cluster.ConfigSpec
}

func (r *tmplRender) Render(file *cluster.FileSpec) (string, error) {
	t, err := template.New(file.Name).Funcs(sprig.TxtFuncMap()).Parse(file.Template)
	s := bytes.NewBufferString("")
	err = t.Execute(s, r.Configs)
	if err != nil {
		return "", err
	}
	return s.String(), nil
}

func newRenderer(instance *cluster.Instance) (r *tmplRender, err error) {
	r = &tmplRender{
		Configs: make(map[string]cluster.ConfigSpec),
	}

	for _, c := range instance.Spec.Configs {
		r.Configs[c.Name] = c
	}
	return
}

func Generate(instance *cluster.Instance) ([]byte, error) {
	result := &cloudinit.CloudConfig{}
	renderer, err := newRenderer(instance)
	if err != nil {
		return nil, err
	}

	files := make([]cloudinit.File, 0)
	units := make([]cloudinit.Unit, 0)
	for _, t := range instance.Spec.Files {
		f := cloudinit.File{
			Encoding:           t.Encoding,
			Owner:              t.Owner,
			Path:               t.Path,
			RawFilePermissions: t.RawFilePermissions,
		}
		if t.Content == "" {
			f.Content, err = renderer.Render(&t)
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
			units = append(units, u)
		} else {
			files = append(files, f)
		}

	}
	result.WriteFiles = files
	result.CoreOS = &cloudinit.CoreOS{
		Units: units,
	}

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

	return result.Bytes()
}
