//go:build windows

package gengort

import (
	"crypto/sha1"
	"encoding/hex"
	"os"

	"golang.org/x/sys/windows"
)

type windll struct {
	dll *windows.DLL
}

func (w windll) Lookup(name string) uintptr {
	proc, err := w.dll.FindProc(name)
	if err != nil {
		return 0
	}
	return proc.Addr()
}

func LoadLibrary(name string) (LoadedLibrary, error) {
	dll, err := windows.LoadDLL(name)
	if err != nil {
		return nil, err
	}
	return windll{dll: dll}, nil
}

func FindLibrary(name string) (LoadedLibrary, error) {
	return LoadLibrary(name)
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
	name := hex.EncodeToString(hash[:]) + ".gengo.dll"
	path := cache + string(os.PathSeparator) + name
	if stat, err := os.Stat(path); err != nil || stat.Size() != int64(len(data)) {
		err = os.WriteFile(path, data, 0644)
		if err != nil {
			return nil, err
		}
	}
	return LoadLibrary(path)
}
