//go:build !windows

package gengort

import (
	"github.com/ebitengine/purego"
)

type solib struct {
	Handle uintptr
}

func (w solib) Lookup(name string) uintptr {
	addr, err := purego.Dlsym(w.Handle, name)
	if err != nil {
		return 0
	}
	return addr
}

func Dlopen(name string) (Library, error) {
	h, err := purego.Dlopen(name, purego.RTLD_NOW|purego.RTLD_LOCAL)
	if err != nil {
		name = "lib" + name + ".so"
		h, err = purego.Dlopen(name, purego.RTLD_NOW|purego.RTLD_LOCAL)
		if err != nil {
			name = "./" + name
			h, err = purego.Dlopen(name, purego.RTLD_NOW|purego.RTLD_LOCAL)
			if err != nil {
				return nil, err
			}
		}
	}
	return solib{Handle: h}, nil
}
