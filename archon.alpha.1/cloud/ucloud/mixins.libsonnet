local k = import "ksonnet.beta.2/k.libsonnet";

{
    instanceGroup+:: {
        new(name)::
          local spec = self.mixin.spec.template.spec;
          local metadata = self.mixin.spec.template.metadata;
          local secretType = self.mixin.spec.template.spec.secretsType;
          local initializers = self.config.initializers;
          super.new(name) +
          spec.instanceType(self.config.instanceType) +
          spec.networkName(self.config.networkName) +
          spec.image(self.config.image) +
          spec.secrets(secretType.name(self.config.instancePasswordRef)) +
          metadata.annotations({initializers: initializers}),
        config+:: {
            instanceType:: "Normal",
            initializers:: "archon.kubeup.com/public-ip",
        },
    },
    network+:: {
        new(name)::
            local vpcID = self.config.vpcID;
            local subnetID = self.config.subnetID;
            local securityGroupID = self.config.securityGroupID;
            local annotations = {
                "ucloud.archon.kubeup.com/vpc-id": vpcID,
                "ucloud.archon.kubeup.com/subnet-id": subnetID,
                "ucloud.archon.kubeup.com/security-group-id": securityGroupID,
            };
            super.new(name) +
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
