package ggm

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	dynamicstruct "github.com/ompluscator/dynamic-struct"
	"github.com/pkg/errors"
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

func New(t TableName) *model {
	m := &model{}

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
	m.connName = conn
	m.conn(conn)

	info, err := structInfo(t)
	if err != nil {
		m.err = err
	}
	m.modelInfo = info

	return m
}

type model struct {
	connName      string
	client        *sqlx.DB
	err           error
	table         string
	fakeDeleteKey string
	modelInfo     *modelInfo
}

func (m *model) DB() *sqlx.DB {
	return m.client
}

func (m *model) check() error {
	if m.err != nil {
		return m.err
	}
	if m.client == nil {
		return fmt.Errorf("model.Client is nil")
	}
	return nil
}

func (m *model) conn(name string) *model {
	client, exist := xdb(name)
	m.client = client
	if !exist {
		m.err = fmt.Errorf("connection [%s] not exist", name)
	}
	return m
}

func (m *model) Select(dest interface{}, condition ...Option) (err error) {
	var kv []interface{}
	defer dbLog("Select", time.Now(), &err, &kv)
	fields, err := m.tFields()
	if err != nil {
		return errors.Wrap(err, "ggm.Select.structFields")
	}
	condition = append(condition, Table(m.table), Field(fields...))
	if m.fakeDeleteKey != "" {
		condition = append(condition, WhereEq(m.fakeDeleteKey, 0))
	}
	_sql, args := SelectBuilder(condition...)
	kv = append(kv, "sql:", _sql, "args:", args)

	err = m.client.Select(dest, _sql, args...)
	if err != nil {
		return err
	}

	dest, err = m.hasOneData(dest)
	if err != nil {
		return err
	}

	dest, err = m.hasManyData(dest)
	if err != nil {
		return err
	}

	return nil
}

func (m *model) SelectOne(dest interface{}, condition ...Option) (err error) {
	condition = append(condition, Limit(1), Offset(0))
	result := dynamicstruct.ExtendStruct(dest).Build().NewSliceOfStructs()
	err = m.Select(result, condition...)
	if err != nil {
		return errors.Wrap(err, "ggm.SelectOne")
	}
	t := reflect.Indirect(reflect.ValueOf(result))
	if t.Len() == 0 {
		return ErrNilRow
	}

	err = copier.Copy(dest, t.Index(0).Interface())
	if err != nil {
		return errors.Wrap(err, "ggm.SelectOne.copy")
	}
	return nil
}

func (m *model) Count(condition ...Option) (count int, err error) {
	var kv []interface{}
	defer dbLog("Count", time.Now(), &err, &kv)
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

func (m *model) Insert(rows ...interface{}) (id int64, err error) {
	var kv []interface{}
	defer dbLog("Insert", time.Now(), &err, &kv)
	if err := m.check(); err != nil {
		return 0, errors.Wrap(err, "ggm.Insert.check")
	}

	if len(rows) == 0 {
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

func (m *model) Update(row interface{}, opt ...Option) (affectedRows int64, err error) {
	var kv []interface{}
	defer dbLog("Update", time.Now(), &err, &kv)
	if err := m.check(); err != nil {
		return 0, errors.Wrap(err, "ggm.Update")
	}

	pk := m.pk()
	// if pk not exist and opt is empty,
	// then we can not to update some record
	if pk == "" && len(opt) == 0 {
		return 0, errors.New("ggm.Update not found update condition")
	}

	fields, args := m.tFieldValue(row)
	if len(fields) == 0 {
		return 0, errors.Wrap(err, "ggm.Update.structFields is empty")
	}

	index := indexOf(fields, pk)
	// if pk exist and opt is empty, and the struct.{pk} is zeroVal
	// then we can not to update some record
	if pk != "" && len(opt) == 0 && index == -1 {
		return 0, errors.New("ggm.Update if update condition is empty, please set the primary key")
	}

	if pk != "" && index > -1 {
		opt = append(opt, WhereEq(pk, args[index]))
		kv = append(kv, "id:", args[index])
		fields = remove(fields, index)
		args = removeInterface(args, index)
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

func (m *model) Delete(opt ...Option) (ok bool, err error) {
	var kv []interface{}
	defer dbLog("Delete", time.Now(), &err, &kv)

	if len(opt) == 0 {
		return false, errors.New("ggm.Delete condition is empty")
	}
	opt = append(opt, Table(m.table))
	var result sql.Result
	if m.fakeDeleteKey != "" {
		opt = append(opt, Field(m.fakeDeleteKey), Value(1))
		_sql, _args := UpdateBuilder(opt...)
		kv = append(kv, "sql:", _sql, "args:", _args)
		result, err = m.Exec(_sql, _args...)
	} else {
		_sql, _args := DeleteBuilder(opt...)
		kv = append(kv, "sql:", _sql, "args:", _args)
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

func (m *model) tFieldValue(t interface{}) (field []string, value []interface{}) {
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

func (m *model) tFields() (fields []string, err error) {
	for _, f := range m.modelInfo.Fields {
		fields = append(fields, f.Name)
	}
	return fields, nil
}

func (m *model) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return m.client.Exec(sql, args...)
}

func (m *model) pk() string {
	return m.modelInfo.PrimaryKey
}

func (m *model) tInsertFields() (fields []string, err error) {
	for _, f := range m.modelInfo.Fields {
		if !inArr(f.Options, "ii") {
			fields = append(fields, f.Name)
		}
	}
	return fields, nil
}

func dbLog(prefix string, start time.Time, err *error, kv *[]interface{}) {
	tc := time.Since(start)
	log := []interface{}{
		prefix,
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

func dbLog2(prefix string, start time.Time, end time.Time, err *error, kv *[]interface{}) {
	log := []interface{}{
		prefix,
		"ums:", end.UnixMilli() - start.UnixMilli(),
	}
	log = append(log, *kv...)
	if *err != nil {
		log = append(log, "error:", *err)
		_ = logger.Log(LevelErr, log...)
		return
	}
	_ = logger.Log(LevelDebug, log...)
}
