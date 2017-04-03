Single machine Kubernetes cluster on Aliyun
===========================================

Step 1
------

First, please follow the [installation guide] to install `archon-controller`
locally or into your Kubernetes cluster.


Step 2
------

Create a new namespace for this cluster:

```
kubectl create k8s-aliyun
```

Step 3
------

Modify `k8s-user.yaml`. Replace `YOUR_SSH_KEY` with your public key which will be
used for authentication with the server. And create the user resource.

```
kubectl create -f k8s-user.yaml --namespace=k8s-aliyun
```

Step 4
------

Create the vpc network and subnet:

```
kubectl create -f k8s-net.yaml --namespace=k8s-aliyun
```

Step 5
------

Create certificates:

```
kubectl create -f k8s-ca.yaml --namespace=k8s-bootkube
kubectl create -f k8s-apiserver.yaml --namespace=k8s-bootkube
kubectl create -f k8s-kubelet.yaml --namespace=k8s-bootkube
kubectl create -f k8s-serviceaccount.yaml --namespace=k8s-bootkube
```

Step 6
------

Replace `DOCKER_REGISTRY_MIRROR` in `k8s-aliyun.yaml` with the one you get from aliyun.

Create the instance group and let the `archon-controller` create the master instance for you:

```
kubectl create -f k8s-aliyun.yaml --namespace=k8s-aliyun
```

Step 7
------

Replace `DOCKER_REGISTRY_MIRROR` in `k8s-aliyun.yaml` with the one you get from aliyun.
And replace `MASTER_HOSTNAME` with the hostname of the master instance.

Create the instance group and let the `archon-controller` create the node instances for you:

```
kubectl create -f k8s-aliyun.yaml --namespace=k8s-node-aliyun
```

Step 7
------

After the new master is created, create a new secret containing aliyun credentials on it:

```bash
kubectl config set-cluster aliyun --server=`k get instances -o yaml|grep publicIP|awk '{ print $2 }'` --insecure-skip-tls-verify
kubectl config set-credentials aliyun --token=kubeup
kubectl config set-context aliyun --cluster=aliyun --user=aliyun
kubectl create secret generic aliyun-creds --context=aliyun --namespace=kube-system --from-literal=accessKey=YourAccessKey --from-literal=accessKeySecret=YourAccessKeySecret
```

[installation guide]: https://github.com/kubeup/archon/blob/master/docs/installation_aliyun.md
