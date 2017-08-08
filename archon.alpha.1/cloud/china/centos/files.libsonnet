local archon = import "archon.alpha.1/archon.libsonnet";
local file = archon.v1.instance.mixin.spec.filesType;

local runMaster = |||
  - sh
  - -c
  - "curl -s https://raw.githubusercontent.com/kubeup/okdc/master/okdc-centos.sh|NOINPUT=true TOKEN=%(token)s sh"
|||;

local runNode = |||
  - sh
  - -c
  - "curl -s https://raw.githubusercontent.com/kubeup/okdc/master/okdc-centos.sh|NOINPUT=true TOKEN=%(token)s MASTER=%(master-ip)s sh"
|||;

{
    master:: {
        i80kubeadm(config):: file.new() + file.name("kubeadm") + file.path("/config/runcmd/kubeadm") + file.content(runMaster % config.k8s),
    },
    node:: {
        i80kubeadm(config):: file.new() + file.name("kubeadm") + file.path("/config/runcmd/kubeadm") + file.content(runNode % config.k8s),

    },
}
