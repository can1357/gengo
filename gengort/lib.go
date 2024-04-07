package gengort

type Library interface {
	Lookup(name string) uintptr
}

//func Dlopen(name string) (Library, error)
