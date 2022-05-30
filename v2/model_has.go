package ggm

import (
	"fmt"
	"github.com/jmoiron/sqlx/reflectx"
	dynamicstruct "github.com/ompluscator/dynamic-struct"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"reflect"
	"time"
)

type reflectValueCache struct {
	Value   reflect.Value
	Pk      string
	PkValue reflect.Value
	Index   int
}

func (m model) hasOneData(list interface{}) (result interface{}, err error) {
	var kv []interface{}
	defer dbLog("hasOne", time.Now(), &err, &kv)
	for _, opt := range m.modelInfo.HasOne {
		_db, exist := xdb(opt.Conn)
		if !exist {
			return nil, fmt.Errorf("can not find database conf [%s]", opt.Conn)
		}

		var pkv []interface{}
		var rc []reflectValueCache
		s := reflect.Indirect(reflect.ValueOf(list))
		hasKey := snakeCaseToCamelCase(opt.LocalKey)
		for i := 0; i < s.Len(); i++ {
			av := reflect.Indirect(s.Index(i))
			f, err := m.getStructFieldNameByTagName(opt.LocalKey)
			if err != nil {
				return nil, err
			}
			pk := reflect.Indirect(av).FieldByName(f)
			pkv = append(pkv, pk.Interface())
			rc = append(rc, reflectValueCache{
				Value:   av,
				Pk:      hasKey,
				PkValue: pk,
				Index:   i,
			})
		}

		if len(pkv) == 0 {
			return nil, errors.New("hasOne pk value is empty")
		}

		fields := append(opt.OtherKeys, opt.ForeignKey+" as "+opt.LocalKey)
		_sql, args := SelectBuilder(Field(fields...), Database(opt.DB), Table(opt.Table), WhereIn(opt.ForeignKey, pkv))

		kv = append(kv, "sql:", _sql, "args:", args)

		_fields := append(opt.OtherKeys, opt.LocalKey)
		builder := dynamicstruct.NewStruct()
		for _, v := range _fields {
			builder.AddField(snakeCaseToCamelCase(v), 0, fmt.Sprintf(`db:"%s"`, v))
		}
		rows, err := _db.Queryx(_sql, args...)
		if err != nil {
			return nil, errors.Wrap(err, "ggm.HasOne err")
		}
		var _result []interface{}
		for rows.Next() {
			tmp := builder.Build().New()
			err := rows.StructScan(tmp)
			if err != nil {
				return nil, errors.Wrap(err, "ggm.HasOne.HasData")
			}
			_result = append(_result, tmp)
		}
		_ = rows.Close()

		for _, av := range rc {
			for _, b := range _result {
				bv := reflect.Indirect(reflect.ValueOf(b))
				fk := bv.FieldByName(av.Pk)
				if cast.ToString(av.PkValue.Interface()) != cast.ToString(fk.Interface()) {
					continue
				}
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
						s.Index(av.Index).Set(_av)
					}
				}
			}
		}
	}

	return list, nil
}

func (m model) hasManyData(list interface{}) (result interface{}, err error) {
	var kv []interface{}
	defer dbLog("hasMany", time.Now(), &err, &kv)
	for _, opt := range m.modelInfo.HasMany {
		_db, exist := xdb(opt.Conn)
		if !exist {
			return nil, fmt.Errorf("can not find database conf [%s]", opt.Conn)
		}

		var pkv []interface{}
		var rc []reflectValueCache
		s := reflect.Indirect(reflect.ValueOf(list))
		hasKey := snakeCaseToCamelCase(opt.LocalKey)
		for i := 0; i < s.Len(); i++ {
			av := reflect.Indirect(s.Index(i))
			f, err := m.getStructFieldNameByTagName(opt.LocalKey)
			if err != nil {
				return nil, err
			}
			pk := av.FieldByName(f)
			pkv = append(pkv, pk.Interface())
			rc = append(rc, reflectValueCache{
				Value:   av,
				Pk:      hasKey,
				PkValue: pk,
				Index:   i,
			})
		}

		if len(pkv) == 0 {
			return nil, errors.New("hasMany pk value is empty")
		}

		fields := append(opt.OtherKeys, opt.ForeignKey+" as "+"my_id")
		_sql, args := SelectBuilder(Field(fields...), Database(opt.DB), Table(opt.Table), WhereIn(opt.ForeignKey, pkv))

		kv = append(kv, "sql:", _sql, "args:", args)

		rows, err := _db.Queryx(_sql, args...)
		if err != nil {
			return nil, err
		}

		ty := reflectx.Deref(opt.RefType.Elem())
		source := reflect.New(ty).Interface()
		builder := dynamicstruct.ExtendStruct(source).AddField("MyId", 0, `db:"my_id"`)

		var result []interface{}
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
				if cast.ToString(av.PkValue.Interface()) == cast.ToString(fk.Interface()) {
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
				if opt.RefType.Elem().Kind() == reflect.Ptr {
					arr = reflect.Append(arr, v)
				} else {
					arr = reflect.Append(arr, reflect.Indirect(v))
				}
			}
			if av.Value.CanSet() {
				av.Value.FieldByName(opt.StructField).Set(arr)
			} else {
				_av := reflect.Indirect(reflect.New(av.Value.Type()))
				for _, item := range reflect.VisibleFields(av.Value.Type()) {
					_av.FieldByName(item.Name).Set(av.Value.FieldByName(item.Name))
				}
				_av.FieldByName(opt.StructField).Set(arr)
				s.Index(av.Index).Set(_av)
			}
		}
	}

	return list, nil
}

func (m model) getStructFieldNameByTagName(name string) (string, error) {
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
