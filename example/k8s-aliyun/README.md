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

Modify `k8s-aliyun.yaml`. Replace `PUT YOUR CA CERTIFICATE HERE` with the content of
`ca.pem` file you generated with `cfssl` during the installation process.

Step 6
------

Create the instance group and let the `archon-controller` create the instance for you:

```
kubectl create -f k8s-aliyun.yaml --namespace=k8s-aliyun
```

[installation guide]: https://github.com/kubeup/archon/blob/master/docs/installation_aliyun.md
