package client

import (
	"io"
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
	if conf.Os == "ubuntu" {
		return GetUbuntuCmds(conf)
	} else if conf.Os == "centos" {
		return GetUbuntuCmds(conf)
	} else if conf.Os == "coreos" {
		return GetCoreOSCmds(conf)
	}
	return []string{}
}
