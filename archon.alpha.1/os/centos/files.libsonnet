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
  - %(pod-ip-range)s
  - --token
  - %(token)s
|||;

local nodeKubeadm = |||
  - kubeadm
  - join
  - --token
  - %(token)s
  - %(master-ip)s
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
    i80kubeadm(config):: file.new() + file.name("kubeadm") + file.path("/config/runcmd/kubeadm") + file.content(masterKubeadm % config.k8s),
  },
  node:: self.shared + {
    i80kubeadm(config):: file.new() + file.name("kubeadm") + file.path("/config/runcmd/kubeadm") + file.content(nodeKubeadm % config.k8s),
  },
}
