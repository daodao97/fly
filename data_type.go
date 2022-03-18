package ggm

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

type DataType[T any] interface {
	Value() (driver.Value, error)
	Scan(value any) error
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(b []byte) error
	Get() T
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

func (j *Json[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.data)
}

func (j *Json[T]) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &j.data)
	return err
}

func (j Json[T]) Get() T {
	return j.data
}
