package jsonnet

import (
	"fmt"
	"path"
)

var profile = ""

func importFunc(base, rel string) (string, string, error) {
	fullpath := rel
	dir, _ := path.Split(rel)
	if dir == "" {
		fullpath = path.Join(profile, rel)
	}
	data, err := Asset(fullpath)
	if err != nil {
		return "", "", fmt.Errorf("Import not available %v", rel)
	}
	return string(data), rel, nil
}

func setProfile(p string) {
	profile = p
}
