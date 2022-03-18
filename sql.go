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

const ASC OrderByType = "ASC"
const DESC OrderByType = "DESC"

type Option interface {
	apply(opts *Options)
}

type Options struct {
	table   string
	field   []string
	where   []where
	orderBy string
	groupBy string
	limit   int
	offset  int
	value   []any
}

type table string

func (t *table) apply(opts *Options) {
	opts.table = string(*t)
}

func Table(name string) Option {
	t := table(name)
	return &t
}

type Offset int

func (o Offset) apply(opts *Options) {
	opts.offset = int(o)
}

type Limit int

func (o Limit) apply(opts *Options) {
	opts.limit = int(o)
}

func Pagination(pageNumber, pageSize int) []Option {
	return []Option{
		Limit(pageSize),
		Offset((pageNumber - 1) * pageSize),
	}
}

type fields []string

func (o fields) apply(opts *Options) {
	opts.field = o
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
	fs := fields(_name)
	return &fs
}

func FieldRaw(name string) Option {
	var _name []string
	_name = append(_name, name)
	fs := fields(_name)
	return &fs
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

type value []any

func (o value) apply(opts *Options) {
	opts.value = o
}

func Value(val ...any) Option {
	v := value(val)
	return &v
}

type where struct {
	field    string
	operator string
	value    any
	logic    string
	sub      []where
}

func (o *where) apply(opts *Options) {
	opts.where = append(opts.where, where{
		field:    o.field,
		operator: o.operator,
		value:    o.value,
		logic:    o.logic,
		sub:      o.sub,
	})
}

func Where(field, operator string, value any) Option {
	return &where{
		field:    field,
		operator: operator,
		value:    value,
	}
}

func WhereEq(field string, value any) Option {
	return &where{
		field:    field,
		operator: "=",
		value:    value,
	}
}

func WhereNotEq(field string, value any) Option {
	return &where{
		field:    field,
		operator: "!=",
		value:    value,
	}
}

func WhereGt(field string, value any) Option {
	return &where{
		field:    field,
		operator: ">",
		value:    value,
	}
}

func WhereGe(field string, value any) Option {
	return &where{
		field:    field,
		operator: ">=",
		value:    value,
	}
}

func WhereLt(field string, value any) Option {
	return &where{
		field:    field,
		operator: "<",
		value:    value,
	}
}

func WhereLe(field string, value any) Option {
	return &where{
		field:    field,
		operator: "<=",
		value:    value,
	}
}

func WhereIn(field string, value []any) Option {
	return &where{
		field:    field,
		operator: "in",
		value:    value,
	}
}

func WhereNotIn(field string, value []any) Option {
	return &where{
		field:    field,
		operator: "not in",
		value:    value,
	}
}

func WhereOr(field, operator string, value any) Option {
	return &where{
		field:    field,
		operator: operator,
		value:    value,
		logic:    "or",
	}
}

func WhereOrEq(field string, value any) Option {
	return &where{
		field:    field,
		operator: "=",
		value:    value,
		logic:    "or",
	}
}

func WhereOrNotEq(field string, value any) Option {
	return &where{
		field:    field,
		operator: "!=",
		value:    value,
		logic:    "or",
	}
}

func WhereOrGt(field string, value any) Option {
	return &where{
		field:    field,
		operator: ">",
		value:    value,
		logic:    "or",
	}
}

func WhereOrGe(field string, value any) Option {
	return &where{
		field:    field,
		operator: ">=",
		value:    value,
		logic:    "or",
	}
}

func WhereOrLt(field string, value any) Option {
	return &where{
		field:    field,
		operator: "<",
		value:    value,
		logic:    "or",
	}
}

func WhereOrLe(field string, value any) Option {
	return &where{
		field:    field,
		operator: "<=",
		value:    value,
		logic:    "or",
	}
}

func WhereOrIn(field string, value []any) Option {
	return &where{
		field:    field,
		operator: "in",
		value:    value,
		logic:    "or",
	}
}

func WhereOrNotIn(field string, value []any) Option {
	return &where{
		field:    field,
		operator: "not in",
		value:    value,
		logic:    "or",
	}
}

func WhereGroup(opts ...Option) Option {
	opt := &Options{}
	for _, v := range opts {
		switch v.(type) {
		case *where:
			v.apply(opt)
		}
	}
	return &where{
		logic: "and",
		sub:   opt.where,
	}
}

func WhereOrGroup(opts ...Option) Option {
	opt := &Options{}
	for _, v := range opts {
		switch v.(type) {
		case *where:
			v.apply(opt)
		}
	}
	return &where{
		logic: "or",
		sub:   opt.where,
	}
}

func WhereLike(field string, value any) Option {
	return &where{
		field:    field,
		operator: "like",
		value:    value,
	}
}

func WhereOrLike(field string, value any) Option {
	return &where{
		field:    field,
		operator: "like",
		value:    value,
		logic:    "or",
	}
}

func WhereOrNotLike(field string, value any) Option {
	return &where{
		field:    field,
		operator: "not like",
		value:    value,
	}
}

func WhereBetween(field string, value any) Option {
	return &where{
		field:    field,
		operator: "between",
		value:    value,
	}
}

type orderBy string

func (t *orderBy) apply(opts *Options) {
	opts.orderBy = string(*t)
}

func OrderBy(field string, mod OrderByType) Option {
	t := orderBy("`" + field + "` " + string(mod))
	return &t
}

type groupBy string

func (t *groupBy) apply(opts *Options) {
	opts.groupBy = string(*t)
}

func GroupBy(field string) Option {
	t := groupBy(field)
	return &t
}

func whereBuilder(condition []where) (sql string, args []any) {
	if len(condition) == 0 {
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
				val := v.value.([]any)
				var placeholder []string
				for range val {
					placeholder = append(placeholder, "?")
				}
				tokens = append(tokens, fmt.Sprintf("`%s` %s (%s)", v.field, v.operator, strings.Join(placeholder, ",")))
				args = append(args, val...)
			case "between":
				val := v.value.([]any)
				tokens = append(tokens, fmt.Sprintf("`%s` %s ? and ?", v.field, v.operator))
				args = append(args, val...)
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

func SelectBuilder(opts ...Option) (sql string, args []any) {
	_opts := &Options{}
	for _, v := range opts {
		v.apply(_opts)
	}

	_where, args := whereBuilder(_opts.where)
	_field := "*"

	if len(_opts.field) > 0 {
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

func InsertBuilder(opts ...Option) (sql string, args []any) {
	_opts := &Options{}
	for _, v := range opts {
		v.apply(_opts)
	}
	var _val []string
	for range _opts.field {
		_val = append(_val, "?")
	}
	sql = fmt.Sprintf(insertMod, _opts.table, strings.Join(_opts.field, ", "), strings.Join(_val, ","))
	args = _opts.value
	return sql, args
}

func InsertNamedBuilder(opts ...Option) (sql string) {
	_opts := &Options{}
	for _, v := range opts {
		v.apply(_opts)
	}
	var _val []string
	for _, k := range _opts.field {
		_val = append(_val, ":"+strings.ReplaceAll(k, "`", ""))
	}
	return fmt.Sprintf(insertMod, _opts.table, strings.Join(_opts.field, ", "), strings.Join(_val, ", "))
}

func UpdateBuilder(opts ...Option) (sql string, args []any) {
	_opts := &Options{}
	for _, v := range opts {
		v.apply(_opts)
	}
	var _val []string
	for _, v := range _opts.field {
		_val = append(_val, v+" = ?")
	}
	sql = fmt.Sprintf(updateMod, _opts.table, strings.Join(_val, ","))
	args = _opts.value
	if len(_opts.where) > 0 {
		_where, _args := whereBuilder(_opts.where)
		sql = sql + " where " + _where
		args = append(args, _args...)
	}
	return sql, args
}

func DeleteBuilder(opts ...Option) (sql string, args []any) {
	_opts := &Options{}
	for _, v := range opts {
		v.apply(_opts)
	}
	sql = fmt.Sprintf(deleteMod, _opts.table)
	if len(_opts.where) > 0 {
		_where, _args := whereBuilder(_opts.where)
		sql = sql + " where " + _where
		args = append(args, _args...)
	}
	return sql, args
}
