Three nodes etcd cluster
========================

Step 1
------

First, please follow the [installation guide] to install `archon-controller`
locally or into your Kubernetes cluster.


Step 2
------

Create a new namespace for this cluster:

```
kubectl create etcd-cluster
```

Step 3
------

Modify `etcd-user.yaml`. Replace `YOUR_SSH_KEY` with your public key which will be
used for authentication with the server. And create the user resource.

```
kubectl create -f etcd-user.yaml --namespace=etcd-cluster
```

Step 4
------

Create the vpc network and subnet:

```
kubectl create -f etcd-net.yaml --namespace=etcd-cluster
```

Step 5
------

Generate a discovery token at `https://discovery.etcd.io/new?size=3` and put it in
`etcd-cluster.yaml`.

Create the instance group and let the `archon-controller` create the instances for you:

```
kubectl create -f etcd-cluster.yaml --namespace=etcd-cluster
```

[installation guide]: https://github.com/kubeup/archon/blob/master/docs/installation.md
