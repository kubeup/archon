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
GENERATED_FILENAME=zz_generated.deepcopy

deepcopy-gen >/dev/null 2>&1
if [ "$?" -eq "0" ]; then
  DEEPCOPY_GEN=deepcopy-gen
else
  DEEPCOPY_GEN=$(pwd)/deepcopy-gen
  go build -o "${DEEPCOPY_GEN}" ./vendor/k8s.io/kubernetes/cmd/libs/go2idl/deepcopy-gen
fi

generated_dirs=($(
  grep --color=never -l '+k8s:deepcopy-gen=' vendor -r| xargs -n1 dirname | LC_ALL=C sort -u
))

for i in ${generated_dirs[@]}; do
  pushd "$i" > /dev/null
  echo Generating for $i
  rm -f $GENERATED_FILENAME
  $DEEPCOPY_GEN -i . -O $GENERATED_FILENAME
  popd > /dev/null
done


