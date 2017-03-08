Single machine Kubernetes cluster
=================================

Step 1
------

First, please follow the [installation guide] to install `archon-controller`
locally or into your Kubernetes cluster.


Step 2
------

Create a new namespace for this cluster:

```
kubectl create k8s-simple
```

Step 3
------

Create the default user. The username is `core` and password is `archon`:

```
kubectl create -f k8s-user.yaml --namespace=k8s-simple
```

Step 4
------

Create the vpc network and subnet:

```
kubectl create -f k8s-net.yaml --namespace=k8s-simple
```

Step 5
------

Create the ca:

```
kubectl create -f k8s-ca.yaml --namespace=k8s-simple
```

Step 6
------

Create the instance group and let the `archon-controller` create the instance for you:

```
kubectl create -f k8s-simple.yaml --namespace=k8s-simple
```

[installation guide]: https://github.com/kubeup/archon/blob/master/docs/installation.md
