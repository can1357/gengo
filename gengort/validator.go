package gengort

import (
	"fmt"
	"reflect"
	"unsafe"
)

func Validate[Struct any](size, align uintptr, fields ...any) {
	var val Struct
	rtype := reflect.TypeOf(val)

	// Validate size
	if size != unsafe.Sizeof(val) {
		panic(fmt.Sprintf("Mismatching sizof(%s) 0x%d, expected 0x%x", rtype.Name(), unsafe.Sizeof(val), size))
	}

	// Validate alignment
	if align != unsafe.Alignof(val) {
		panic(fmt.Sprintf("Mismatching alignof(%s) 0x%d, expected 0x%x", rtype.Name(), unsafe.Alignof(val), align))
	}

	// Validate fields
	for i := 0; i < len(fields); i += 2 {
		fieldName := fields[i].(string)
		fieldOffset := reflect.ValueOf(fields[i+1]).Int()
		field, ok := rtype.FieldByName(fieldName)
		if !ok {
			panic(fmt.Sprintf("Field %s not found", fieldName))
		}
		if field.Offset != uintptr(fieldOffset) {
			panic(fmt.Sprintf("Mismatching offsetof(%s::%s): 0x%x, expected 0x%x", rtype.Name(), fieldName, field.Offset, fieldOffset))
		}
	}
}
