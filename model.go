package ggm

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"reflect"
	"strings"
	"time"
)

var ErrTableNameNotDefine = errors.New("table name is not define")
var ErrNilRow = errors.New("select row is nil")
var ErrPrimaryKeyNotDefined = errors.New("pk is not defined")

type ConnName interface {
	Conn() string
}

type TableName interface {
	Table() string
}

type FakeDeleteKey interface {
	FakeDeleteKey() string
}

func New[T TableName]() *model[T] {
	m := &model[T]{}

	t := reflectNew[T]()
	m.table = t.(TableName).Table()
	if m.table == "" {
		m.err = ErrTableNameNotDefine
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

func NewConn[T TableName](c *Config) *model[T] {
	m := New[T]()
	client, err := newDb(c)
	m.client = client
	if err != nil {
		m.err = err
	}
	return m
}

type model[T TableName] struct {
	client        *sqlx.DB
	err           error
	table         string
	fakeDeleteKey string
	modelInfo     *modelInfo
}

func (m *model[T]) DB() *sqlx.DB {
	return m.client
}

func (m *model[T]) check() error {
	if m.err != nil {
		return m.err
	}
	if m.client == nil {
		return fmt.Errorf("model.Client is nil")
	}
	return nil
}

func (m *model[T]) conn(name string) *model[T] {
	client, exist := xdb(name)
	m.client = client
	if !exist {
		m.err = fmt.Errorf("connection [%s] not exist", name)
	}
	return m
}

func (m *model[T]) Select(condition ...Option) (result []T, err error) {
	var kv []any
	defer dbLog(time.Now(), &err, &kv)
	if err := m.check(); err != nil {
		return nil, errors.Wrap(err, "ggm.Select.check")
	}

	fields, err := m.tFields()
	if err != nil {
		return nil, errors.Wrap(err, "ggm.Select.structFields")
	}
	condition = append(condition, Table(m.table), Field(fields...))
	if m.fakeDeleteKey != "" {
		condition = append(condition, WhereEq(m.fakeDeleteKey, 0))
	}
	_sql, args := SelectBuilder(condition...)
	kv = append(kv, "sql:", _sql, "args:", args)

	err = m.client.Select(&result, _sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "ggm.Select.sqlx_select")
	}

	if lenT(result) == 0 {
		return result, nil
	}

	result, err = m.hasOneData(result)
	if err != nil {
		return nil, errors.Wrap(err, "ggm.Select.hasOneData")
	}

	result, err = m.hasManyData(result)
	if err != nil {
		return nil, errors.Wrap(err, "ggm.Select.hasManyData")
	}

	return result, nil
}

func (m *model[T]) SelectOne(condition ...Option) (row T, err error) {
	condition = append(condition, Limit(1), Offset(0))
	result, err := m.Select(condition...)
	if err != nil {
		return row, errors.Wrap(err, "ggm.SelectOne")
	}
	if lenT(result) == 0 {
		return row, ErrNilRow
	}
	return result[0], nil
}

func (m *model[T]) Count(condition ...Option) (count int, err error) {
	var kv []any
	defer dbLog(time.Now(), &err, &kv)
	if err := m.check(); err != nil {
		return 0, err
	}

	condition = append(condition, Table(m.table), AggregateCount("*"))
	if m.fakeDeleteKey != "" {
		condition = append(condition, WhereNotEq(m.fakeDeleteKey, 1))
	}
	_sql, args := SelectBuilder(condition...)
	kv = append(kv, "sql:", _sql, "args:", args)

	result := struct {
		Count int `db:"count"`
	}{}

	err = m.client.Get(&result, _sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "ggm.Count.sqlx_get")
	}

	return result.Count, nil
}

func (m *model[T]) Insert(rows ...T) (id int64, err error) {
	var kv []any
	defer dbLog(time.Now(), &err, &kv)
	if err := m.check(); err != nil {
		return 0, errors.Wrap(err, "ggm.Insert.check")
	}

	if lenT(rows) == 0 {
		return 0, errors.New("ggm.Insert data is empty")
	}

	fields, err := m.tInsertFields()
	if err != nil {
		return 0, errors.Wrap(err, "ggm.Insert.structFields")
	}

	_sql := InsertNamedBuilder(Table(m.table), Field(fields...))
	kv = append(kv, "sql:", _sql, "rows:", rows)

	result, err := m.client.NamedExec(_sql, rows)
	if err != nil {
		return 0, errors.Wrap(err, "ggm.Insert.sqlx_exec")
	}

	insertID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("ggm.Insert id failed, err:%v\n", err)
	}

	return insertID, nil
}

