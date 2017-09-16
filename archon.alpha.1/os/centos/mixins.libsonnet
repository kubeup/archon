local archon = import "archon.alpha.1/a.libsonnet";
local f = import "archon.alpha.1/os/centos/files.libsonnet";

{
    instanceGroup+:: {
        new(name, config={})::
            super.new(name, config) +
            self.mixin.spec.template.spec.os("CentOS"),
    },
    master+:: {
        files+:: f.master,
    },
    node+:: {
        files+:: f.node,
    },
    user+:: {
        new(name, config={})::
            local spec = self.mixin.spec;
            super.new(name, config) +
            spec.name("centos") +
            spec.sudo("ALL=(ALL) NOPASSWD:ALL") +
            spec.shell("/bin/bash"),
    },
}

