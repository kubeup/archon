local archon = import "archon.alpha.1/cloud/ucloud/centos/archon.libsonnet";
local config = import "config.libsonnet";
local file = archon.v1.instance.mixin.spec.filesType;

local mixin =  {
    password+:: {
        new(name, config={})::
            local finalConfig = self.config + config;
            if finalConfig.instancePassword == "" then
                error "config instancePassword should not be empty"
            else
                super.new(name, finalConfig.instancePassword),
        config:: {
            instancePassword:: config.instancePassword,
        }
    },
    user+:: {
        new(name, config={})::
            local finalConfig = self.config + config;
            if finalConfig.sshAuthorizedKeys == "" then
                error "config sshAuthorizedKeys should not be empty"
            else
                super.new(name) + self.mixin.spec.sshAuthorizedKeys(finalConfig.sshAuthorizedKeys),
        config+:: {
            sshAuthorizedKeys:: config.sshAuthorizedKeys,
        }
    },
    instanceGroup+:: {
        new(name, config={}):: super.new(name, config) + self.mixin.spec.template.spec.users({name: "k8s-user"}),
    },
    network+:: {
        config+:: {
            region:: config.networkRegion,
            zone:: config.networkZone,
            subnet:: config.networkSubnet,
        }
    },
    master+:: {
        new(name, config={})::
            local finalConfig = self.config + config + {instanceLabel+:: {group: name}};
            local pki = [
                {name: "k8s-ca"},
                {name: "serviceaccount"},
            ];
            super.new(name, finalConfig) + if finalConfig.k8sMasterHA == "true" then archon.v1.instanceGroup.mixin.spec.template.spec.secrets(pki) else {},
        config+:: {
            k8sToken: config.k8sToken,
            k8sMasterIP:: config.k8sMasterIP,
            k8sMasterHA:: config.k8sMasterHA,
            k8sMasterState:: config.k8sMasterState,
            networkName:: "k8s-net",
            instancePasswordRef:: "password",
            etcd1IP:: config.etcd1IP,
            etcd2IP:: config.etcd2IP,
            etcd3IP:: config.etcd3IP,
            ucloudPublicKey:: config.ucloudPublicKey,
            ucloudPrivateKey:: config.ucloudPrivateKey,
            ucloudZone:: config.networkZone,
            ucloudRegion:: config.networkRegion,
            ucloudProject:: config.ucloudProject,
        },
        files+:: {
            i40caCrt(config)::
                local caCrt = |||
                  {{ index .Secrets "ca" "tls-cert" | printf "%s" }}
                |||;
                if config.k8sMasterHA == "true" then
                    file.new() + file.name("ca.crt") + file.path("/etc/kubernetes/pki/ca.crt") + file.template(caCrt)
                else null,
            i40caKey(config)::
                local caKey = |||
                  {{ index .Secrets "ca" "tls-key" | printf "%s" }}
                |||;
                if config.k8sMasterHA == "true" then
                    file.new() + file.name("ca.key") + file.path("/etc/kubernetes/pki/ca.key") + file.template(caKey) + file.permissions("0600")
                else null,
            i40saPub(config)::
                local saPub = |||
                  {{ index .Secrets "serviceaccount" "tls-cert" | printf "%s" }}
                |||;
                if config.k8sMasterHA == "true" then
                    file.new() + file.name("sa.pub") + file.path("/etc/kubernetes/pki/sa.pub") + file.template(saPub)
                else null,
            i40saKey(config)::
                local saKey = |||
                  {{ index .Secrets "serviceaccount" "tls-key" | printf "%s" }}
                |||;
                if config.k8sMasterHA == "true" then
                    file.new() + file.name("sa.key") + file.path("/etc/kubernetes/pki/sa.key") + file.template(saKey) + file.permissions("0600")
                else null,
            i40localServerDropIn(config)::
                local localServer = |||
                  [Service]
                  Environment="KUBELET_KUBECONFIG_ARGS=--kubeconfig=/etc/kubernetes/kubelet.conf --api-servers=https://127.0.0.1:6443"
                |||;
                if config.k8sMasterHA == "true" then
                    file.new() + file.name("20-local-server.conf") + file.path("/etc/systemd/system/kubelet.service.d/20-local-server.conf") + file.content(localServer)
                else null,
            i40apiserverService(config)::
                local apiserverService = |||
                  kind: Service
                  apiVersion: v1
                  metadata:
                    name: kube-apiserver
                    annotations:
                      ucloud.archon.kubeup.com/network-type: internal
                      ucloud.archon.kubeup.com/listen-type: PacketsTransmit
                  spec:
                    externalTrafficPolicy: Local
                    selector:
                      component: kube-apiserver
                      tier: control-plane
                    ports:
                    - protocol: TCP
                      port: 6443
                      targetPort: 6443
                    type: LoadBalancer
                |||;
                if config.k8sMasterHA == "true" && config.k8sMasterState == "new" then
                    file.new() + file.name("kube-apiserver-svc.yaml") + file.path("/tmp/kube-apiserver-svc.yaml") + file.content(apiserverService)
                else
                    null,
            i40ucloudController(config)::
                local ucloudController = |||
                  apiVersion: v1
                  kind: Pod
                  metadata:
                    name: ucloud-controller
                    namespace: kube-system
                  spec:
                    hostNetwork: true
                    containers:
                    - name: ucloud-controller
                      image: registry.aliyuncs.com/kubeup/kube-ucloud:latest
                      command:
                      - /ucloud-controller
                      - --server=https://127.0.0.1:6443
                      - --leader-elect=true
                      - --cluster-cidr=10.244.0.0/16
                      - --kubeconfig=/etc/kubernetes/admin.conf
                      - --v=2
                      env:
                      - name: UCLOUD_PUBLIC_KEY
                        value: %(ucloudPublicKey)s
                      - name: UCLOUD_PRIVATE_KEY
                        value: %(ucloudPrivateKey)s
                      - name: UCLOUD_ZONE
                        value: %(ucloudZone)s
                      - name: UCLOUD_REGION
                        value: %(ucloudRegion)s
                      - name: UCLOUD_PROJECT
                        value: %(ucloudProject)s
                      volumeMounts:
                      - mountPath: /etc/kubernetes
                        name: k8s
                        readOnly: true
                    volumes:
                    - hostPath:
                        path: /etc/kubernetes
                      name: k8s
                |||;
                if config.k8sMasterHA == "true" then
                    file.new() + file.name("ucloud-controller.yaml") + file.path("/etc/kubernetes/manifests/ucloud-controller.yaml") + file.content(ucloudController % config) + file.permissions("0600")
                else null,
            i30kubeadmConf(config)::
                local kubeadmConf = |||
                  apiVersion: kubeadm.k8s.io/v1alpha1
                  kind: MasterConfiguration
                  networking:
                    podSubnet: 10.244.0.0/16
                  kubernetesVersion: v1.7.0
                  token: %(k8sToken)s
                  apiServerCertSANs:
                  - 127.0.0.1
                  - %(k8sMasterIP)s
                  controllerManagerExtraArgs:
                    master: https://127.0.0.1:6443
                  schedulerExtraArgs:
                    master: https://127.0.0.1:6443
                |||;
                local etcdConf = |||
                  etcd:
                    endpoints:
                    - http://%(etcd1IP)s:2379
                    - http://%(etcd2IP)s:2379
                    - http://%(etcd3IP)s:2379
                |||;
                local advertiseAddress = |||
                  api:
                    advertiseAddress: %(k8sMasterIP)s
                |||;
                local config2 = if config.etcd1IP != "" && config.etcd2IP != "" && config.etcd3IP != "" then kubeadmConf + etcdConf else kubeadmConf;
                local config3 = if config.k8sMasterState != "new" then config2 + advertiseAddress else config2;
                if config.k8sMasterHA == "true" then
                    file.new() + file.name("kubeadm.conf") + file.path("/tmp/kubeadm.conf") + file.content(config3 % config)
                else null,
            i90finalize(config)::
                local runFinalize = |||
                  - sh
                  - -c
                  - "kubectl --kubeconfig=/etc/kubernetes/admin.conf label node `hostname` node-role.kubernetes.io/master- node-role.archon.kubeup.com/master=''"
                |||;
                if config.k8sMasterHA == "true" then
                    file.new() + file.name("finalize") + file.path("/config/runcmd/finalize") + file.content(runFinalize)
                else null,
            i90finalize2(config)::
                local runFinalize = |||
                  - sh
                  - -c
                  - "kubectl --kubeconfig=/etc/kubernetes/admin.conf apply -f /tmp/kube-apiserver-svc.yaml -n kube-system"
                |||;
                if config.k8sMasterHA == "true" && config.k8sMasterState == "new" then
                    file.new() + file.name("finalize2") + file.path("/config/runcmd/finalize2") + file.content(runFinalize)
                else null,
        },
    },
    node+:: {
        config+:: {
            k8sToken:: config.k8sToken,
            k8sMasterIP:: config.k8sMasterIP,
            networkName:: "k8s-net",
            instancePasswordRef:: "password",
        }
    },
    etcd+:: {
        config+:: {
            networkName:: "k8s-net",
            instancePasswordRef:: "password",
            etcdInitialClusterState:: config.etcdInitialClusterState,
            nodeIP:: "",
            etcd1IP:: config.etcd1IP,
            etcd2IP:: config.etcd2IP,
            etcd3IP:: config.etcd3IP,
        },
        files+:: {
            i30etcdSetIP(config)::
                local etcdSetIPDropIn = |||
                  [Service]
                  ExecStartPre=-/usr/sbin/ip addr add %(nodeIP)s/16 dev eth0
                  ExecStartPre=-/usr/sbin/arping -U %(nodeIP)s -c 5
                  PermissionsStartOnly=true
                |||;
                if config.nodeIP == "" then
                    error "config nodeIP should not be empty"
                else
                    file.new() + file.name("10-setip.conf") + file.path("/etc/systemd/system/etcd.service.d/10-setip.conf") + file.content(etcdSetIPDropIn % config),
        },
        generateConfig(name, ip, config)::
            local nodeConfig = {
                name:: name,
                ip:: ip,
            };
            {
                instanceLabel:: {
                    app: "etcd-cluster",
                    etcd: name,
                },
                nodeIP:: nodeConfig.ip,
                etcdName:: name,
                etcdDataDir:: "/var/lib/etcd/%(name)s.etcd" % nodeConfig,
                etcdListenClientUrls:: "http://localhost:2379,http://%(ip)s:2379" % nodeConfig,
                etcdAdvertiseClientURLs:: "http://%(ip)s:2379" % nodeConfig,
                etcdInitialClusterState:: config.etcdInitialClusterState,
                etcdListenPeerUrls:: "http://%(ip)s:2380" % nodeConfig,
                etcdInitialAdvertisePeerURLs:: "http://%(ip)s:2380" % nodeConfig,
                etcdInitialCluster:: "etcd1=http://%(etcd1IP)s:2380,etcd2=http://%(etcd2IP)s:2380,etcd3=http://%(etcd3IP)s:2380" % config,
                etcdInitialClusterToken:: config.etcdInitialClusterToken,
                etcdOtherPeerClientURLs:: std.join(",", std.map(function(x) x + ":2379", std.filter(function(x) x != ip, [config.etcd1IP, config.etcd2IP, config.etcd3IP]))),
            },
    },
    etcd1:: self.etcd + {
        new(name, config={})::
            local finalConfig1 = self.config + config;
            local finalConfig = finalConfig1 + self.generateConfig("etcd1", finalConfig1.etcd1IP, finalConfig1);
            if finalConfig.etcd1IP == "" || finalConfig.etcd2IP == "" || finalConfig.etcd3IP == "" then
                error "config etcdIP should not be empty"
            else
                super.new(name, finalConfig),
    },
    etcd2:: self.etcd + {
        new(name, config={})::
            local finalConfig1 = self.config + config;
            local finalConfig = finalConfig1 + self.generateConfig("etcd2", finalConfig1.etcd2IP, finalConfig1);
            if finalConfig.etcd1IP == "" || finalConfig.etcd2IP == "" || finalConfig.etcd3IP == "" then
                error "config etcdIP should not be empty"
            else
                super.new(name, finalConfig),
    },
    etcd3:: self.etcd + {
        new(name, config={})::
            local finalConfig1 = self.config + config;
            local finalConfig = finalConfig1 + self.generateConfig("etcd3", finalConfig1.etcd3IP, finalConfig1);
            if finalConfig.etcd1IP == "" || finalConfig.etcd2IP == "" || finalConfig.etcd3IP == "" then
                error "config etcdIP should not be empty"
            else
                super.new(name, finalConfig),
    },
};

archon + {
    v1+:: archon.v1 + mixin,
}
