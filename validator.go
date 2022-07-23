package fly

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/daodao97/fly/interval/xtype"
)

type Valid = func(v *ValidInfo) error

type ValidInfo struct {
	Field string
	Row   map[string]interface{}
	Model Model
	Label string
	Msg   string
}

func mergeOpt(v1, v2 *ValidInfo) *ValidInfo {
	if v2.Row != nil {
		v1.Row = v2.Row
	}
	if v2.Model != nil {
		v1.Model = v2.Model
	}
	if v2.Field != "" {
		v1.Field = v2.Field
	}
	if v2.Label != "" {
		v1.Label = v2.Label
	}
	if v2.Msg != "" {
		v1.Msg = v2.Msg
	}
	return v1
}

func ValidWrap(valid Valid, v1 *ValidInfo) Valid {
	return func(v *ValidInfo) error {
		return valid(mergeOpt(v, v1))
	}
}

type ValidOpt = func(*ValidInfo)

func WithMsg(msg string) ValidOpt {
	return func(v *ValidInfo) {
		v.Msg = msg
	}
}

func withField(field string) ValidOpt {
	return func(v *ValidInfo) {
		v.Field = field
	}
}

func WithLabel(label string) ValidOpt {
	return func(v *ValidInfo) {
		v.Label = label
	}
}

func withRow(row map[string]interface{}) ValidOpt {
	return func(v *ValidInfo) {
		v.Row = row
	}
}

func WithModel(m Model) ValidOpt {
	return func(v *ValidInfo) {
		v.Model = m
	}
}

func ExtendValidOpt(v *ValidInfo, opt ...ValidOpt) *ValidInfo {
	for _, o := range opt {
		o(v)
	}
	return v
}

func NewValidOpt(opt ...ValidOpt) *ValidInfo {
	v := &ValidInfo{}
	for _, o := range opt {
		o(v)
	}
	return v
}

func msg(msg1, msg2 string) string {
	if msg2 != "" {
		return msg2
	}
	return msg1
}

func Required(opt ...ValidOpt) Valid {
	v1 := NewValidOpt(opt...)
	return ValidWrap(func(v *ValidInfo) error {
		val, ok := v.Row[v.Field]
		if !ok {
			return errors.New(msg(fmt.Sprintf("%s not found", v.Field), v.Msg))
		}
		if !xtype.Bool(val) {
			return errors.New(msg(fmt.Sprintf("%s value is zero value", v.Field), v.Msg))
		}
		return nil
	}, v1)
}

func IfRequired(ifField string, opt ...ValidOpt) Valid {
	v1 := NewValidOpt(opt...)
	return ValidWrap(func(v *ValidInfo) error {
		if err := Required()(NewValidOpt(withField(ifField), withRow(v.Row), WithModel(v.Model))); err != nil {
			return nil
		}

		return Required()(NewValidOpt(
			withField(v.Field),
			withRow(v.Row),
			WithModel(v.Model),
			WithMsg(msg(fmt.Sprintf("当 %s 存在时 %s 是必须的", ifField, v.Field), v.Msg)),
		))
	}, v1)
}

func Unique(opt ...ValidOpt) Valid {
	v1 := NewValidOpt(opt...)
	return ValidWrap(func(v *ValidInfo) error {
		opts := []Option{WhereEq(v.Field, v.Row[v.Field])}
		if id, ok := v.Row[v.Model.PrimaryKey()]; ok {
			opts = append(opts, WhereNotEq(v.Model.PrimaryKey(), id))
		}
		count, err := v.Model.Count(opts...)
		if err != nil {
			return err
		}
		if count != 0 {
			return errors.New(msg("Duplicate data", v.Msg))
		}
		return nil
	}, v1)
}
