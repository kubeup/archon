Kubernetes cluster with RedHat and kubeadm
==========================================

In this guide, we'll demonstrate how to bootstrap a Kubernetes cluster with
redhat using [kubeadm].

Step 1
------

First, please follow the [installation guide] to install `archon-controller`
locally or into your Kubernetes cluster.


Step 2
------

Create a new namespace for this cluster:

```
kubectl create k8s-redhat
```

Step 3
------

Modify `k8s-user.yaml`. Replace `YOUR_SSH_KEY` with your public key which will be
used for authentication with the server. And create the user resource.

```
kubectl create -f k8s-user.yaml --namespace=k8s-redhat
```

Step 4
------

Create the vpc network and subnet:

```
kubectl create -f k8s-net.yaml --namespace=k8s-redhat
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

Generate a token with `python generate_token.py` and replace `TOKEN` in the `k8s-master.yaml` file.
Then create the master with:

```
kubectl create -f k8s-master.yaml --namespace=k8s-redhat
```

Step 7
------

SSH to the server. Wait for the the Kubernetes master boots up. Then install `flannel` into the cluster:

```
wget https://github.com/coreos/flannel/raw/master/Documentation/kube-flannel.yml
kubectl create -f kube-flannel.yml -n kube-system
```

Step 8
------

Replace `TOKEN` with the token generated in step 6. Replace `MASTER_IP` with the internal ip of the master server.
Then create the node with:

```
kubectl create -f k8s-node.yaml --namespace=k8s-redhat
```

[installation guide]: https://github.com/kubeup/archon#installation
[kubeadm]: https://kubernetes.io/docs/getting-started-guides/kubeadm/
