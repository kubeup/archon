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
    etcd+:: {
        files+:: f.etcd,
        config+:: {
            etcdName:: "default",
            etcdDataDir:: "/var/lib/etcd/default.etcd",
            etcdListenClientUrls:: "http://localhost:2379",
            etcdAdvertiseClientURLs:: "http://localhost:2379",
            etcdInitialClusterState:: "new",
            etcdListenPeerUrls:: "http://localhost:2380",
            etcdInitialAdvertisePeerURLs:: "http://localhost:2380",
            etcdInitialCluster:: "default=http://localhost:2380",
            etcdInitialClusterToken:: "etcd-cluster",
            etcdOtherPeerClientURLs:: "",
        },
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

