package ggm

import (
	"fmt"
	"github.com/jmoiron/sqlx/reflectx"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"reflect"
	"strconv"
	"strings"
)

func Len[T any](a []T) int {
	count := 0
	for range a {
		count++
	}
	return count
}

func InArr[T comparable](arr []T, el T) bool {
	return lo.IndexOf(arr, el) > 0
}

func Remove[T any](arr []T, s int) []T {
	return append(arr[:s], arr[s+1:]...)
}

func Join[T int | string](arr []T) string {
	var tmp []string
	for _, v := range arr {
		tmp = append(tmp, fmt.Sprintf("%v", v))
	}
	return strings.Join(tmp, ",")
}

func Split[T int | string](str string) []T {
	parts := strings.Split(str, ",")
	var tmp []T
	rt := reflectx.Deref(reflect.TypeOf(new(T)))

	for _, v := range parts {
		_v := new(T)
		el := reflect.ValueOf(_v).Elem()

		switch rt.Kind() {
		case reflect.Int:
			t, _ := strconv.Atoi(v)
			el.SetInt(int64(t))
		case reflect.String:
			el.SetString(v)
		}
		tmp = append(tmp, *_v)
	}

	return tmp
}

func filterEmptyStr(arr []string) []string {
	return lo.Filter[string](arr, func(v string, _ int) bool {
		return v != ""
	})
}

func reflectNew[T any]() any {
	t := reflect.TypeOf(new(T))
	if t.Elem().Kind() == reflect.Pointer {
		return reflect.New(t.Elem().Elem()).Interface()
	}

	return reflect.New(t.Elem()).Elem().Interface()
}

func structFields[T any]() ([]string, error) {
	t := reflect.TypeOf(new(T))
	if t.Elem().Kind() == reflect.Pointer {
		t = reflect.TypeOf(*new(T))
	}
	sType := t.Elem()

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

func merge(a, b any) any {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)
	for i := 0; i < bv.NumField(); i++ {
		af := av.Field(i)
		bf := bv.Field(i)
		if !bf.IsZero() {
			af.Set(reflect.ValueOf(bf.Interface()))
		}
	}

	return a
}
