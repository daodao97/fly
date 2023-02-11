package util

import (
	"reflect"
)

// Deref is Indirect for reflect.Types
func Deref(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

var typeChecker = map[string]func(t reflect.Type) bool{
	"*[]struct":             isPtrSliceStruct,
	"*[]*struct":            isPtrSlicePtrStruct,
	"struct":                isStruct,
	"*struct":               isPtrStruct,
	"**struct":              isPtrPtrStruct,
	"map[string]interface":  isMapStrInterface,
	"*map[string]interface": isPtrMapStrInterface,
}

func AllowType(v interface{}, types []string) (ok bool) {
	ty := reflect.TypeOf(v)
	for _, t := range types {
		if ok {
			return ok
		}

		if checker, has := typeChecker[t]; has {
			ok = checker(ty)
		}
	}

	return ok
}

func isMapStrInterface(t reflect.Type) bool {
	return t.Kind() == reflect.Map && t.Key().Kind() == reflect.String && t.Elem().Kind() == reflect.Interface
}

func isPtrMapStrInterface(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Map && t.Elem().Key().Kind() == reflect.String && t.Elem().Elem().Kind() == reflect.Interface
}

func isStruct(t reflect.Type) bool {
	return t.Kind() == reflect.Struct
}

func isPtrStruct(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}

func isPtrPtrStruct(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Ptr && t.Elem().Elem().Kind() == reflect.Struct
}

func isPtrSliceStruct(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Slice && t.Elem().Elem().Kind() == reflect.Struct
}

func isPtrSlicePtrStruct(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Slice && t.Elem().Elem().Kind() == reflect.Ptr && t.Elem().Elem().Elem().Kind() == reflect.Struct
}
