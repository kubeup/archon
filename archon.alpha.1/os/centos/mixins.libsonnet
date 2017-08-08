local archon = import "archon.alpha.1/a.libsonnet";
local f = import "archon.alpha.1/os/centos/files.libsonnet";

{
    instanceGroup+:: {
        new(name)::
            super.new(name) +
            self.mixin.spec.template.spec.os("CentOS"),
    },
    master+:: {
        files+:: f.master,
    },
    node+:: {
        files+:: f.node,
    },
    user+:: {
        new(name)::
            local spec = self.mixin.spec;
            super.new(name) +
            spec.name("centos") +
            spec.sudo("ALL=(ALL) NOPASSWD:ALL") +
            spec.shell("/bin/bash"),
    },
}

