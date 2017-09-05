package client

import (
	"fmt"
	"strings"
)

func GetUbuntuCmds(conf *Config) []string {
	s := conf.Server
	if !strings.HasSuffix(s, "/") {
		s = s + "/"
	}
	return []string{
		"sudo apt-get update",
		"sudo apt-get -y install cloud-init curl wget",
		fmt.Sprintf("sudo echo 'apt_preserve_sources_list: true\ncloud_init_modules: [write-files, update_etc_hosts, users-groups]\ncloud_final_modules: [scripts-vendor, scripts-per-once, scripts-per-boot, scripts-per-instance, scripts-user]\nusers: []\ndatasource_list: [NoCloud]\ndatasource: \n  NoCloud: \n    seedfrom: %s' > /etc/cloud/cloud.cfg.d/95_nocloud.cfg", s),
		"sudo rm -rf /var/lib/cloud/instance/*",
		"sudo cloud-init init",
		"sudo cloud-init modules -m config",
		"sudo cloud-init modules -m final",
	}
}

func GetCoreOSCmds(conf *Config) []string {
	return []string{
		fmt.Sprintf("sudo coreos-cloudinit --from-url=%s", conf.Server),
	}
}
