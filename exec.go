package fly

import (
	"database/sql"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// 一般用Prepared Statements和Exec()完成INSERT, UPDATE, DELETE操作
func exec(db *sql.DB, _sql string, args ...interface{}) (res sql.Result, err error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	stmt, err := tx.Prepare(_sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	res, err = stmt.Exec(args...)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func query(db *sql.DB, _sql string, args ...interface{}) (result []Row, err error) {
	stmt, err := db.Prepare(_sql)
	if err != nil {
		return nil, errors.Wrap(err, "fly.exec.Prepare err")
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, errors.Wrap(err, "fly.exec.Query err")
	}
	defer rows.Close()

	return rows2SliceMap(rows)
}

func destination(columnTypes []*sql.ColumnType) func() []interface{} {
	dest := make([]func() interface{}, 0, len(columnTypes))
	for _, v := range columnTypes {
		switch strings.ToUpper(v.DatabaseTypeName()) {
		case "VARCHAR", "CHAR", "TEXT", "NVARCHAR", "LONGTEXT", "LONGBLOB", "MEDIUMTEXT", "MEDIUMBLOB", "BLOB", "TINYTEXT", "DECIMAL":
			if nullable, _ := v.Nullable(); nullable {
				dest = append(dest, func() interface{} {
					return new(sql.NullString)
				})
			} else {
				dest = append(dest, func() interface{} {
					return new(string)
				})
			}
		case "INT", "TINYINT", "INTEGER", "SMALLINT", "MEDIUMINT":
			dest = append(dest, func() interface{} {
				return new(int)
			})
		case "BIGINT":
			dest = append(dest, func() interface{} {
				return new(int64)
			})
		case "DATETIME", "DATE", "TIMESTAMP", "TIME":
			dest = append(dest, func() interface{} {
				return new(time.Time)
			})
		case "DOUBLE", "FLOAT":
			dest = append(dest, func() interface{} {
				return new(float64)
			})
		default:
			dest = append(dest, func() interface{} {
				return new(string)
			})
		}
	}
	return func() []interface{} {
		tmp := make([]interface{}, 0, len(dest))
		for _, d := range dest {
			tmp = append(tmp, d())
		}
		return tmp
	}
}

func rows2SliceMap(rows *sql.Rows) (list []Row, err error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, errors.Wrap(err, "fly.rows2SliceMap.columns err")
	}
	length := len(columns)

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, errors.Wrap(err, "fly.rows2SliceMap.ColumnTypes err")
	}

	dest := destination(columnTypes)

	for rows.Next() {
		tmp := dest()
		err = rows.Scan(tmp...)
		if err != nil {
			return nil, errors.Wrap(err, "fly.rows2SliceMap.Scan err")
		}
		row := new(Row)
		row.Data = map[string]interface{}{}
		for i := 0; i < length; i++ {
			if val, ok := tmp[i].(*sql.NullString); ok {
				row.Data[columns[i]] = val.String
			} else {
				row.Data[columns[i]] = tmp[i]
			}
		}
		list = append(list, *row)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "fly.rows2SliceMap.rows.Err err")
	}
	return list, nil
}
