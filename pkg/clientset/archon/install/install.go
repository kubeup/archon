package install

import (
	"kubeup.com/archon/pkg/clientset/archon"
)

func init() {
	archon.AddToScheme(archon.Scheme)
}
