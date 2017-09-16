# Copyright 2016 The Archon Authors.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#     http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

GO    := GO15VENDOREXPERIMENT=1 go
pkgs   = $(shell $(GO) list ./... | grep -v /vendor/)

DOCKER_IMAGE_NAME       ?= archon-controller
DOCKER_IMAGE_TAG        ?= $(subst /,-,$(shell git rev-parse --abbrev-ref HEAD))

all: format test

test:
	@echo ">> running tests"
	@$(GO) test -short $(pkgs)

style:
	@echo ">> checking code style"
	@! gofmt -d $(shell find . -path ./vendor -prune -o -name '*.go' -print) | grep '^'

format:
	@echo ">> formatting code"
	@$(GO) fmt $(pkgs)

vet:
	@echo ">> vetting code"
	@$(GO) vet $(pkgs)

docker:
	@echo ">> building docker image"
	@docker build -t "$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)" .

test-docker:
	@echo ">> test docker image"
	@docker run "$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)" /archon-controller --test-run

update-jsonnet:
	@echo ">> updating jsonnet"
	@mkdir -p assets
	@git clone https://github.com/ksonnet/ksonnet-lib.git
	@cp -r ksonnet-lib/ksonnet.beta.2 assets/
	@cp -r archon.alpha.1 assets/
	@mkdir -p assets/aws-centos
	@cp example/k8s-aws-centos-jsonnet/*.libsonnet assets/aws-centos
	@mkdir -p assets/ucloud-centos
	@cp example/k8s-ucloud-centos-jsonnet/*.libsonnet assets/ucloud-centos
	@go-bindata -pkg jsonnet -prefix assets/ -o pkg/jsonnet/assets.go assets/...
	@rm -rf assets ksonnet-lib

.PHONY: all style format test vet docker update-jsonnet

