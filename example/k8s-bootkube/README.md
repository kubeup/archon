Self-hosted Kubernetes cluster with bootkube
============================================

In this guide, we'll demonstrate how to bootstrap a Kubernetes cluster with
self-hosted etcd using [bootkube].

Since archon could generate the tls and manifests assets, there's no need to call `bootkube render`.
We can call `bootkube start` directly after bootup.

Step 1
------

First, please follow the [installation guide] to install `archon-controller`
locally or into your Kubernetes cluster.


Step 2
------

Create a new namespace for this cluster:

```
kubectl create k8s-bootkube
```

Step 3
------

Create the default user. The username is `core` and password is `archon`:

```
kubectl create -f k8s-user.yaml --namespace=k8s-bootkube
```

Step 4
------

Create the vpc network and subnet:

```
kubectl create -f k8s-net.yaml --namespace=k8s-bootkube
```

Step 5
------

Create a new instance profile named `k8s-master` with content below:

```
{
    "Version": "2012-10-17",
        "Statement": [{
            "Effect": "Allow",
            "Action": [
                "ec2:DescribeInstances",
                "ec2:AttachVolume",
                "ec2:DetachVolume",
                "ec2:DescribeVolumes",
                "ec2:DescribeSecurityGroups",
                "ec2:CreateSecurityGroup",
                "ec2:DeleteSecurityGroup",
                "ec2:AuthorizeSecurityGroupIngress",
                "ec2:RevokeSecurityGroupIngress",
                "ec2:DescribeSubnets",
                "ec2:CreateTags",
                "ec2:DescribeRouteTables",
                "ec2:CreateRoute",
                "ec2:DeleteRoute",
                "ec2:ModifyInstanceAttribute",
                "ecr:GetAuthorizationToken"
            ],
            "Resource": [
                "*"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "s3:GetObject"
            ],
            "Resource": [
                "*"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "elasticloadbalancing:*"
            ],
            "Resource": [
                "*"
            ]
        }]
}
```

Step 6
------

Create a new s3 bucket in the same region. Modify `k8s-bootkube.yaml`. Replace `YOUR_S3_BUCKET`
with your bucket name.


Step 7
------

Create certificates:

```
kubectl create -f k8s-ca.yaml --namespace=k8s-bootkube
kubectl create -f k8s-apiserver.yaml --namespace=k8s-bootkube
kubectl create -f k8s-kubelet.yaml --namespace=k8s-bootkube
kubectl create -f k8s-serviceaccount.yaml --namespace=k8s-bootkube
```

Step 8
------

Create the instance group and let the `archon-controller` create the instance for you:

```
kubectl create -f k8s-bootkube.yaml --namespace=k8s-bootkube
```

Step 9
------

Modify `k8s-master.yaml` and `k8s-node.yaml`. Replace `INTERNAL_APISERVER_LB` with the
dns name you get from the aws console. Then create the master and slave nodes.
```
kubectl create -f k8s-master.yaml --namespace=k8s-bootkube
kubectl create -f k8s-node.yaml --namespace=k8s-bootkube
```

Step 10
-------

Scale `kube-scheduler`, `kube-controller-manager` and `etcd-cluster` to 3 replicas. Then
drain the `k8s-bootkube` node and remove it from the cluster.

Now you have a self-hosted Kubernetes cluster with HA enabled.

[installation guide]: https://github.com/kubeup/archon#installation
[bootkube]: https://github.com/kubernetes-incubator/bootkube
