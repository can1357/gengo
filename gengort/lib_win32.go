//go:build windows

package gengort

import (
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

func Dlopen(name string) (Library, error) {
	dll, err := windows.LoadDLL(name)
	if err != nil {
		return nil, err
	}
	return windll{dll: dll}, nil
}