func (m *model[T]) Update(row T, opt ...Option) (affectedRows int64, err error) {
	var kv []any
	defer dbLog(time.Now(), &err, &kv)
	if err := m.check(); err != nil {
		return 0, errors.Wrap(err, "ggm.Update")
	}

	pk := m.pk()
	// if pk not exist and opt is empty,
	// then we can not to update some record
	if pk == "" && lenT(opt) == 0 {
		return 0, errors.New("ggm.Update not found update condition")
	}

	fields, args := m.tFieldValue(row)
	if lenT(fields) == 0 {
		return 0, errors.Wrap(err, "ggm.Update.structFields is empty")
	}

	index := lo.IndexOf(fields, pk)
	// if pk exist and opt is empty, and the struct.{pk} is zeroVal
	// then we can not to update some record
	if pk != "" && lenT(opt) == 0 && index == -1 {
		return 0, errors.New("ggm.Update if update condition is empty, please set the primary key")
	}

	if pk != "" && index > -1 {
		opt = append(opt, WhereEq(pk, args[index]))
		kv = append(kv, "id:", args[index])
		fields = remove(fields, index)
		args = remove(args, index)
	}

	opt = append(opt, Table(m.table), Field(fields...), Value(args...))
	_sql, _args := UpdateBuilder(opt...)

	kv = append(kv, "sql:", _sql, "args:", args)

	result, err := m.Exec(_sql, _args...)
	if err != nil {
		return 0, errors.Wrap(err, "ggm.Update.exec")
	}

	affectedRows, err = result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("ggm.Update get affected failed, err:%v\n", err)
	}

	return affectedRows, nil
}

func (m *model[T]) Delete(opt ...Option) (ok bool, err error) {
	var kv []any
	defer dbLog(time.Now(), &err, &kv)

	if lenT(opt) == 0 {
		return false, errors.New("ggm.Delete condition is empty")
	}
	opt = append(opt, Table(m.table))
	var result sql.Result
	if m.fakeDeleteKey != "" {
		opt = append(opt, Field(m.fakeDeleteKey), Value(1))
		_sql, _args := UpdateBuilder(opt...)
		kv = append(kv, "slq:", _sql, "args:", _args)
		result, err = m.Exec(_sql, _args...)
	} else {
		_sql, _args := DeleteBuilder(opt...)
		kv = append(kv, "slq:", _sql, "args:", _args)
		result, err = m.Exec(_sql, _args...)
	}

	if err != nil {
		return false, errors.Wrap(err, "ggm.Delete exec")
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("ggm.Delete get affected failed, err:%v\n", err)
	}

	return affectedRows > 0, nil
}

func (m *model[T]) Exec(sql string, args ...any) (sql.Result, error) {
	return m.client.Exec(sql, args...)
}

func (m *model[T]) pk() string {
	return m.modelInfo.PrimaryKey
}

func (m *model[T]) tInsertFields() (fields []string, err error) {
	for _, f := range m.modelInfo.Fields {
		if !inArr[string](f.Options, "ii") {
			fields = append(fields, f.Name)
		}
	}
	return fields, nil
}

func (m *model[T]) tFields() (fields []string, err error) {
	for _, f := range m.modelInfo.Fields {
		fields = append(fields, f.Name)
	}
	return fields, nil
}

func (m *model[T]) tFieldValue(t T) (field []string, value []any) {
	rm := reflectx.NewMapper("db")
	fm := rm.FieldMap(reflect.ValueOf(t))

	for k, v := range fm {
		if strings.Contains(k, ".") || v.IsZero() {
			continue
		}
		field = append(field, k)
		value = append(value, v.Interface())
	}

	return field, value
}

func dbLog(start time.Time, err *error, kv *[]any) {
	tc := time.Since(start)
	log := []any{
		"ums:", tc.Milliseconds(),
	}
	log = append(log, *kv...)
	if *err != nil {
		log = append(log, "error:", *err)
		_ = logger.Log(LevelErr, log...)
		return
	}
	_ = logger.Log(LevelDebug, log...)
}
