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
