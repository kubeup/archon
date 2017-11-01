{
    sshAuthorizedKeys:: "YOUR_SSH_PUBLIC_KEY",
    instancePassword:: "archon",
    networkRegion:: "cn-sh2",
    networkZone:: "cn-sh2-02",
    networkSubnet:: "10.99.0.0/24",
    k8sToken:: "YOUR_KUBEADM_TOKEN",
    k8sMasterIP:: "MASTER_IP",
    // Set configs below for a HA cluster
    k8sMasterHA:: "false",
    k8sMasterState:: "new",
    etcdInitialClusterState:: "new",
    etcd1IP:: "",
    etcd2IP:: "",
    etcd3IP:: "",
    ucloudPublicKey:: "",
    ucloudPrivateKey:: "",
    ucloudProject:: "",
}
