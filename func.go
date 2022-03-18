package ggm

import (
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"reflect"
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

func Remove[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}

func reflectNew[T any]() any {
	t := reflect.TypeOf(new(T))
	if t.Elem().Kind() == reflect.Pointer {
		return reflect.New(t.Elem().Elem()).Interface()
	}

	return reflect.New(t.Elem()).Elem().Interface()
}

func structInfo[T any]() ([]field, error) {
	t := reflect.TypeOf(new(T))
	if t.Elem().Kind() == reflect.Pointer {
		t = reflect.TypeOf(*new(T))
	}
	sType := t.Elem()

	if sType.Kind() != reflect.Struct {
		return nil, errors.New("generic type T must be a struct")
	}
	var fields []field
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
		fields = append(fields, field{
			Name:         tokens[0],
			IsPrimaryKey: InArr(tokens, "pk"),
		})
	}
	return fields, nil
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
