#!/bin/bash

# Copyright 2017 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

KUBE_ROOT=$(dirname "${BASH_SOURCE}")/..
PACKAGE_BASE=${PACKAGE_BASE:-"github.com/kube-node/nodeset"}
CLIENT_PATH=pkg/client
CLIENT_NAME=clientset_v1alpha1

client-gen --help >/dev/null 2>/dev/null
if [ "$?" -eq "2" ]; then
  CLIENT_GEN=client-gen
else
  CLIENT_GEN=$(pwd)/client-gen
  go build -o "${CLIENT_GEN}" ./vendor/k8s.io/kubernetes/cmd/libs/go2idl/client-gen
fi

echo Removing old clientset
rm -rf "pkg/client/$CLIENT_NAME"

echo Generating clientset
${CLIENT_GEN} --input-base "${PACKAGE_BASE}/pkg" --input "nodeset/v1alpha1" --clientset-path "${PACKAGE_BASE}/${CLIENT_PATH}" --clientset-name "$CLIENT_NAME" --fake-clientset=true

# Inject namespace into client coz TPR requires a namespace
find ${CLIENT_PATH} -name '*.go' | xargs sed -i -e '/\tNamespace(c\.ns)\./d' \
  -e 's/\tResource(/\tNamespace(v1alpha1.TPRNamespace).\n\t\tResource(/g'  

# Use json serializer instead of yaml coz we don't have protobuf now
CLIENTFILE=${CLIENT_PATH}/${CLIENT_NAME}/typed/nodeset/v1alpha1/nodeset_client.go
sed -i -e 's/import (/import (\n\t"k8s.io\/apimachinery\/pkg\/runtime"/g' \
  -e 's/config\.APIPath/config.ContentType = runtime.ContentTypeJSON\n\tconfig.APIPath/g' \
  ${CLIENTFILE}





