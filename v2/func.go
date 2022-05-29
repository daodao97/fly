package ggm

import (
	"github.com/jmoiron/sqlx/reflectx"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

func filterEmptyStr(arr []string) []string {
	var tmp []string
	for _, v := range arr {
		if v != "" {
			tmp = append(tmp, v)
		}
	}
	return tmp
}

func InArr(arr []string, el string) bool {
	for _, v := range arr {
		if v == el {
			return true
		}
	}
	return false
}

func IndexOf(arr []string, el string) int {
	for i, v := range arr {
		if v == el {
			return i
		}
	}
	return -1
}

func Remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func RemoveInterface(slice []interface{}, s int) []interface{} {
	return append(slice[:s], slice[s+1:]...)
}

func structFields(el interface{}) ([]string, error) {
	sType := reflectx.Deref(reflect.TypeOf(el))

	if sType.Kind() != reflect.Struct {
		return nil, errors.New("generic type T must be a struct")
	}
	var fields []string
	for i := 0; i < sType.NumField(); i++ {
		f := sType.Field(i)
		tag := f.Tag.Get("db")
		if tag == "" {
			continue
		}
		tokens := strings.Split(tag, ",")
		if tokens[0] == "" {
			continue
		}
		fields = append(fields, tokens[0])
	}
	return fields, nil
}

//snake_case to camelCase
func snakeCaseToCamelCase(inputUnderScoreStr string) (camelCase string) {
	isToUpper := false
	for k, v := range inputUnderScoreStr {
		if k == 0 {
			camelCase = strings.ToUpper(string(inputUnderScoreStr[0]))
		} else {
			if isToUpper {
				camelCase += strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if v == '_' {
					isToUpper = true
				} else {
					camelCase += string(v)
				}
			}
		}
	}
	return
}
