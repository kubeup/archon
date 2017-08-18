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
          spec.image(self.config.image),
        config+:: {
            instanceType:: "t2.small",
        },
    },
    master+:: {
        new(name)::
          local metadata = self.mixin.spec.template.metadata;
          local instanceProfile = self.config.instanceProfile;
          local annotations = {
              "aws.archon.kubeup.com/instance-profile": instanceProfile,
          };
          super.new(name) +
          metadata.annotations(annotations),
    },
    network+:: {
        new(name)::
            local nameServers = self.config.nameServers;
            local domainName = self.config.region + ".compute.internal";
            local annotations = {
                "aws.archon.kubeup.com/name-servers": nameServers,
                "aws.archon.kubeup.com/domain-name": domainName,
            };
            local labels = {
                "KubernetesCluster": "kubernetes",
            };
            super.new(name) +
            self.mixin.metadata.annotations(annotations) +
            self.mixin.metadata.labels(labels) +
            self.mixin.status.phase("Pending"),
        config+:: {
            nameServers:: "169.254.169.253",
        },
    },
}
