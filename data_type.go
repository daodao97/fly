package fly

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/pkg/errors"
)

type ExampleStruct struct {
	Name  string
	Value string
}

type ExampleStructSlice []ExampleStruct

func (l ExampleStructSlice) Value() (driver.Value, error) {
	bytes, err := json.Marshal(l)
	return string(bytes), err
}

func (l *ExampleStructSlice) Scan(input interface{}) error {
	switch value := input.(type) {
	case string:
		return json.Unmarshal([]byte(value), l)
	case []byte:
		return json.Unmarshal(value, l)
	default:
		return errors.New("not supported")
	}
}
