package fly

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

type HasOpts struct {
	Conn       string
	Database   string
	Table      string
	LocalKey   string
	ForeignKey string
	OtherKeys  []string
}

var space = regexp.MustCompile(`\s+`)

func filterStr(arr []string) (real []string) {
	for _, v := range arr {
		if v != "" {
			real = append(real, v)
		}
	}
	return real
}

func otherKeys(key []string) (rel []string, err error) {
	for _, v := range key {
		split := filterStr(space.Split(strings.Trim(v, " "), -1))
		l := len(split)
		if l == 1 || l == 3 {
			rel = append(rel, split[l-1])
			continue
		}
		return nil, errors.New("otherKey must be like a, a as b, a AS c")
	}
	return rel, nil
}

func (m *model) hasOneData(rows []Row, opt HasOpts) ([]Row, error) {
	otherKeys, err := otherKeys(opt.OtherKeys)
	if err != nil {
		return nil, err
	}
	var localKeys []interface{}
	for _, v := range rows {
		if val, ok := v.Data[opt.LocalKey]; ok {
			localKeys = append(localKeys, val)
		}
	}

	if len(localKeys) == 0 {
		_ = logger.Log(LevelDebug, "hasOneData empty localKeys", fmt.Sprintf("%+v", opt))
		return rows, nil
	}

	opt.OtherKeys = append(opt.OtherKeys, opt.ForeignKey)

	_rows := New(opt.Table, WithConn(opt.Conn)).Select(Field(opt.OtherKeys...), WhereIn(opt.ForeignKey, localKeys))
	if _rows.Err != nil {
		return nil, errors.Wrap(_rows.Err, "hasOne err")
	}

	for i, left := range rows {
		for _, right := range _rows.List {
			l, err := cast.ToStringE(left.Data[opt.LocalKey])
			if err != nil {
				return nil, err
			}
			r, err := cast.ToStringE(right.Data[opt.ForeignKey])
			if err != nil {
				return nil, err
			}
			if l == r {
				for _, k := range otherKeys {
					rows[i].Data[k] = right.Data[k]
				}
			}
		}
	}

	return rows, nil
}

func (m *model) hasManyData(rows []Row, opt HasOpts) ([]Row, error) {
	otherKeys, err := otherKeys(opt.OtherKeys)
	if err != nil {
		return nil, err
	}
	var localKeys []interface{}
	for _, v := range rows {
		if val, ok := v.Data[opt.LocalKey]; ok {
			localKeys = append(localKeys, val)
		}
	}

	if len(localKeys) == 0 {
		_ = logger.Log(LevelDebug, "hasManyData empty localKeys", fmt.Sprintf("%+v", opt))
		return rows, nil
	}

	opt.OtherKeys = append(opt.OtherKeys, opt.ForeignKey)

	_rows := New(opt.Table, WithConn(opt.Conn)).Select(Field(opt.OtherKeys...), WhereIn(opt.ForeignKey, localKeys))
	if _rows.Err != nil {
		return nil, errors.Wrap(_rows.Err, "hasOne err")
	}

	for i, left := range rows {
		tmp := make(map[string][]interface{})
		for _, k := range otherKeys {
			tmp[k] = make([]interface{}, 0)
		}
		for _, right := range _rows.List {
			l, err := cast.ToStringE(left.Data[opt.LocalKey])
			if err != nil {
				return nil, err
			}
			r, err := cast.ToStringE(right.Data[opt.ForeignKey])
			if err != nil {
				return nil, err
			}
			if l == r {
				for _, k := range otherKeys {
					tmp[k] = append(tmp[k], right.Data[k])
				}
			}
		}
		for k, v := range tmp {
			rows[i].Data[k] = v
		}
	}

	return rows, nil
}
