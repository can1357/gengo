//go:build !windows

package gengort

import (
	"crypto/sha1"
	"encoding/hex"
	"os"
	"path/filepath"
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
	if err != nil {
		rel, e := filepath.Rel(".", name)
		if e == nil {
			if h, e2 := purego.Dlopen(rel, purego.RTLD_NOW|purego.RTLD_LOCAL); e2 == nil {
				return solib{Handle: h}, nil
			}
		}
		return nil, err
	}
	return solib{Handle: h}, nil
}

func FindLibrary(name string) (LoadedLibrary, error) {
	lib, err := LoadLibrary(name)
	if err == nil {
		return lib, nil
	}
	org := err
	if !strings.HasSuffix(name, ".so") {
		name += ".so"
		lib, err = LoadLibrary(name)
		if err == nil {
			return lib, nil
		}
	}
	if strings.ContainsAny(name, "/\\") {
		return nil, err
	}
	if !strings.HasPrefix(name, "lib") {
		lib, err = LoadLibrary("lib" + name)
		if err == nil {
			return lib, nil
		}
	}
	return nil, org
}

func LoadLibraryEmbed(data []byte) (LoadedLibrary, error) {
	cache, err := os.UserCacheDir()
	if err != nil {
		cache, err = os.UserHomeDir()
		if err != nil {
			cache = os.TempDir()
		}
	}
	hash := sha1.Sum(data)
	filename := hex.EncodeToString(hash[:]) + ".gengo.so"
	path := cache + string(os.PathSeparator) + filename
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return nil, err
	}
	return LoadLibrary(path)
}
