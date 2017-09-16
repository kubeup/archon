local archon = import "archon.alpha.1/archon.libsonnet";
local centos = import "archon.alpha.1/os/centos/mixins.libsonnet";
local china = import "archon.alpha.1/cloud/china/centos/mixins.libsonnet";
local ucloud = import "archon.alpha.1/cloud/ucloud/mixins.libsonnet";
local mixins = import "archon.alpha.1/cloud/ucloud/centos/mixins.libsonnet";

archon + {
    v1:: archon.v1 + ucloud + centos + china + mixins,
}
