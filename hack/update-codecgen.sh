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

codecgen >/dev/null
if [ "$?" -eq "0" ]; then
  CODECGEN=codecgen
else
  CODECGEN=$(pwd)/codecgen_binary
  go build -o "${CODECGEN}" ./vendor/github.com/ugorji/go/codec/codecgen
fi

generated_files=($(
  find . -not \( \
      \( \
        -wholename './output' \
        -o -wholename './_output' \
        -o -wholename './staging' \
        -o -wholename './release' \
        -o -wholename './target' \
        -o -wholename '*/third_party/*' \
        -o -wholename '*/vendor/*' \
        -o -wholename '*/codecgen-*-1234.generated.go' \
      \) -prune \
    \) -name '*.generated.go' | LC_ALL=C sort -r
))

for i in ${generated_files[@]}; do
  file=${i/\.generated\.go/.go}
  base_file=$(basename "${file}")
  base_generated_file=$(basename "${i}")
  pushd "$(dirname ${file})" > /dev/null
  echo Generating for $file
  rm -f ${base_generated_file} 
  $CODECGEN -d 1234 -o "${base_generated_file}" "${base_file}"
  popd > /dev/null
done

