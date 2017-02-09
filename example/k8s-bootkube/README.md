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

Modify `k8s-user.yaml`. Replace `YOUR_SSH_KEY` with your public key which will be
used for authentication with the server. And create the user resource.

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

Modify `k8s-bootkube.yaml`. Replace `PUT YOUR CA CERTIFICATE HERE` with the content of
`ca.pem` file you generated with `cfssl` during the installation process.

Step 8
------

Create the instance group and let the `archon-controller` create the instance for you:

```
kubectl create -f k8s-bootkube.yaml --namespace=k8s-bootkube
```

[installation guide]: https://github.com/kubeup/archon/blob/master/docs/installation.md
[bootkube]: https://github.com/kubernetes-incubator/bootkube
