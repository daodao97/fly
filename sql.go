package ggm

import (
	"fmt"
	"strings"
)

const selectMod = "select %s from `%s`"
const insertMod = "insert into %s (%s) values (%s)"
const updateMod = "update %s set %s"
const deleteMod = "delete from %s"

type OrderByType string

const OrderByASC OrderByType = "ASC"
const OrderByDESC OrderByType = "DESC"

type Option = func(opts *Options)

type Options struct {
	database string
	table    string
	field    []string
	where    []where
	orderBy  string
	groupBy  string
	limit    int
	offset   int
	value    []interface{}
}

func Table(table string) Option {
	return func(opts *Options) {
		opts.table = table
	}
}

func Database(database string) Option {
	return func(opts *Options) {
		opts.database = database
	}
}

func Offset(offset int) Option {
	return func(opts *Options) {
		opts.offset = offset
	}
}

func Limit(offset int) Option {
	return func(opts *Options) {
		opts.limit = offset
	}
}

func Pagination(pageNumber, pageSize int) []Option {
	return []Option{
		Limit(pageSize),
		Offset((pageNumber - 1) * pageSize),
	}
}

func Field(name ...string) Option {
	var _name []string
	for _, v := range name {
		if strings.Contains(v, " as ") {
			tmp := strings.Split(v, " as ")
			_name = append(_name, "`"+strings.Trim(tmp[0], " ")+"` as `"+strings.Trim(tmp[1], " ")+"`")
		} else if strings.Contains(v, " AS ") {
			tmp := strings.Split(v, " AS ")
			_name = append(_name, "`"+strings.Trim(tmp[0], " ")+"` as `"+strings.Trim(tmp[1], " ")+"`")
			_name = append(_name, "`"+v+"`")
		} else {
			_name = append(_name, "`"+v+"`")
		}
	}
	return func(opts *Options) {
		opts.field = _name
	}
}

func FieldRaw(name string) Option {
	return func(opts *Options) {
		opts.field = append(opts.field, name)
	}
}

func AggregateSum(name string) Option {
	return Field("sum(" + name + ") as aggregate")
}

func AggregateCount(name string) Option {
	return FieldRaw("count(" + name + ") as count")
}

func AggregateMax(name string) Option {
	return Field("max(" + name + ") as aggregate")
}

func Value(val ...interface{}) Option {
	return func(opts *Options) {
		opts.value = val
	}
}

type where struct {
	field    string
	operator string
	value    interface{}
	logic    string
	sub      []where
}

func Where(field, operator string, value interface{}) Option {
	return func(opts *Options) {
		opts.where = append(opts.where, where{
			field:    field,
			operator: operator,
			value:    value,
			logic:    "and",
		})
	}
}

func WhereEq(field string, value interface{}) Option {
	return func(opts *Options) {
		opts.where = append(opts.where, where{
			field:    field,
			operator: "=",
			value:    value,
		})
	}
}

func WhereNotEq(field string, value interface{}) Option {
	return func(opts *Options) {
		opts.where = append(opts.where, where{
			field:    field,
			operator: "!=",
			value:    value,
		})
	}
}

func WhereGt(field string, value interface{}) Option {
	return func(opts *Options) {
		opts.where = append(opts.where, where{
			field:    field,
			operator: ">",
			value:    value,
		})
	}
}

func WhereGe(field string, value interface{}) Option {
	return func(opts *Options) {
		opts.where = append(opts.where, where{
			field:    field,
			operator: ">=",
			value:    value,
		})
	}
}

func WhereLt(field string, value interface{}) Option {
	return func(opts *Options) {
		opts.where = append(opts.where, where{
			field:    field,
			operator: "<",
			value:    value,
		})
	}
}

func WhereLe(field string, value interface{}) Option {
	return func(opts *Options) {
		opts.where = append(opts.where, where{
			field:    field,
			operator: "<=",
			value:    value,
		})
	}
}

func WhereIn(field string, value []interface{}) Option {
	return func(opts *Options) {
		opts.where = append(opts.where, where{
			field:    field,
			operator: "in",
			value:    value,
		})
	}
}

func WhereNotIn(field string, value []interface{}) Option {
	return func(opts *Options) {
		opts.where = append(opts.where, where{
			field:    field,
			operator: "not in",
			value:    value,
		})
	}
}

