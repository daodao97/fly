package ggm

import (
	"fmt"
	"github.com/jmoiron/sqlx/reflectx"
	dynamicstruct "github.com/ompluscator/dynamic-struct"
	"github.com/pkg/errors"
	"reflect"
)

type reflectValueCache struct {
	Value   reflect.Value
	Pk      string
	PkValue reflect.Value
	Index   int
}

func (m model[T]) hasOneData(list []T) ([]T, error) {
	for _, opt := range m.modelInfo.HasOne {
		_db, exist := DB(opt.Conn)
		if !exist {
			return nil, fmt.Errorf("can not find database conf [%s]", opt.Conn)
		}

		var pkv []any
		var rc []reflectValueCache
		for i, a := range list {
			av := reflect.Indirect(reflect.ValueOf(a))
			f, err := m.getStructFieldNameByTagName(opt.LocalKey)
			if err != nil {
				return nil, err
			}
			pk := reflect.Indirect(av).FieldByName(f)
			pkv = append(pkv, pk.Interface())
			rc = append(rc, reflectValueCache{
				Value:   av,
				Pk:      f,
				PkValue: pk,
				Index:   i,
			})
		}

		if len(pkv) == 0 {
			return nil, errors.New("hasOne pk value is empty")
		}

		fields := append(opt.OtherKeys, opt.ForeignKey+" as "+opt.LocalKey)
		_sql, args := SelectBuilder(Field(fields...), Database(opt.DB), Table(opt.Table), WhereIn(opt.ForeignKey, pkv))

		var _result []T
		err := _db.Select(&_result, _sql, args...)
		if err != nil {
			return nil, err
		}

		for _, av := range rc {
			for _, b := range _result {
				bv := reflect.Indirect(reflect.ValueOf(b))
				fk := bv.FieldByName(av.Pk)
				if av.PkValue.Interface() == fk.Interface() {
					for _, k := range opt.OtherKeys {
						_f, err := m.getStructFieldNameByTagName(k)
						if err != nil {
							return nil, err
						}
						if av.Value.CanSet() {
							av.Value.FieldByName(_f).Set(bv.FieldByName(_f))
						} else {
							_av := reflect.Indirect(reflect.New(av.Value.Type()))
							for _, item := range reflect.VisibleFields(av.Value.Type()) {
								_av.FieldByName(item.Name).Set(av.Value.FieldByName(item.Name))
							}
							_av.FieldByName(_f).Set(bv.FieldByName(_f))
							list[av.Index] = _av.Interface().(T)
						}
					}
				}
			}
		}
	}

	return list, nil
}

func (m model[T]) hasManyData(list []T) ([]T, error) {
	for _, opt := range m.modelInfo.HasMany {
		_db, exist := DB(opt.Conn)
		if !exist {
			return nil, fmt.Errorf("can not find database conf [%s]", opt.Conn)
		}

		var pkv []any
		var rc []reflectValueCache
		for i, a := range list {
			av := reflect.Indirect(reflect.ValueOf(a))
			f, err := m.getStructFieldNameByTagName(opt.LocalKey)
			if err != nil {
				return nil, err
			}
			pk := av.FieldByName(f)
			pkv = append(pkv, pk.Interface())
			rc = append(rc, reflectValueCache{
				Value:   av,
				Pk:      f,
				PkValue: pk,
				Index:   i,
			})
		}

		if len(pkv) == 0 {
			return nil, errors.New("hasMany pk value is empty")
		}

		fields := append(opt.OtherKeys, opt.ForeignKey+" as "+"my_id")
		_sql, args := SelectBuilder(Field(fields...), Database(opt.DB), Table(opt.Table), WhereIn(opt.ForeignKey, pkv))

		rows, err := _db.Queryx(_sql, args...)
		if err != nil {
			return nil, err
		}

		ty := reflectx.Deref(opt.RefType.Elem())
		source := reflect.New(ty).Interface()
		builder := dynamicstruct.ExtendStruct(source).AddField("MyId", 0, `db:"my_id"`)

		var result []any
		for rows.Next() {
			_result := builder.Build().New()
			err := rows.StructScan(_result)
			if err != nil {
				return nil, err
			}
			result = append(result, _result)
		}
		_ = rows.Close()

		for _, av := range rc {
			var tmp []reflect.Value
			for _, b := range result {
				bv := reflect.Indirect(reflect.ValueOf(b))
				fk := bv.FieldByName("MyId")
				if av.PkValue.Interface() == fk.Interface() {
					hasM := reflect.New(ty)
					for i := 0; i < ty.NumField(); i++ {
						f := ty.Field(i)
						reflect.Indirect(hasM).FieldByName(f.Name).Set(bv.FieldByName(f.Name))
					}
					tmp = append(tmp, hasM)
				}
			}
			arr := reflect.MakeSlice(opt.RefType, 0, len(tmp))
			for _, v := range tmp {
				arr = reflect.Append(arr, v)
			}
			if av.Value.CanSet() {
				av.Value.FieldByName(opt.StructField).Set(arr)
			} else {
				_av := reflect.Indirect(reflect.New(av.Value.Type()))
				for _, item := range reflect.VisibleFields(av.Value.Type()) {
					_av.FieldByName(item.Name).Set(av.Value.FieldByName(item.Name))
				}
				_av.FieldByName(opt.StructField).Set(arr)
				list[av.Index] = _av.Interface().(T)
			}
		}
	}

	return list, nil
}

func (m model[T]) getStructFieldNameByTagName(name string) (string, error) {
	for _, v := range m.modelInfo.Fields {
		if v.Name == name {
			return v.Field, nil
		}
	}
	for _, v := range m.modelInfo.OtherFields {
		if v.Name == name {
			return v.Field, nil
		}
	}
	return "", errors.New("not found db tag: " + name)
}
