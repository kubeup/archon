Baremetal Kubernetes cluster with PXE and Matchbox
==================================================

Step 1
------

First, please follow the [installation guide] to install `archon-controller`
locally or into your Kubernetes cluster.

Config these environment variables to use your matchbox server:

- `MATCHBOX_HTTP_ENDPOINT`
- `MATCHBOX_RPC_ENDPOINT`
- `MATCHBOX_CA_FILE`
- `MATCHBOX_CERT_FILE`
- `MATCHBOX_KEY_FILE`

Step 2
------

Create a new namespace for this cluster:

```
kubectl create k8s-matchbox
```

Step 3
------

Create the default user. The username is `core` and password is `archon`:

```
kubectl create -f k8s-user.yaml --namespace=k8s-matchbox
```

Step 4
------

Create the vpc network and subnet:

```
kubectl create -f k8s-net.yaml --namespace=k8s-matchbox
```

Step 5
------

Create the ca and certificates:

```
kubectl create -f k8s-ca.yaml --namespace=k8s-matchbox
kubectl create -f k8s-apiserver.yaml --namespace=k8s-matchbox
kubectl create -f k8s-kubelet.yaml --namespace=k8s-matchbox
kubectl create -f k8s-serviceaccount.yaml --namespace=k8s-matchbox
```

Step 6
------

Edit `k8s-servers.yaml`. Replace server mac and ip address with your server settings.
Create `ReservedInstance`:

```
kubectl create -f k8s-servers.yaml --namespace=k8s-matchbox
```

Step 7
------

Create `installer`, `master`, `node` instances. And let archon populate matchbox configurations
with instance definitions:

```
kubectl create -f k8s-installer.yaml --namespace=k8s-matchbox
kubectl create -f k8s-master.yaml --namespace=k8s-matchbox
kubectl create -f k8s-node.yaml --namespace=k8s-matchbox
```

Step 8
------

Turn on your servers. They will boot with PXE. The master server will install CoreOS to disk
and the node servers will boot in memory as immutable computing nodes.

[installation guide]: https://github.com/kubeup/archon/blob/master/docs/installation.md
