local archon = import "archon.alpha.1/cloud/aws/centos/archon.libsonnet";
local config = import "config.libsonnet";

local mixin =  {
    network+:: {
        config+:: {
            region:: config.networkRegion,
            zone:: config.networkZone,
            subnet:: config.networkSubnet,
        }
    },
    user+:: {
        new(name, config={})::
            local finalConfig = self.config + config;
            super.new(name, finalConfig) + self.mixin.spec.sshAuthorizedKeys(finalConfig.sshAuthorizedKeys),
        config+:: {
            sshAuthorizedKeys:: config.sshAuthorizedKeys,
        }
    },
    instanceGroup+:: {
        local file = archon.v1.instance.mixin.spec.filesType,
        local sysctl = |||
          - sysctl
          - -p
        |||,
        new(name, config={}):: super.new(name, config) + self.mixin.spec.template.spec.users({name: "k8s-user"}),
        files+:: {
            i01fixIptable(config):: file.new() + file.name("fix-iptable") + file.path("/etc/sysctl.d/10-iptable.conf") + file.content("net.bridge.bridge-nf-call-iptables = 1"),
            i02sysctl(config):: file.new() + file.name("sysctl") + file.path("/config/runcmd/sysctl") + file.content(sysctl),
        },
    },
    master+:: {
        config+:: {
            k8sToken:: config.k8sToken,
            k8sPodIPRange:: config.k8sPodIPRange,
            networkName:: "k8s-net",
            instanceProfile:: config.masterInstanceProfile,
        }
    },
    node+:: {
        config+:: {
            k8sToken:: config.k8sToken,
            k8sMasterIP:: config.k8sMasterIP,
            networkName:: "k8s-net",
        }
    },
};

archon + {
    v1+:: archon.v1 + mixin,
}
