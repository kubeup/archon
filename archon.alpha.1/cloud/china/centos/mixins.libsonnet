local f = import "archon.alpha.1/cloud/china/centos/files.libsonnet";

{
    master+:: {
        files:: f.master,
    },
    node+:: {
        files:: f.node,
    },
}
