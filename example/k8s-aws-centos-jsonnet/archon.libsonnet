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
        new(name):: super.new(name) + self.mixin.spec.sshAuthorizedKeys(config.sshAuthorizedKeys),
    },
    instanceGroup+:: {
        local file = archon.v1.instance.mixin.spec.filesType,
        local sysctl = |||
          - sysctl
          - -p
        |||,
        new(name):: super.new(name) + self.mixin.spec.template.spec.users({name: "k8s-user"}),
        files+:: {
            i01fixIptable(config):: file.new() + file.name("fix-iptable") + file.path("/etc/sysctl.d/10-iptable.conf") + file.content("net.bridge.bridge-nf-call-iptables = 1"),
            i02sysctl(config):: file.new() + file.name("sysctl") + file.path("/config/runcmd/sysctl") + file.content(sysctl),
        },
    },
    master+:: {
        config+:: {
            k8s+:: {
                "token": config.token,
                "pod-ip-range": config.podIPRange,
            },
            networkName:: "k8s-net",
            instanceProfile:: config.masterInstanceProfile,
        }
    },
    node+:: {
        config+:: {
            k8s+:: {
                "token":: config.token,
                "master-ip":: config.masterIP,
            },
            networkName:: "k8s-net",
        }
    },
};

archon + {
    v1+:: archon.v1 + mixin,
}
