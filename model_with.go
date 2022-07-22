package fly

import (
	"database/sql"
)

type With = func(*model)

func WithDB(db *sql.DB) With {
	return func(b *model) {
		b.client = db
	}
}

func WithConn(name string) With {
	return func(b *model) {
		b.connection = name
	}
}

func WithFakeDelKey(name string) With {
	return func(b *model) {
		b.fakeDelKey = name
	}
}

func WithPrimaryKey(name string) With {
	return func(b *model) {
		b.primaryKey = name
	}
}

func ColumnHook(columnHook ...Hook) With {
	return func(b *model) {
		if b.columnHook == nil {
			b.columnHook = make(map[string]HookData)
		}
		for _, v := range columnHook {
			f, h := v()
			b.columnHook[f] = h
		}
	}
}

// ColumnValidator while validate data by validator when create or update event
func ColumnValidator(validator ...[]Valid) With {
	return func(b *model) {
		if b.columnValidator == nil {
			b.columnValidator = make([]Valid, 0, len(validator))
		}
		for _, v := range validator {
			b.columnValidator = append(b.columnValidator, v...)
		}
	}
}

func Validate(field string, vf ...Valid) (v []Valid) {
	for _, each := range vf {
		v = append(v, ValidWrap(each, NewValidOpt(withField(field))))
	}
	return v
}

func WithSaveZero() With {
	return func(b *model) {
		b.saveZero = true
	}
}

func HasOne(opts ...HasOpts) With {
	return func(b *model) {
		if b.hasOne == nil {
			b.hasOne = make([]HasOpts, 0)
		}
		for i, v := range opts {
			if v.Conn == "" {
				opts[i].Conn = "default"
			}
			if v.LocalKey == "" {
				opts[i].LocalKey = "id"
			}
			if v.ForeignKey == "" {
				opts[i].ForeignKey = "id"
			}
		}
		b.hasOne = append(b.hasOne, opts...)
	}
}

func HasMany(opts ...HasOpts) With {
	return func(b *model) {
		if b.hasMany == nil {
			b.hasMany = make([]HasOpts, 0)
		}
		for i, v := range opts {
			if v.Conn == "" {
				opts[i].Conn = "default"
			}
			if v.LocalKey == "" {
				opts[i].LocalKey = "id"
			}
			if v.ForeignKey == "" {
				opts[i].ForeignKey = "id"
			}
		}
		b.hasMany = append(b.hasMany, opts...)
	}
}
