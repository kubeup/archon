package client

import (
	"fmt"
	"io"
	"strings"
)

type Config struct {
	Hostname                 string
	Port                     int
	User                     string
	Password                 string
	PublicKeyFile            string
	Cmds                     []string
	Os                       string
	Server                   string
	DontCleanCloudInitStatus bool
	UseCloudDataSource       bool
	UserData                 string
	Callback                 string
	Stdout                   io.Writer
}

func (conf *Config) GetCmds() []string {
	if len(conf.Cmds) > 0 {
		return conf.Cmds
	}
	os := strings.ToLower(conf.Os)
	if os == "ubuntu" {
		return GetUbuntuCmds(conf)
	} else if os == "centos" {
		return GetCentOSCmds(conf)
	} else if os == "coreos" {
		return GetCoreOSCmds(conf)
	} else {
		if conf.Stdout != nil {
			fmt.Fprintf(conf.Stdout, "Warning: unsupported system for ssh-cloudinit: %s", conf.Os)
		}

	}
	return []string{}
}
