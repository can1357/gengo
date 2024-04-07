//go:build !windows

package gengort

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

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

func LoadLibrary(name string) (LoadedLibrary, error) {
	h, err := purego.Dlopen(name, purego.RTLD_NOW|purego.RTLD_LOCAL)
	if err == nil {
		return solib{Handle: h}, nil
	}
	if !strings.ContainsAny(name, "/\\") {
		h, elocal := purego.Dlopen("./"+name, purego.RTLD_NOW|purego.RTLD_LOCAL)
		if elocal == nil {
			return solib{Handle: h}, nil
		}
	}
	return nil, err
}

func FindLibrary(name string) (LoadedLibrary, error) {
	lib, err := LoadLibrary(name)
	if err == nil {
		return lib, nil
	}
	org := err
	if !strings.HasSuffix(name, ".so") {
		name += ".so"
		if lib, err = LoadLibrary(name); err == nil {
			return lib, nil
		}
	}
	if !strings.HasPrefix(name, "lib") {
		name = "lib" + name
		if lib, err = LoadLibrary(name); err == nil {
			return lib, nil
		}
	}
	return nil, org
}

func LoadLibraryEmbed(data []byte) (LoadedLibrary, error) {
	cache := getTmpDir()
	hash := sha1.Sum(data)
	name := "." + hex.EncodeToString(hash[:4]) + ".gengo.so"
	path := cache + name
	if stat, err := os.Stat(path); err != nil || stat.Size() != int64(len(data)) {
		os.MkdirAll(cache, 0755)
		err = os.WriteFile(path, data, 0755)
		if err != nil {
			fmt.Println("write file error: ", err)
			return nil, err
		}
	}
	return LoadLibrary(path)
}
