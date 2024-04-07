package gengort

import (
	"sync"
	"sync/atomic"
)

var libRegistry = make(map[string]Library)
var libMutex sync.Mutex

func refLibrary(name string) (Library, error) {
	libMutex.Lock()
	defer libMutex.Unlock()
	if lib, ok := libRegistry[name]; ok {
		return lib, nil
	}
	lib, err := Dlopen(name)
	if err != nil {
		return nil, err
	}
	libRegistry[name] = lib
	return lib, nil
}

type Proc struct {
	Library string
	Name    string
	Found   atomic.Uintptr
}

func (lp *Proc) Find() uintptr {
	proc := lp.Found.Load()
	if proc == 0 {
		lib, err := refLibrary(lp.Library)
		if err != nil {
			panic("failed to load library: " + lp.Library)
		}
		proc = lib.Lookup(lp.Name)
		if proc == 0 {
			panic("proc not found: " + lp.Name)
		}
		lp.Found.Store(proc)
	}
	return proc
}
func Import(lib, name string) Proc {
	return Proc{
		Name:    name,
		Library: lib,
	}
}

//go:uintptrescapes
func invoke0(proc uintptr) uintptr

//go:uintptrescapes
func invoke1(proc uintptr, a uintptr) uintptr

//go:uintptrescapes
func invoke2(proc uintptr, a, b uintptr) uintptr

//go:uintptrescapes
func invoke3(proc uintptr, a, b, c uintptr) uintptr

//go:uintptrescapes
func invoke4(proc uintptr, a, b, c, d uintptr) uintptr

//go:uintptrescapes
func invoke5(proc uintptr, a, b, c, d, e uintptr) uintptr

//go:uintptrescapes
func invoke6(proc uintptr, a, b, c, d, e, f uintptr) uintptr

//go:uintptrescapes
func invoke7(proc uintptr, a, b, c, d, e, f, g uintptr) uintptr

//go:uintptrescapes
func invoke8(proc uintptr, a, b, c, d, e, f, g, h uintptr) uintptr

//go:uintptrescapes
func invoke9(proc uintptr, a, b, c, d, e, f, g, h, i uintptr) uintptr

//go:uintptrescapes
func invoke10(proc uintptr, a, b, c, d, e, f, g, h, i, j uintptr) uintptr

//go:uintptrescapes
func invoke11(proc uintptr, a, b, c, d, e, f, g, h, i, j, k uintptr) uintptr

//go:uintptrescapes
func invoke12(proc uintptr, a, b, c, d, e, f, g, h, i, j, k, l uintptr) uintptr

//go:uintptrescapes
func (lp *Proc) Call0() (r1 uintptr) {
	return invoke0(lp.Find())
}

//go:uintptrescapes
func (lp *Proc) Call1(a uintptr) (r1 uintptr) {
	return invoke1(lp.Find(), a)
}

//go:uintptrescapes
func (lp *Proc) Call2(a, b uintptr) (r1 uintptr) {
	return invoke2(lp.Find(), a, b)
}

//go:uintptrescapes
func (lp *Proc) Call3(a, b, c uintptr) (r1 uintptr) {
	return invoke3(lp.Find(), a, b, c)
}

//go:uintptrescapes
func (lp *Proc) Call4(a, b, c, d uintptr) (r1 uintptr) {
	return invoke4(lp.Find(), a, b, c, d)
}

//go:uintptrescapes
func (lp *Proc) Call5(a, b, c, d, e uintptr) (r1 uintptr) {
	return invoke5(lp.Find(), a, b, c, d, e)
}

//go:uintptrescapes
func (lp *Proc) Call6(a, b, c, d, e, f uintptr) (r1 uintptr) {
	return invoke6(lp.Find(), a, b, c, d, e, f)
}

//go:uintptrescapes
func (lp *Proc) Call7(a, b, c, d, e, f, g uintptr) (r1 uintptr) {
	return invoke7(lp.Find(), a, b, c, d, e, f, g)
}

//go:uintptrescapes
func (lp *Proc) Call8(a, b, c, d, e, f, g, h uintptr) (r1 uintptr) {
	return invoke8(lp.Find(), a, b, c, d, e, f, g, h)
}

//go:uintptrescapes
func (lp *Proc) Call9(a, b, c, d, e, f, g, h, i uintptr) (r1 uintptr) {
	return invoke9(lp.Find(), a, b, c, d, e, f, g, h, i)
}

//go:uintptrescapes
func (lp *Proc) Call10(a, b, c, d, e, f, g, h, i, j uintptr) (r1 uintptr) {
	return invoke10(lp.Find(), a, b, c, d, e, f, g, h, i, j)
}

//go:uintptrescapes
func (lp *Proc) Call11(a, b, c, d, e, f, g, h, i, j, k uintptr) (r1 uintptr) {
	return invoke11(lp.Find(), a, b, c, d, e, f, g, h, i, j, k)
}

//go:uintptrescapes
func (lp *Proc) Call12(a, b, c, d, e, f, g, h, i, j, k, l uintptr) (r1 uintptr) {
	return invoke12(lp.Find(), a, b, c, d, e, f, g, h, i, j, k, l)
}
