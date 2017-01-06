Installation
============

The only service you have to run is the `archon-controller`. You can launch it
locally or deploy it into your Kubernetes cluster.

Prerequisites
-------------

You need a running Kubernetes cluster to use Archon. [Google Container Engine]
is a good choice if you need one.

Generating CA certificates
--------------------------

If you don't use ssl certificates in your cluster. You can just skip this step.
For Kubernetes clusters, ssl certificates are needed in `apiserver`, `kubelet`
and other places.

We use [cfssl] to generate the CA certificates with the simple configuration below:

```
{
	"hosts": [
		"ca.example.com"
	],
	"key": {
		"algo": "rsa",
		"size": 4096
	},
	"names": [
		{
			"C": "US",
			"L": "San Francisco",
			"O": "Internet Widgets, LLC",
			"OU": "Certificate Authority",
			"ST": "California"
		}
	]
}
```

Saving the configuration as `ca-csr.json` and create the certificates with
`cfssl gencert -initca ca-csr.json | cfssljson -bare ca -`.

Now you have `ca.pem` and `ca-key.pem` which are needed in following steps.

Launch locally
--------------

You can just launch `archon-controller` locally when you want to make modifications
to the cluster.

First install it with `go get`:

```
go get -u kubeup.com/archon/cmd/archon-controller
```

Then config AWS credentials and run it:

```
export AWS_ACCESS_KEY_ID=YOUR_AWS_KEY_ID
export AWS_SECRET_ACCESS_KEY=YOUR_AWS_SECRET
archon-controller --kubeconfig ~/.kube/config --cloud-provider aws --cluster-signing-cert-file ca.pem --cluster-signing-key-file ca-key.pem
```

Deploy to Kubernetes
--------------------

Before you begin. You should use roles and policies to protect secrets making them
only available to sysadmins.

Create a `secret` containing the CA certificates:

```
kubectl create secret tls archon-ca --cert=ca.pem --key=ca-key.pem --namespace kube-system
```

Create another `secret` containing the AWS credentials:

```
kubectl create secret generic archon-aws --from-literal=AWS_ACCESS_KEY_ID=YOUR_AWS_KEY_ID --from-literal=AWS_SECRET_ACCESS_KEY=YOUR_AWS_SECRET --namespace=kube-system
```

Save the following configuration as `archon-controller.yaml`:

```
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: archon-controller
spec:
  replicas: 1
  template:
    metadata:
      labels:
        name: archon-controller
    spec:

      containers:
      - name: archon-controller
        image: kubeup/archon-controller
        command:
        - "/archon-controller"
        - "--master"
        - kubernetes.default.svc
        - "--cloud-provider"
        - aws
        - "--cluster-signing-cert-file"
        - "/etc/ca/tls.crt"
        - "--cluster-signing-key-file"
        - "/etc/ca/tls.key"
        env:
        - name: AWS_ACCESS_KEY_ID
          valueFrom:
            secretKeyRef:
              name: archon-aws
              key: AWS_ACCESS_KEY_ID
        - name: AWS_SECRET_ACCESS_KEY
          valueFrom:
            secretKeyRef:
              name: archon-aws
              key: AWS_SECRET_ACCESS_KEY
        volumeMounts:
        - mountPath: "/etc/ca"
          name: archon-ca
      volumes:
      - name: archon-ca
        secret:
          secretName: archon-ca
```

And create the deployment with `kubectl create -f archon-controller.yaml --namespace kube-system`

[Google Container Engine]: https://cloud.google.com/container-engine/
[cfssl]: https://github.com/cloudflare/cfssl
