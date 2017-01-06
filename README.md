Archon
======

Archon is an open source tool for cluster creation and daily operations.
It is primarily targeted for Kubernetes clusters but extending to other
distributed clusters should be easy due to its declarative nature.

Archon is designed following the principles of Kubernetes and it works as
an extension of Kubernetes using [ThirdPartyResource]. You can define your
cluster with yaml files then use `kubectl` to create and scale it.

WARNING: Archon is currently in alpha status.

Why Archon
----------

We've already have tools like [kubeadm] and [kops]. Why do we need another tool
for similar job? Here are bunch of reasons:

  1. Tools exist today are mainly imperative. Declarative tools will help
  users share their experience by publish their own cluster definitions.
  Also users could easily make small customization based on existing template
  using declarative tools.
  2. [Self-hosted Kubernetes] sounds like a promising idea and we need cluster
  management tools which could be easily integrated into current Kubernetes
  architecture.
  3. There're many different environment to launch a Kubernetes cluster in,
  an unified experience and easy to extend codebase will help developers to
  collaborate by adding more cloud provider and os support.

How it works
------------

Archon use Kubernetes as its base. Just like the way you use Kubernetes,
define resources in your cluster by creating yaml files using our customized
resource types. Then manage its lifecycle using `kubectl`. Until now, we
support these resources:

  - User
  - Network
  - InstanceGroup
  - Instance

The `archon-controller` which you should launch beforehand in the cluster will
create and manage its status based on your definition.

You might find the idea that bootstrapping a new Kubernetes cluster with an existing
one too complex. But after your cluster got initialized, you can move all the
definitions and controllers into the new cluster and let it manage itself. It will
be very convenient that you can manage both applications running in your cluster
and the cluster itself with just the `kubectl` cli tool.

Supported platforms
-------------------

At the moment, we only support [AWS] and [CoreOS]. More cloud providers and operating
systems support will be added when the core is stable.

Installation
------------

See the [installation instructions].


Example
-------

[Here][simple-example] is a simple example showing how to create a single machine
Kubernetes cluster with Archon. More examples are on their way.

[ThirdPartyResource]: http://kubernetes.io/docs/user-guide/thirdpartyresources/
[kubeadm]: https://github.com/kubernetes/kubeadm
[kops]: https://github.com/kubernetes/kops
[AWS]: https://aws.amazon.com
[CoreOS]: https://coreos.com/os/docs/latest/
[simple-example]: https://github.com/kubeup/archon/tree/master/example/k8s-simple
[Self-hosted Kubernetes]: https://github.com/kubernetes/community/pull/206
[installation instructions]: https://github.com/kubeup/archon/blob/master/docs/installation.md
