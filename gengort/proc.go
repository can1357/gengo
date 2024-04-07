package gengort

import (
	"sync/atomic"
)

type Proc struct {
	library *Library
	name    string
	cache   atomic.Uintptr
}

//go:noinline
func (lp *Proc) addrSlow() uintptr {
	proc := lp.cache.Load()
	if proc == 0 {
		lib, err := lp.library.Get()
		if err != nil {
			panic("failed to load library: " + err.Error())
		}

		proc = lib.Lookup(lp.name)
		if proc == 0 {
			panic("proc not found: " + lp.name)
		}
		lp.cache.Store(proc)
	}
	return proc
}

//go:registerparams
func (lp *Proc) Addr() uintptr {
	proc := lp.cache.Load()
	if proc == 0 {
		proc = lp.addrSlow()
	}
	return proc
}

//go:uintptrescapes
//go:nosplit
func invoke0(proc uintptr) uintptr

//go:uintptrescapes
//go:nosplit
func invoke1(proc uintptr, a uintptr) uintptr

//go:uintptrescapes
//go:nosplit
func invoke2(proc uintptr, a, b uintptr) uintptr

//go:uintptrescapes
//go:nosplit
func invoke3(proc uintptr, a, b, c uintptr) uintptr

//go:uintptrescapes
//go:nosplit
func invoke4(proc uintptr, a, b, c, d uintptr) uintptr

//go:uintptrescapes
//go:nosplit
func invoke5(proc uintptr, a, b, c, d, e uintptr) uintptr

//go:uintptrescapes
//go:nosplit
func invoke6(proc uintptr, a, b, c, d, e, f uintptr) uintptr

//go:uintptrescapes
//go:nosplit
func invoke7(proc uintptr, a, b, c, d, e, f, g uintptr) uintptr

//go:uintptrescapes
//go:nosplit
func invoke8(proc uintptr, a, b, c, d, e, f, g, h uintptr) uintptr

//go:uintptrescapes
//go:nosplit
func invoke9(proc uintptr, a, b, c, d, e, f, g, h, i uintptr) uintptr

//go:uintptrescapes
//go:nosplit
func invoke10(proc uintptr, a, b, c, d, e, f, g, h, i, j uintptr) uintptr

//go:uintptrescapes
//go:nosplit
func invoke11(proc uintptr, a, b, c, d, e, f, g, h, i, j, k uintptr) uintptr

//go:uintptrescapes
//go:nosplit
func invoke12(proc uintptr, a, b, c, d, e, f, g, h, i, j, k, l uintptr) uintptr

//go:uintptrescapes
//go:nosplit
//go:registerparams
func (lp *Proc) Call0() (r1 uintptr) {
	return invoke0(lp.Addr())
}

//go:uintptrescapes
//go:nosplit
//go:registerparams
func (lp *Proc) Call1(a uintptr) (r1 uintptr) {
	return invoke1(lp.Addr(), a)
}

//go:uintptrescapes
//go:nosplit
//go:registerparams
func (lp *Proc) Call2(a, b uintptr) (r1 uintptr) {
	return invoke2(lp.Addr(), a, b)
}

//go:uintptrescapes
//go:nosplit
//go:registerparams
func (lp *Proc) Call3(a, b, c uintptr) (r1 uintptr) {
	return invoke3(lp.Addr(), a, b, c)
}

//go:uintptrescapes
//go:nosplit
//go:registerparams
func (lp *Proc) Call4(a, b, c, d uintptr) (r1 uintptr) {
	return invoke4(lp.Addr(), a, b, c, d)
}

//go:uintptrescapes
//go:nosplit
//go:registerparams
func (lp *Proc) Call5(a, b, c, d, e uintptr) (r1 uintptr) {
	return invoke5(lp.Addr(), a, b, c, d, e)
}

//go:uintptrescapes
//go:nosplit
//go:registerparams
func (lp *Proc) Call6(a, b, c, d, e, f uintptr) (r1 uintptr) {
	return invoke6(lp.Addr(), a, b, c, d, e, f)
}

//go:uintptrescapes
//go:nosplit
//go:registerparams
func (lp *Proc) Call7(a, b, c, d, e, f, g uintptr) (r1 uintptr) {
	return invoke7(lp.Addr(), a, b, c, d, e, f, g)
}

//go:uintptrescapes
//go:nosplit
//go:registerparams
func (lp *Proc) Call8(a, b, c, d, e, f, g, h uintptr) (r1 uintptr) {
	return invoke8(lp.Addr(), a, b, c, d, e, f, g, h)
}

//go:uintptrescapes
//go:nosplit
//go:registerparams
func (lp *Proc) Call9(a, b, c, d, e, f, g, h, i uintptr) (r1 uintptr) {
	return invoke9(lp.Addr(), a, b, c, d, e, f, g, h, i)
}

//go:uintptrescapes
//go:nosplit
//go:registerparams
func (lp *Proc) Call10(a, b, c, d, e, f, g, h, i, j uintptr) (r1 uintptr) {
	return invoke10(lp.Addr(), a, b, c, d, e, f, g, h, i, j)
}

//go:uintptrescapes
//go:nosplit
//go:registerparams
func (lp *Proc) Call11(a, b, c, d, e, f, g, h, i, j, k uintptr) (r1 uintptr) {
	return invoke11(lp.Addr(), a, b, c, d, e, f, g, h, i, j, k)
}

//go:uintptrescapes
//go:nosplit
//go:registerparams
func (lp *Proc) Call12(a, b, c, d, e, f, g, h, i, j, k, l uintptr) (r1 uintptr) {
	return invoke12(lp.Addr(), a, b, c, d, e, f, g, h, i, j, k, l)
}
