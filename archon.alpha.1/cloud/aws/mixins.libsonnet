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
          spec.image(finalConfig.image),
        config+:: {
            instanceType:: "t2.small",
        },
    },
    master+:: {
        new(name, config={})::
          local finalConfig = self.config + config;
          local metadata = self.mixin.spec.template.metadata;
          local instanceProfile = finalConfig.instanceProfile;
          local annotations = {
              "aws.archon.kubeup.com/instance-profile": instanceProfile,
          };
          super.new(name, finalConfig) +
          metadata.annotations(annotations),
    },
    network+:: {
        new(name, config={})::
            local finalConfig = self.config + config;
            local nameServers = finalConfig.nameServers;
            local domainName = finalConfig.region + ".compute.internal";
            local annotations = {
                "aws.archon.kubeup.com/name-servers": nameServers,
                "aws.archon.kubeup.com/domain-name": domainName,
            };
            local labels = {
                "KubernetesCluster": "kubernetes",
            };
            super.new(name, finalConfig) +
            self.mixin.metadata.annotations(annotations) +
            self.mixin.metadata.labels(labels) +
            self.mixin.status.phase("Pending"),
        config+:: {
            nameServers:: "169.254.169.253",
        },
    },
}
