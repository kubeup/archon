package userdata

import (
	"github.com/coreos/yaml"
	"kubeup.com/archon/cloudinit"
	"kubeup.com/archon/pkg/cluster"
	"kubeup.com/archon/pkg/render"
	"strings"
)

func Generate(instance *cluster.Instance) ([]byte, error) {
	result := &cloudinit.CloudConfig{}
	renderer, err := render.NewInstanceRenderer(instance)
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
