local k = import "ksonnet.beta.2/k.libsonnet";

{
    instanceGroup+:: {
        new(name, config={})::
          local finalConfig = self.config + config;
          local spec = self.mixin.spec.template.spec;
          local metadata = self.mixin.spec.template.metadata;
          local secretType = self.mixin.spec.template.spec.secretsType;
          local initializers = finalConfig.initializers;
          super.new(name, finalConfig) +
          spec.instanceType(finalConfig.instanceType) +
          spec.networkName(finalConfig.networkName) +
          spec.image(finalConfig.image) +
          spec.secrets(secretType.name(finalConfig.instancePasswordRef)) +
          metadata.annotations({initializers: initializers}),
        config+:: {
            instanceType:: "Normal",
            initializers:: "archon.kubeup.com/public-ip",
        },
    },
    network+:: {
        new(name, config={})::
            local finalConfig = self.config + config;
            local vpcID = finalConfig.vpcID;
            local subnetID = finalConfig.subnetID;
            local securityGroupID = finalConfig.securityGroupID;
            local annotations = {
                "ucloud.archon.kubeup.com/vpc-id": vpcID,
                "ucloud.archon.kubeup.com/subnet-id": subnetID,
                "ucloud.archon.kubeup.com/security-group-id": securityGroupID,
            };
            super.new(name, finalConfig) +
            self.mixin.metadata.annotations(annotations) +
            self.mixin.status.phase("Running"),
        config+:: {
            vpcID:: "",
            subnetID:: "",
            securityGroupID:: "",
        },
    },
    password+:: k.core.v1.secret + {
        new(name, password)::
            super.new() +
            self.mixin.metadata.name(name) +
            self.type("kubernetes.io/basic-auth") +
            self.data({password: std.base64(password)}),
    },
}
