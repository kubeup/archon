package jsonnet

import (
	"encoding/json"
	"errors"

	"github.com/strickyak/jsonnet_cgo"
)

const entrypoint = "entrypoint.jsonnet"

type VM struct {
	vm      *jsonnet.VM
	profile string
	config  map[string]string
}

func Make(profile string) (*VM, error) {
	vm := jsonnet.Make()
	setProfile(profile)
	vm.ImportCallback(importFunc)
	snippet := `import "archon.libsonnet"`
	_, err := vm.EvaluateSnippet(entrypoint, snippet)
	if err != nil {
		return nil, errors.New("Invalid Profile")
	}
	return &VM{
		vm:      vm,
		profile: profile,
		config:  make(map[string]string),
	}, nil
}

func (vm *VM) New(resource, name string) (string, error) {
	snippet := `
	local archon = import "archon.libsonnet";
	function(resource, name, config)
		archon.v1[resource].new(name, config)
	`
	vm.vm.TlaVar("resource", resource)
	vm.vm.TlaVar("name", name)
	config, err := json.Marshal(vm.config)
	if err != nil {
		return "", err
	}
	vm.vm.TlaCode("config", string(config))
	return vm.vm.EvaluateSnippet(entrypoint, snippet)
}

func (vm *VM) Config(key, value string) {
	vm.config[key] = value
}

func (vm *VM) Destroy() {
	vm.vm.Destroy()
}
