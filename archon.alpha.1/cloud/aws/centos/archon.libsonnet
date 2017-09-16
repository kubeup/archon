local archon = import "archon.alpha.1/archon.libsonnet";
local centos = import "archon.alpha.1/os/centos/mixins.libsonnet";
local aws = import "archon.alpha.1/cloud/aws/mixins.libsonnet";
local mixins = import "archon.alpha.1/cloud/aws/centos/mixins.libsonnet";

archon + {
    v1:: archon.v1 + aws + centos + mixins,
}

