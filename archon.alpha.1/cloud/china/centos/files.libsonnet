local archon = import "archon.alpha.1/archon.libsonnet";
local file = archon.v1.instance.mixin.spec.filesType;

local runMaster = |||
  - sh
  - -c
  - "curl -s https://raw.githubusercontent.com/kubeup/okdc/master/okdc-centos.sh|NOINPUT=true TOKEN=%(k8sToken)s sh"
|||;

local runNode = |||
  - sh
  - -c
  - "curl -s https://raw.githubusercontent.com/kubeup/okdc/master/okdc-centos.sh|NOINPUT=true TOKEN=%(k8sToken)s MASTER=%(k8sMasterIP)s:6443 sh"
|||;

{
    master:: {
        i80kubeadm(config)::
            if config.k8sToken == "" then
                error "config k8sToken should not be empty"
            else
                file.new() + file.name("kubeadm") + file.path("/config/runcmd/kubeadm") + file.content(runMaster % config),
    },
    node:: {
        i80kubeadm(config)::
            if config.k8sToken == "" then
                error "config k8sToken should not be empty"
            else if config.k8sMasterIP == "" then
                error "config k8sMasterIP should not be empty"
            else
                file.new() + file.name("kubeadm") + file.path("/config/runcmd/kubeadm") + file.content(runNode % config),

    },
}
