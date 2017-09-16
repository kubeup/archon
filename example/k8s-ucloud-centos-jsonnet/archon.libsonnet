local archon = import "archon.alpha.1/cloud/ucloud/centos/archon.libsonnet";
local config = import "config.libsonnet";

local mixin =  {
    password+:: {
        new(name, config={})::
            local finalConfig = self.config + config;
            super.new(name, finalConfig.instancePassword),
        config:: {
            instancePassword:: config.instancePassword,
        }
    },
    user+:: {
        new(name, config={})::
            local finalConfig = self.config + config;
            super.new(name) + self.mixin.spec.sshAuthorizedKeys(finalConfig.sshAuthorizedKeys),
        config+:: {
            sshAuthorizedKeys:: config.sshAuthorizedKeys,
        }
    },
    instanceGroup+:: {
        new(name, config={}):: super.new(name, config) + self.mixin.spec.template.spec.users({name: "k8s-user"}),
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
            k8sToken: config.k8sToken,
            networkName:: "k8s-net",
            instancePasswordRef:: "password",
        }
    },
    node+:: {
        config+:: {
            k8sToken:: config.k8sToken,
            k8sMasterIP:: config.k8sMasterIP,
            networkName:: "k8s-net",
            instancePasswordRef:: "password",
        }
    },
};

archon + {
    v1+:: archon.v1 + mixin,
}
