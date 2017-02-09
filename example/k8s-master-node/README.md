One master multiple nodes Kubernetes cluster
=================================

Step 1
------

First, please follow the [installation guide] to install `archon-controller`
locally or into your Kubernetes cluster.


Step 2
------

Create a new namespace for this cluster:

```
kubectl create k8s-master-node
```

Step 3
------

Modify `k8s-user.yaml`. Replace `YOUR_SSH_KEY` with your public key which will be
used for authentication with the server. And create the user resource.

```
kubectl create -f k8s-user.yaml --namespace=k8s-master-node
```

Step 4
------

Create the vpc network and subnet:

```
kubectl create -f k8s-net.yaml --namespace=k8s-master-node
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

And `k8s-node` with content below:

```
{
    "Version": "2012-10-17",
        "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ec2:DescribeInstances",
                "ec2:AttachVolume",
                "ec2:DetachVolume",
                "ecr:GetAuthorizationToken",
                "ec2:DescribeVolumes"
            ],
            "Resource": [
                "*"
            ]
        }]
}
```

Step 6
------

Modify `k8s-master.yaml` and `k8s-node.yaml`. Replace `PUT YOUR CA CERTIFICATE HERE` with the content of
`ca.pem` file you generated with `cfssl` during the installation process.

Step 7
------

Create the master instance group and let the `archon-controller` create the instance for you:

```
kubectl create -f k8s-master.yaml --namespace=k8s-master-node
```

Step 8
------

Use `kubectl get instance -o yaml` to get the `PrivateIP` status of your master instance. And
set the ip address in `k8s-node.yaml` replacing `MASTER_PRIVATE_IP`.

Create the node instance group:

```
kubectl create -f k8s-node.yaml --namespace=k8s-master-node
```

Step 9
-------

Scale your node instance group by editing `replicas` field with `kubectl edit`.


[installation guide]: https://github.com/kubeup/archon/blob/master/docs/installation.md
