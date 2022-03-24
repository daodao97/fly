package ggm

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"time"
)

type DataTypeDB[T any] interface {
	Value() (driver.Value, error) // Variables in a program to db value
	Scan(value any) error         // db value scan to program var
}

type DataTypeJson[T any] interface {
	MarshalJSON() ([]byte, error) // var to json string
	UnmarshalJSON(b []byte) error // json string to var
}

// NewJson struct <==> json_str
// db json string data to struct
func NewJson[T any](data T) *Json[T] {
	return &Json[T]{data: data}
}

type Json[T any] struct {
	data T
}

func (j Json[T]) Value() (driver.Value, error) {
	return json.Marshal(j.data)
}

func (j *Json[T]) Scan(value any) error {
	if value == nil {
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	err := json.Unmarshal(bytes, &j.data)
	return err
}

func (j Json[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.data)
}

func (j Json[T]) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &j.data)
	return err
}

func (j Json[T]) Get() T {
	return j.data
}

const TimeFormat = "2006-01-02 15:04:05"

type Time time.Time

func (t Time) Value() (driver.Value, error) {
	return t, nil
}

func (t *Time) Scan(value any) error {
	v, ok := value.(time.Time)
	if ok {
		*t = Time(v)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (t Time) MarshalJSON() ([]byte, error) {
	_t := (time.Time)(t)
	if _t.IsZero() {
		return []byte("null"), nil
	}
	formatted := fmt.Sprintf("\"%s\"", _t.Format(TimeFormat))
	return []byte(formatted), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	if len(data) == 2 {
		*t = Time(time.Time{})
		return nil
	}
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now, err := time.ParseInLocation(`"`+TimeFormat+`"`, string(data), loc)
	*t = Time(now)
	return err
}

func (t *Time) Get() *time.Time {
	return (*time.Time)(t)
}

type CommaSlice[T int | string] []T

func (t CommaSlice[T]) Value() (driver.Value, error) {
	return Join(t), nil
}

func (t *CommaSlice[T]) Scan(value any) error {
	*t = []T{}
	if value == nil {
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	v := string(bytes)
	if v != "" {
		*t = Split[T](v)
	}
	return nil
}

func (t *CommaSlice[T]) Get() []T {
	return *t
}
