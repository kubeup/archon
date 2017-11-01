local archon = import "archon.alpha.1/archon.libsonnet";
local file = archon.v1.instance.mixin.spec.filesType;

local yumRepos = |||
  kubernetes:
    name: Kubernetes
    baseurl: http://yum.kubernetes.io/repos/kubernetes-el7-x86_64
    enabled: true
    gpgcheck: true
    gpgkey: https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
|||;

local packages = |||
  - docker
  - kubelet
  - kubeadm
  - kubectl
  - kubernetes-cni
|||;

local setenforce = |||
  - setenforce
  - "0"
|||;

local enableDocker = |||
  - systemctl
  - enable
  - docker
|||;

local startDocker = |||
  - systemctl
  - start
  - docker
|||;

local enableKubelet = |||
  - systemctl
  - enable
  - kubelet
|||;

local startKubelet = |||
  - systemctl
  - start
  - kubelet
|||;

local masterKubeadm = |||
  - kubeadm
  - init
  - --pod-network-cidr
  - %(k8sPodIPRange)s
  - --token
  - %(k8sToken)s
|||;

local nodeKubeadm = |||
  - kubeadm
  - join
  - --token
  - %(k8sToken)s
  - %(k8sMasterIP)s
|||;

local etcdPackages = |||
  - etcd
|||;


local enableEtcd = |||
  - systemctl
  - enable
  - etcd
|||;

local startEtcd = |||
  - systemctl
  - start
  - etcd
  - --no-block
|||;

local etcdConfig = |||
  ETCD_NAME="%(etcdName)s"
  ETCD_DATA_DIR="%(etcdDataDir)s"
  ETCD_LISTEN_CLIENT_URLS="%(etcdListenClientUrls)s"
  ETCD_ADVERTISE_CLIENT_URLS="%(etcdAdvertiseClientURLs)s"
  ETCD_INITIAL_CLUSTER_STATE="%(etcdInitialClusterState)s"
  ETCD_LISTEN_PEER_URLS="%(etcdListenPeerUrls)s"
  ETCD_INITIAL_ADVERTISE_PEER_URLS="%(etcdInitialAdvertisePeerURLs)s"
  ETCD_INITIAL_CLUSTER="%(etcdInitialCluster)s"
  ETCD_INITIAL_CLUSTER_TOKEN="%(etcdInitialClusterToken)s"
  ETCD_OTHER_PEERS_CLIENT_URLS="%(etcdOtherPeerClientURLs)s"
|||;

local etcdCleanup = |||
  #!/bin/bash

  export ETCDCTL_API=3
  if [[ ! -z "${ETCD_OTHER_PEERS_CLIENT_URLS}" && $ETCD_INITIAL_CLUSTER_STATE == "existing" && ! -d $ETCD_DATA_DIR ]]
  then
    member_id=`etcdctl --endpoints $ETCD_OTHER_PEERS_CLIENT_URLS member list | grep "${ETCD_NAME}" | cut -d , -f 1`
    if [[ $member_id != "" ]]
    then
      etcdctl --endpoints $ETCD_OTHER_PEERS_CLIENT_URLS member remove ${member_id}
      etcdctl --endpoints $ETCD_OTHER_PEERS_CLIENT_URLS member add ${ETCD_NAME} --peer-urls ${ETCD_INITIAL_ADVERTISE_PEER_URLS}
    fi
  fi
|||;

local etcdCleanupDropIn = |||
  [Service]
  ExecStartPre=-/usr/bin/etcd_cleanup
|||;

local reloadSystemd = |||
  - systemctl
  - daemon-reload
|||;

{
  shared:: {
    i10yumRepos(config):: file.new() + file.name("yum-repos") + file.path("/config/yum_repos") + file.content(yumRepos),
    i20packages(config):: file.new() + file.name("packages") + file.path("/config/packages") + file.content(packages),
    i30setenforce(config):: file.new() + file.name("setenforce") + file.path("/config/runcmd/enable-docker") + file.content(setenforce),
    i40enableDocker(config):: file.new() + file.name("enable-docker") + file.path("/config/runcmd/enable-docker") + file.content(enableDocker),
    i50startDocker(config):: file.new() + file.name("start-docker") + file.path("/config/runcmd/start-docker") + file.content(startDocker),
    i60enableKubelet(config):: file.new() + file.name("enable-kubelet") + file.path("/config/runcmd/enable-kubelet") + file.content(enableKubelet),
    i70startKubelet(config):: file.new() + file.name("start-kubelet") + file.path("/config/runcmd/start-kubelet") + file.content(startKubelet),
  },
  master:: self. shared + {
    i80kubeadm(config):: file.new() + file.name("kubeadm") + file.path("/config/runcmd/kubeadm") + file.content(masterKubeadm % config),
  },
  node:: self.shared + {
    i80kubeadm(config):: file.new() + file.name("kubeadm") + file.path("/config/runcmd/kubeadm") + file.content(nodeKubeadm % config),
  },
  etcd:: {
    i10packages(config):: file.new() + file.name("packages") + file.path("/config/packages") + file.content(etcdPackages),
    i20etcdConfig(config):: file.new() + file.name("etcd.conf") + file.path("/etc/etcd/etcd.conf") + file.content(etcdConfig % config),
    i30etcdCleanup(config):: file.new() + file.name("20-cleanup.conf") + file.path("/etc/systemd/system/etcd.service.d/20-cleanup.conf") + file.content(etcdCleanupDropIn),
    i40etcdCleanup(config):: file.new() + file.name("etcd_cleanup") + file.path("/usr/bin/etcd_cleanup") + file.content(etcdCleanup) + file.permissions("0755"),
    i50reloadSystemd(config):: file.new() + file.name("reload-systemd") + file.path("/config/runcmd/reload-systemd") + file.content(reloadSystemd),
    i60enableEtcd(config):: file.new() + file.name("enable-etcd") + file.path("/config/runcmd/enable-etcd") + file.content(enableEtcd),
    i70startEtcd(config):: file.new() + file.name("start-etcd") + file.path("/config/runcmd/start-etcd") + file.content(startEtcd),
  },
}
