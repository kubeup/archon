local schema = import "archon.alpha.1/schema.libsonnet";

local schemaMixins = {
    user+:: {
        new(name, config={}):: super.new(name),
        newWithPasswordHash(name, passwordHash):: self.new(name) + self.mixin.spec.passwordHash(passwordHash),
        newWithSSHAuthorizedKeys(name, sshAuthorizedKeys):: self.new(name) + self.mixin.spec.sshAuthorizedKeys(sshAuthorizedKeys),
        config:: {
        },
    },
    instanceGroup+:: {
        new(name, config={})::
            local finalConfig = self.config + config;
            local spec = self.mixin.spec;
            local template = self.mixin.spec.template;
            super.new(name) +
            spec.replicas(finalConfig.replicas) +
            spec.selector.matchLabels(finalConfig.instanceLabel) +
            template.metadata.labels(finalConfig.instanceLabel) +
            template.spec.files(std.prune([self.files[x](finalConfig) for x in std.objectFieldsAll(self.files)])),
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
    etcd:: self.instanceGroup + {
        config+:: {
            instanceLabel:: { app: "etcd-cluster" },
        },
    },
    network+:: {
        new(name, config={})::
            local finalConfig = self.config + config;
            local spec = self.mixin.spec;
            super.new(name) +
            spec.region(finalConfig.region) +
            spec.zone(finalConfig.zone) +
            spec.subnet(finalConfig.subnet),
        config:: {
        },
    },
};

schema + {
    v1:: schema.v1 + schemaMixins,
}
