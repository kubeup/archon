#! /bin/bash

src=vendor/k8s.io/kubernetes/staging/src/k8s.io/
dst=vendor/k8s.io/
packages="apimachinery client-go apiserver"

for i in $packages; do
  rm -rf $dst$i
  cp -ar $src$i $dst
done