func WhereOr(field, operator string, value interface{}) Option {
	return func(opts *Options) {
		opts.where = append(opts.where, where{
			field:    field,
			operator: operator,
			value:    value,
			logic:    "or",
		})
	}
}

func WhereOrEq(field string, value interface{}) Option {
	return func(opts *Options) {
		opts.where = append(opts.where, where{
			field:    field,
			operator: "=",
			value:    value,
			logic:    "or",
		})
	}
}

func WhereOrNotEq(field string, value interface{}) Option {
	return func(opts *Options) {
		opts.where = append(opts.where, where{
			field:    field,
			operator: "!=",
			value:    value,
			logic:    "or",
		})
	}
}

func WhereOrGt(field string, value interface{}) Option {
	return func(opts *Options) {
		opts.where = append(opts.where, where{
			field:    field,
			operator: ">",
			value:    value,
			logic:    "or",
		})
	}
}

func WhereOrGe(field string, value interface{}) Option {
	return func(opts *Options) {
		opts.where = append(opts.where, where{
			field:    field,
			operator: ">=",
			value:    value,
			logic:    "or",
		})
	}
}

func WhereOrLt(field string, value interface{}) Option {
	return func(opts *Options) {
		opts.where = append(opts.where, where{
			field:    field,
			operator: "<",
			value:    value,
			logic:    "or",
		})
	}
}

func WhereOrLe(field string, value interface{}) Option {
	return func(opts *Options) {
		opts.where = append(opts.where, where{
			field:    field,
			operator: "<=",
			value:    value,
			logic:    "or",
		})
	}
}

func WhereOrIn(field string, value []interface{}) Option {
	return func(opts *Options) {
		opts.where = append(opts.where, where{
			field:    field,
			operator: "in",
			value:    value,
			logic:    "or",
		})
	}
}

func WhereOrNotIn(field string, value []interface{}) Option {
	return func(opts *Options) {
		opts.where = append(opts.where, where{
			field:    field,
			operator: "not in",
			value:    value,
			logic:    "or",
		})
	}
}

func WhereGroup(opts ...Option) Option {
	opt := &Options{}
	for _, v := range opts {
		v(opt)
	}
	return func(opts *Options) {
		opts.where = append(opts.where, where{
			logic: "and",
			sub:   opt.where,
		})
	}
}

func WhereOrGroup(opts ...Option) Option {
	opt := &Options{}
	for _, v := range opts {
		v(opt)
	}
	return func(opts *Options) {
		opts.where = append(opts.where, where{
			logic: "or",
			sub:   opt.where,
		})
	}
}

func WhereLike(field string, value interface{}) Option {
	return func(opts *Options) {
		opts.where = append(opts.where, where{
			field:    field,
			operator: "like",
			value:    value,
		})
	}
}

func WhereOrLike(field string, value interface{}) Option {
	return func(opts *Options) {
		opts.where = append(opts.where, where{
			field:    field,
			operator: "like",
			value:    value,
			logic:    "or",
		})
	}
}

func WhereOrNotLike(field string, value interface{}) Option {
	return func(opts *Options) {
		opts.where = append(opts.where, where{
			field:    field,
			operator: "not like",
			value:    value,
		})
	}
}

func WhereBetween(field string, value interface{}) Option {
	return func(opts *Options) {
		opts.where = append(opts.where, where{
			field:    field,
			operator: "between",
			value:    value,
		})
	}
}

func WhereFindInSet(field string, value interface{}) Option {
	return func(opts *Options) {
		opts.where = append(opts.where, where{
			field:    field,
			operator: "find_in_set",
			value:    value,
		})
	}
}

func WhereOrFindInSet(field string, value interface{}) Option {
	return func(opts *Options) {
		opts.where = append(opts.where, where{
			field:    field,
			operator: "find_in_set",
			value:    value,
			logic:    "or",
		})
	}
}

func OrderBy(field string, mod OrderByType) Option {
	return func(opts *Options) {
		opts.orderBy = "`" + field + "` " + string(mod)
	}
}

func GroupBy(field string) Option {
	return func(opts *Options) {
		opts.groupBy = field
	}
}

