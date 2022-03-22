package ggm

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

type modelInfo struct {
	PrimaryKey            string
	PrimaryKeyStructField string
	Fields                []*dbFields
	OtherFields           []*dbFields
	HasOne                []*hasOpts
	HasMany               []*hasOpts
}

type dbFields struct {
	Name    string
	Type    string
	Options []string
	Field   string
}

func getTElem[T any]() reflect.Type {
	t := reflect.TypeOf(new(T))
	if t.Elem().Kind() == reflect.Pointer {
		t = reflect.TypeOf(*new(T))
	}
	return t.Elem()
}

func getRealType(p reflect.Type) reflect.Type {
	if p.Kind() == reflect.Ptr {
		return p.Elem()
	}
	return p
}

func getModelInfo(t reflect.Type) (*modelInfo, error) {
	t = getRealType(t)
	if t.Kind() != reflect.Struct {
		return nil, errors.New("generic type T must be a struct")
	}
	m := &modelInfo{
		Fields:      []*dbFields{},
		OtherFields: []*dbFields{},
		HasOne:      []*hasOpts{},
		HasMany:     []*hasOpts{},
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		dbTag := f.Tag.Get("db")
		tagHasOne := f.Tag.Get("hasOne")
		tagHasMany := f.Tag.Get("hasMany")

		if dbTag == "" && tagHasOne == "" && tagHasMany == "" {
			continue
		}

		// hasMany 时, db tag 在 子 struct 中定义
		if tagHasMany != "" {
			hasManyOpt, err := explodeHasStr(tagHasMany + ",")
			if err != nil {
				return nil, err
			}
			if f.Type.Kind() != reflect.Slice {
				return nil, errors.New("hasMany field must be a slice")
			}
			mm, err := getModelInfo(getRealType(f.Type).Elem())
			if err != nil {
				return nil, err
			}
			for _, v := range mm.Fields {
				hasManyOpt.OtherKeys = append(hasManyOpt.OtherKeys, v.Name)
			}

			hasManyOpt.RefType = f.Type
			hasManyOpt.StructField = f.Name
			m.HasMany = append(m.HasMany, hasManyOpt)
		}

		// 当 为本地字段和 hasOne 字段时, 需要定义db tag
		if dbTag == "" {
			continue
		}

		tokens := strings.Split(dbTag, ",")
		if tokens[0] == "" {
			continue
		}

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
			m.PrimaryKeyStructField = f.Name
		}

		if tagHasOne != "" {
			hasOneOpt, err := explodeHasStr(tagHasOne + "," + dbTag)
			if err != nil {
				return nil, err
			}
			m.HasOne = append(m.HasOne, hasOneOpt)
		}

		fmt.Println(tagHasMany)

	}
	m.HasOne = uniqueHas(m.HasOne)
	m.HasMany = uniqueHas(m.HasMany)
	return m, nil
}

func structInfo[T any]() (*modelInfo, error) {
	t := getTElem[T]()
	m, err := getModelInfo(t)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func uniqueHas(opts []*hasOpts) (uniq []*hasOpts) {
	tmp := make(map[string]*hasOpts)
	for _, v := range opts {
		k := v.Conn + v.DB + v.Table
		if _v, ok := tmp[k]; ok {
			_v.OtherKeys = append(tmp[k].OtherKeys, v.OtherKeys...)
			tmp[k] = _v
		} else {
			tmp[k] = v
		}
	}
	for _, v := range tmp {
		uniq = append(uniq, v)
	}
	return uniq
}
