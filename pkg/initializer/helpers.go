package initializer

import (
	"k8s.io/kubernetes/pkg/util/sets"
)

func AddInitializer(dst Object, additions ...string) {
	if len(additions) == 0 {
		return
	}

	s := dst.GetInitializers()
	dst.SetInitializers(addStrings(s, additions...))
}

func RemoveInitializer(dst Object, removals ...string) {
	if len(removals) == 0 {
		return
	}

	s := dst.GetInitializers()
	dst.SetInitializers(removeStrings(s, removals...))
}

func HasInitializer(obj Object, name string) bool {
	s := obj.GetInitializers()
	for _, i := range s {
		if i == name {
			return true
		}
	}

	return false
}

func AddFinalizer(dst Object, additions ...string) {
	if len(additions) == 0 {
		return
	}

	s := dst.GetFinalizers()
	dst.SetFinalizers(addStrings(s, additions...))
}

func RemoveFinalizer(dst Object, removals ...string) {
	if len(removals) == 0 {
		return
	}

	s := dst.GetFinalizers()
	dst.SetFinalizers(removeStrings(s, removals...))
}

func HasFinalizer(obj Object, name string) bool {
	s := obj.GetFinalizers()
	for _, i := range s {
		if i == name {
			return true
		}
	}

	return false
}

func addStrings(dst []string, additions ...string) []string {
	s := sets.NewString(dst...)
	for _, a := range additions {
		if !s.Has(a) {
			dst = append(dst, a)
			s.Insert(a)
		}
	}
	return dst
}

func removeStrings(src []string, removals ...string) (dst []string) {
	s := sets.NewString(removals...)
	for _, i := range src {
		if !s.Has(i) {
			dst = append(dst, i)
		}
	}
	return dst
}