func whereBuilder(condition []where) (sql string, args []interface{}) {
	if lenT(condition) == 0 {
		return "", nil
	}
	var tokens []string
	for i, v := range condition {
		if i != 0 {
			if v.logic != "" {
				tokens = append(tokens, v.logic)
			} else {
				tokens = append(tokens, "and")
			}
		}

		if v.field != "" {
			switch v.operator {
			case "in", "not in":
				val := v.value.([]interface{})
				var placeholder []string
				for range val {
					placeholder = append(placeholder, "?")
				}
				tokens = append(tokens, fmt.Sprintf("`%s` %s (%s)", v.field, v.operator, strings.Join(placeholder, ",")))
				args = append(args, val...)
			case "between":
				val := v.value.([]interface{})
				tokens = append(tokens, fmt.Sprintf("`%s` %s ? and ?", v.field, v.operator))
				args = append(args, val...)
			case "find_in_set":
				tokens = append(tokens, fmt.Sprintf("find_in_set(?, %s)", v.field))
				args = append(args, v.value)
			default:
				tokens = append(tokens, fmt.Sprintf("`%s` %s ?", v.field, v.operator))
				args = append(args, v.value)
			}
		}

		if v.sub != nil {
			_sql, _args := whereBuilder(v.sub)
			tokens = append(tokens, "("+_sql+")")
			args = append(args, _args...)
		}
	}
	return strings.Join(tokens, " "), args
}

func SelectBuilder(opts ...Option) (sql string, args []interface{}) {
	_opts := &Options{}
	for _, v := range opts {
		v(_opts)
	}

	_where, args := whereBuilder(_opts.where)
	_field := "*"

	if lenT(_opts.field) > 0 {
		_field = strings.Join(_opts.field, ", ")
	}

	sql = fmt.Sprintf(selectMod, _field, _opts.table)

	if _where != "" {
		sql = sql + " where " + _where
	}

	if _opts.orderBy != "" {
		sql = sql + " order by " + _opts.orderBy
	}

	if _opts.groupBy != "" {
		sql = sql + " group by " + _opts.groupBy
	}

	if _opts.limit != 0 {
		sql = sql + " limit ? offset ? "
		args = append(args, _opts.limit, _opts.offset)
	}

	return sql, args
}

func getTable(opt *Options) string {
	if opt.database == "" {
		return opt.table
	}
	return opt.database + "." + opt.table
}

func InsertBuilder(opts ...Option) (sql string, args []interface{}) {
	_opts := &Options{}
	for _, v := range opts {
		v(_opts)
	}
	var _val []string
	for range _opts.field {
		_val = append(_val, "?")
	}
	sql = fmt.Sprintf(insertMod, getTable(_opts), strings.Join(_opts.field, ", "), strings.Join(_val, ","))
	args = _opts.value
	return sql, args
}

func InsertNamedBuilder(opts ...Option) (sql string) {
	_opts := &Options{}
	for _, v := range opts {
		v(_opts)
	}
	var _val []string
	for _, k := range _opts.field {
		_val = append(_val, ":"+strings.ReplaceAll(k, "`", ""))
	}
	return fmt.Sprintf(insertMod, getTable(_opts), strings.Join(_opts.field, ", "), strings.Join(_val, ", "))
}

func UpdateBuilder(opts ...Option) (sql string, args []interface{}) {
	_opts := &Options{}
	for _, v := range opts {
		v(_opts)
	}
	var _val []string
	for _, v := range _opts.field {
		_val = append(_val, v+" = ?")
	}
	sql = fmt.Sprintf(updateMod, getTable(_opts), strings.Join(_val, ","))
	args = _opts.value
	if lenT(_opts.where) > 0 {
		_where, _args := whereBuilder(_opts.where)
		sql = sql + " where " + _where
		args = append(args, _args...)
	}
	return sql, args
}

func DeleteBuilder(opts ...Option) (sql string, args []interface{}) {
	_opts := &Options{}
	for _, v := range opts {
		v(_opts)
	}
	sql = fmt.Sprintf(deleteMod, _opts.table)
	if lenT(_opts.where) > 0 {
		_where, _args := whereBuilder(_opts.where)
		sql = sql + " where " + _where
		args = append(args, _args...)
	}
	return sql, args
}
