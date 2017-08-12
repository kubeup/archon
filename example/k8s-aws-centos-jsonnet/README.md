Kubernetes cluster with CentOS, jsonnet and kubeadm in AWS
==========================================================

In this guide, we'll demonstrate how to bootstrap a Kubernetes cluster with
CentOS using [kubeadm]. All resources are defined with [jsonnet].

Step 1
------

First, please follow the [installation guide] to install `archon-controller`
locally or into your Kubernetes cluster.


Step 2
------

Create a new namespace for this cluster:

```
kubectl create namespace aws-centos
```

Step 3
------

Fill in your ssh puclic key in `config.libsonnet` which will be
used for authentication with the server. And create the user resource.

```
jsonnet -J PATH_TO_KSONNET_LIB -J PATH_TO_ARCHON k8s-user.jsonnet | kubectl create -f - --namespace=aws-centos
```

Step 4
------

Create the vpc network and subnet:

```
jsonnet -J PATH_TO_KSONNET_LIB -J PATH_TO_ARCHON k8s-net.jsonnet | kubectl create -f - --namespace=aws-centos
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

Generate a token with `python generate_token.py` and replace `TOKEN` in the `config.libsonnet` file.
Then create the master with:

```
jsonnet -J PATH_TO_KSONNET_LIB -J PATH_TO_ARCHON k8s-master.jsonnet | kubectl create -f - --namespace=aws-centos
```

Step 7
------

SSH to the server. Wait for the Kubernetes master to boot up. Then install `flannel` into the cluster:

```
kubeclt apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel-rbac.yml
kubectl apply -f https://github.com/coreos/flannel/raw/master/Documentation/kube-flannel.yml
```

Step 8
------

Replace `MASTER_IP` with the internal ip of the master server in `config.libsonnet`.
Then create the node with:

```
jsonnet -J PATH_TO_KSONNET_LIB -J PATH_TO_ARCHON k8s-node.jsonnet | kubectl create -f - --namespace=aws-centos
```

[installation guide]: https://github.com/kubeup/archon/blob/master/docs/installation_aws.md
[kubeadm]: https://kubernetes.io/docs/admin/kubeadm/
[jsonnet]: http://jsonnet.org
