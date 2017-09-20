package client

import (
	"fmt"
	"path"
	"strings"
)

const CloudDataPath = "/opt/cloud/"

func GetCentOSCmds(conf *Config) []string {
	s := conf.Server
	flags := ""
	source := ""
	if conf.UserData != "" {
		s = CloudDataPath
	} else if !strings.HasSuffix(s, "/") {
		s = s + "/"
	}
	cmds := []string{
		"sudo yum makecache",
		"sudo yum -y install cloud-init curl wget",
	}

	if conf.UseCloudDataSource == false {
		flags = "-l"
		cmds = append(cmds, fmt.Sprintf("sudo echo 'apt_preserve_sources_list: true\ncloud_init_modules: [write-files, update_etc_hosts, users-groups]\ncloud_final_modules: [scripts-vendor, scripts-per-once, scripts-per-boot, scripts-per-instance, scripts-user]\nusers: []\ndatasource_list: [NoCloud]\ndatasource: \n  NoCloud: \n    seedfrom: %s' > /etc/cloud/cloud.cfg.d/95_nocloud.cfg", s))
	} else {
		source = fmt.Sprintf("-f %s", path.Join(CloudDataPath, "user-data"))
	}

	cmds = append(cmds,
		fmt.Sprintf("sudo cloud-init %s init %s", source, flags),
		fmt.Sprintf("sudo cloud-init %s modules -m config", source),
		fmt.Sprintf("sudo cloud-init %s modules -m final", source),
	)

	if conf.DontCleanCloudInitStatus == false {
		cmds = append([]string{"sudo rm -rf /var/lib/cloud/instance/*"}, cmds...)
	}

	return cmds
}
func GetUbuntuCmds(conf *Config) []string {
	s := conf.Server
	flags := ""
	source := ""
	if conf.UserData != "" {
		s = CloudDataPath
	} else if !strings.HasSuffix(s, "/") {
		s = s + "/"
	}
	cmds := []string{
		"sudo apt-get update -q",
		"sudo sh -c 'which cloud-init || apt-get -y -q install cloud-init'",
		"sudo apt-get -y -q install curl wget",
	}

	if conf.UseCloudDataSource == false {
		flags = "-l"
		cmds = append(cmds, fmt.Sprintf("sudo sh -c \"echo 'apt_preserve_sources_list: true\ncloud_init_modules: [write-files, update_etc_hosts, users-groups]\ncloud_final_modules: [scripts-vendor, scripts-per-once, scripts-per-boot, scripts-per-instance, scripts-user]\nusers: []\ndatasource_list: [NoCloud]\ndatasource: \n  NoCloud: \n    seedfrom: %s' > /etc/cloud/cloud.cfg.d/96_nocloud.cfg\"", s))
	} else {
		source = fmt.Sprintf("-f %s", path.Join(CloudDataPath, "user-data"))
	}

	cmds = append(cmds,
		fmt.Sprintf("sudo cloud-init %s init %s", source, flags),
		fmt.Sprintf("sudo cloud-init %s modules -m config", source),
		fmt.Sprintf("sudo cloud-init %s modules -m final", source),
	)

	if conf.DontCleanCloudInitStatus == false {
		cmds = append([]string{"sudo rm -rf /var/lib/cloud/instance/*"}, cmds...)
	}
	return cmds
}

func GetCoreOSCmds(conf *Config) []string {
	filePath := path.Join(CloudDataPath, "user-data")
	if conf.UserData != "" {
		return []string{
			fmt.Sprintf("sudo coreos-cloudinit --from-file=%s", filePath),
		}
	}

	return []string{
		fmt.Sprintf("sudo coreos-cloudinit --from-url=%s", conf.Server),
	}
}
