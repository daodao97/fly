package ggm

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	dynamicstruct "github.com/ompluscator/dynamic-struct"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"reflect"
	"strings"
)

var tableNameNotDefine = errors.New("table name is not define")

type ConnName interface {
	Conn() string
}

type TableName interface {
	Table() string
}

type FakeDeleteKey interface {
	FakeDeleteKey() string
}

func New[T TableName]() *Model[T] {
	m := &Model[T]{}

	t := reflectNew[T]()
	m.table = t.(TableName).Table()
	if m.table == "" {
		m.err = tableNameNotDefine
	}

	if fd, ok := t.(FakeDeleteKey); ok {
		m.fakeDeleteKey = fd.FakeDeleteKey()
	}

	conn := "default"
	if connName, ok := t.(ConnName); ok {
		conn = connName.Conn()
	}
	m.conn(conn)

	info, err := structInfo[T]()
	if err != nil {
		m.err = err
	}
	m.modelInfo = info
	return m
}

func NewClient[T TableName](conn string) *Model[T] {
	m := &Model[T]{}

	t := reflectNew[T]()
	m.table = t.(TableName).Table()
	if m.table == "" {
		m.err = tableNameNotDefine
	}

	if fd, ok := t.(FakeDeleteKey); ok {
		m.fakeDeleteKey = fd.FakeDeleteKey()
	}

	m.conn(conn)

	info, err := structInfo[T]()
	if err != nil {
		m.err = err
	}
	m.modelInfo = info
	return m
}

func NewConn[T TableName](c *Config) *Model[T] {
	m := New[T]()
	client, err := newDb(c)
	m.client = client
	if err != nil {
		m.err = err
	}
	return m
}

type Model[T TableName] struct {
	client        *sqlx.DB
	err           error
	table         string
	fakeDeleteKey string
	modelInfo     *modelInfo
}

func (m *Model[T]) GetDB() *sqlx.DB {
	return m.client
}

func (m *Model[T]) check() error {
	if m.err != nil {
		return m.err
	}
	if m.client == nil {
		return fmt.Errorf("model.Client is nil")
	}
	return nil
}

func (m *Model[T]) conn(name string) *Model[T] {
	client, exist := DB(name)
	m.client = client
	if !exist {
		m.err = fmt.Errorf("connection [%s] not exist", name)
	}
	return m
}

func (m *Model[T]) Select(condition ...Option) ([]T, error) {
	if err := m.check(); err != nil {
		return nil, err
	}

	fields, err := m.tFields()
	if err != nil {
		return nil, err
	}
	condition = append(condition, Table(m.table), Field(fields...))
	if m.fakeDeleteKey != "" {
		condition = append(condition, WhereEq(m.fakeDeleteKey, 0))
	}
	_sql, args := SelectBuilder(condition...)

	var result []T
	err = m.client.Select(&result, _sql, args...)
	if err != nil {
		return nil, err
	}

	result, err = m.hasOneData(result)
	if err != nil {
		return nil, err
	}

	result, err = m.hasManyData(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *Model[T]) SelectOne(condition ...Option) (row T, err error) {
	condition = append(condition, Limit(1), Offset(0))
	result, err := m.Select(condition...)
	if err != nil {
		return row, err
	}
	if Len(result) == 0 {
		return row, fmt.Errorf("result is empty")
	}
	return result[0], nil
}

func (m *Model[T]) Count(condition ...Option) (int, error) {
	if err := m.check(); err != nil {
		return 0, err
	}

	condition = append(condition, Table(m.table), AggregateCount("*"))
	if m.fakeDeleteKey != "" {
		condition = append(condition, WhereNotEq(m.fakeDeleteKey, 1))
	}
	_sql, args := SelectBuilder(condition...)

	result := struct {
		Count int `db:"count"`
	}{}

	err := m.client.Get(&result, _sql, args...)
	if err != nil {
		return 0, err
	}

	return result.Count, nil
}

func (m *Model[T]) Insert(rows ...T) (int64, error) {
	if err := m.check(); err != nil {
		return 0, err
	}

	if Len(rows) == 0 {
		return 0, errors.New("insert data is empty")
	}

	fields, err := m.tFields()
	if err != nil {
		return 0, err
	}

	_sql := InsertNamedBuilder(Table(m.table), Field(fields...))
	result, err := m.client.NamedExec(_sql, rows)
	if err != nil {
		return 0, err
	}

	insertID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get insert id failed, err:%v\n", err)
	}

	return insertID, nil
}

