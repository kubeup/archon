package cluster

import (
	"strings"
)

const InitializerKey = "initializers"

// Just a substitute when native initializer/finalizer for tpr is not ready
func (i *Instance) GetInitializersInAnnotations() []string {
	if i.Annotations == nil {
		return nil
	}
	s, _ := i.Annotations[InitializerKey]
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}

func (i *Instance) SetInitializersInAnnotations(is []string) {
	if i.Annotations == nil {
		if len(is) == 0 {
			return
		}
		i.Annotations = make(map[string]string)
	}
	if len(is) == 0 {
		delete(i.Annotations, InitializerKey)
	} else {
		i.Annotations[InitializerKey] = strings.Join(is, ",")
	}
}
