local archon = import "archon.alpha.1/cloud/ucloud/centos/archon.libsonnet";
local config = import "config.libsonnet";

local mixin =  {
    password+:: {
        new(name)::
            super.new(name, config.instancePassword),
    },
    user+:: {
        new(name):: super.new(name) + self.mixin.spec.sshAuthorizedKeys(config.sshAuthorizedKeys),
    },
    instanceGroup+:: {
        new(name):: super.new(name) + self.mixin.spec.template.spec.users({name: "k8s-user"}),
    },
    network+:: {
        config+:: {
            region:: config.networkRegion,
            zone:: config.networkZone,
            subnet:: config.networkSubnet,
        }
    },
    master+:: {
        config+:: {
            k8s+:: {
                "token": config.token,
            },
            networkName:: "k8s-net",
            instancePasswordRef:: "password",
        }
    },
    node+:: {
        config+:: {
            k8s+:: {
                "token":: config.token,
                "master-ip":: config.masterIP,
            },
            networkName:: "k8s-net",
            instancePasswordRef:: "password",
        }
    },
};

archon + {
    v1+:: archon.v1 + mixin,
}
