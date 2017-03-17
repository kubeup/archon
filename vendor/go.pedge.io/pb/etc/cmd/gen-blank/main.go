package main

import "go.pedge.io/pb/etc/cmd/common"

func main() {
	common.Main(&generateHelper{})
}

type generateHelper struct{}

func (g *generateHelper) ExtraTmplFuncs() map[string]interface{} {
	return nil
}
