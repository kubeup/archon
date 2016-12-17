package cloudinit

import (
	"github.com/coreos/yaml"
)

type UbuntuConfig struct {
	AptSources     []AptSource `yaml:"apt_sources,omitempty"`
	Packages       []string    `yaml:"packages,omitempty"`
	Runcmd         [][]string  `yaml:"runcmd,omitempty"`
	WriteFiles     []File      `yaml:"write_files,omitempty"`
	Hostname       string      `yaml:"hostname,omitempty"`
	ManageEtcHosts string      `yaml:"manage_etc_hosts,omitempty"`
}

type AptSource struct {
	Source string `yaml:"source,omitempty"`
	Key    string `yaml:"key,omitempty"`
}

func (uc UbuntuConfig) Bytes() ([]byte, error) {
	data, err := yaml.Marshal(uc)
	if err != nil {
		return nil, err
	}
	return append([]byte("#cloud-config\n"), data...), nil
}

func (uc UbuntuConfig) String() (string, error) {
	data, err := uc.Bytes()
	if err != nil {
		return "", err
	}
	return string(data), nil
}
