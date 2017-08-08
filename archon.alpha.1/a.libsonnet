local archon = import "archon.alpha.1/archon.libsonnet";

local mixins = {
    user+:: {
        newWithPasswordHash(name, passwordHash):: self.new(name) + self.mixin.spec.passwordHash(passwordHash),
        newWithSSHAuthorizedKeys(name, sshAuthorizedKeys):: self.new(name) + self.mixin.spec.sshAuthorizedKeys(sshAuthorizedKeys),
    },
    instanceGroup+:: {
        new(name)::
            local spec = self.mixin.spec;
            local template = self.mixin.spec.template;
            super.new(name) +
            spec.replicas(self.config.replicas) +
            spec.selector.matchLabels(self.config.instanceLabel) +
            template.metadata.labels(self.config.instanceLabel) +
            template.spec.files([self.files[x](self.config) for x in std.objectFieldsAll(self.files)]),
        files:: {
        },
        config:: {
            replicas:: 1,
        },
    },
    master:: self.instanceGroup + {
        config+:: {
            instanceLabel:: { app: "k8s-master" },
        },
    },
    node:: self.instanceGroup + {
        config+:: {
            instanceLabel:: { app: "k8s-node" },
        },
    },
    network+:: {
        new(name)::
            local spec = self.mixin.spec;
            super.new(name) +
            spec.region(self.config.region) +
            spec.zone(self.config.zone) +
            spec.subnet(self.config.subnet),
        config:: {
        },
    },
};

archon + {
    v1:: archon.v1 + mixins,
}
