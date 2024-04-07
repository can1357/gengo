package gengort

import (
	"unsafe"
)

//go:inline
func ReadBitcast[T any](p unsafe.Pointer) T {
	return *(*T)(p)
}

//go:inline
func WriteBitcast[T any](p unsafe.Pointer, value T) {
	*(*T)(p) = value
}

//go:inline
//go:uintptrescapes
func MarshallSyscall[T any](data T) uintptr {
	if size := unsafe.Sizeof(data); size > unsafe.Sizeof(uintptr(0)) {
		return uintptr(unsafe.Pointer(&data))
	} else if size == unsafe.Sizeof(uintptr(0)) {
		return *(*uintptr)(unsafe.Pointer(&data))
	} else {
		var buf [unsafe.Sizeof(uintptr(0))]byte
		bufp := unsafe.Pointer(&buf[0])
		*(*T)(bufp) = data
		return *(*uintptr)(bufp)
	}
}

//go:inline
//go:uintptrescapes
func UnmarshallSyscall[T any](ptr uintptr) (res T) {
	if size := unsafe.Sizeof(res); size > unsafe.Sizeof(uintptr(0)) {
		return *(*T)(unsafe.Pointer(ptr))
	} else if size == unsafe.Sizeof(uintptr(0)) {
		return *(*T)(unsafe.Pointer(&ptr))
	} else {
		var buf [unsafe.Sizeof(uintptr(0))]byte
		bufp := unsafe.Pointer(&buf[0])
		*(*uintptr)(bufp) = ptr
		return *(*T)(bufp)
	}
}
