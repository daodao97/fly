package ggm

import (
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

type modelInfo struct {
	PrimaryKey  string
	Fields      []*dbFields
	OtherFields []*dbFields
	HasOne      []*hasOpts
	HasMany     []*hasOpts
}

type dbFields struct {
	Name    string
	Type    string
	Options []string
	Field   string
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

func structInfo2[T any]() (*modelInfo, error) {
	t := reflect.TypeOf(new(T))
	if t.Elem().Kind() == reflect.Pointer {
		t = reflect.TypeOf(*new(T))
	}
	sType := t.Elem()

	if sType.Kind() != reflect.Struct {
		return nil, errors.New("generic type T must be a struct")
	}
	m := &modelInfo{
		Fields:      []*dbFields{},
		OtherFields: []*dbFields{},
		HasOne:      []*hasOpts{},
		HasMany:     []*hasOpts{},
	}
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

		tagHasOne := f.Tag.Get("hasOne")
		tagHasMany := f.Tag.Get("hasMany")

		dbf := &dbFields{
			Name:    tokens[0],
			Field:   f.Name,
			Options: tokens[1:],
		}

		if tagHasMany == "" && tagHasOne == "" {
			m.Fields = append(m.Fields, dbf)
		} else {
			m.OtherFields = append(m.OtherFields, dbf)
		}

		if InArr(tokens, "pk") && m.PrimaryKey == "" {
			m.PrimaryKey = tokens[0]
		}

		if tagHasOne != "" {
			hasOneOpt, err := explodeHasStr(tagHasOne + "," + tag)
			if err != nil {
				return nil, err
			}
			m.HasOne = append(m.HasOne, hasOneOpt)
		}

		if tagHasMany != "" {
			hasManyOpt, err := explodeHasStr(tagHasMany)
			if err != nil {
				return nil, err
			}
			m.HasMany = append(m.HasOne, hasManyOpt)
		}
	}
	return m, nil
}