func (m *Model[T]) Update(row T, opt ...Option) (int64, error) {
	if err := m.check(); err != nil {
		return 0, err
	}

	fields, args, err := m.tFieldValue(row)
	if err != nil {
		return 0, err
	}

	pk := m.pk()

	// if pk not exist and opt is empty,
	// then we can not to update some record
	if pk == "" && len(opt) == 0 {
		return 0, errors.New("not found update condition")
	}

	index := lo.IndexOf(fields, pk)
	// if pk exist and opt is empty, and the struct.{pk} is zeroVal
	// then we can not to update some record
	if pk != "" && len(opt) == 0 && index == -1 {
		return 0, errors.New("if update condition is empty, please set the primary key")
	}

	if pk != "" && index > -1 {
		opt = append(opt, WhereEq(pk, args[index]))
		fields = Remove(fields, index)
		args = Remove(args, index)
	}

	opt = append(opt, Table(m.table), Field(fields...), Value(args...))
	_sql, _args := UpdateBuilder(opt...)

	result, err := m.Exec(_sql, _args...)
	if err != nil {
		return 0, err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("get affected failed, err:%v\n", err)
	}

	return affectedRows, nil
}

func (m *Model[T]) Delete(opt ...Option) (bool, error) {
	if len(opt) == 0 {
		return false, errors.New("delete condition is empty")
	}
	opt = append(opt, Table(m.table))
	var result sql.Result
	var err error
	if m.fakeDeleteKey != "" {
		opt = append(opt, Field(m.fakeDeleteKey), Value(1))
		_sql, _args := UpdateBuilder(opt...)
		result, err = m.Exec(_sql, _args...)
	} else {
		_sql, _args := DeleteBuilder(opt...)
		result, err = m.Exec(_sql, _args...)
	}

	if err != nil {
		return false, err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("get affected failed, err:%v\n", err)
	}

	return affectedRows > 0, nil
}

func (m *Model[T]) Exec(sql string, args ...any) (sql.Result, error) {
	return m.client.Exec(sql, args...)
}

func (m *Model[T]) tFields() (fields []string, err error) {
	for _, f := range m.modelInfo.Fields {
		fields = append(fields, f.Name)
	}
	return fields, nil
}

func (m *Model[T]) pk() string {
	return m.modelInfo.PrimaryKey
}

func (m *Model[T]) tFieldValue(t T) (field []string, value []any, err error) {
	rm := reflectx.NewMapper("db")
	fm := rm.FieldMap(reflect.ValueOf(t))

	for k, v := range fm {
		if strings.Contains(k, ".") || v.IsZero() {
			continue
		}
		field = append(field, k)
		value = append(value, v.Interface())
	}

	return field, value, nil
}

type reflectValueCache struct {
	Value   reflect.Value
	Pk      string
	PkValue reflect.Value
}

func (m Model[T]) hasOneData(list []T) ([]T, error) {
	for _, opt := range m.modelInfo.HasOne {
		_db, exist := DB(opt.Conn)
		if !exist {
			return nil, fmt.Errorf("can not find database conf [%s]", opt.Conn)
		}

		var pkv []any
		var rc []reflectValueCache
		for _, a := range list {
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
						av.Value.FieldByName(_f).Set(bv.FieldByName(_f))
					}
				}
			}
		}
	}

	return list, nil
}

func (m Model[T]) hasManyData(list []T) ([]T, error) {
	for _, opt := range m.modelInfo.HasMany {
		_db, exist := DB(opt.Conn)
		if !exist {
			return nil, fmt.Errorf("can not find database conf [%s]", opt.Conn)
		}

		var pkv []any
		var rc []reflectValueCache
		for _, a := range list {
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
			av.Value.FieldByName(opt.StructField).Set(arr)
		}
	}

	return list, nil
}

func (m Model[T]) getStructFieldNameByTagName(name string) (string, error) {
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
