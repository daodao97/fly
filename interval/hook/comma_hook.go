package hook

import (
	"sort"
	"strings"

	"github.com/spf13/cast"
)

type CommaSeparatedInt struct{}

func (CommaSeparatedInt) Input(row map[string]interface{}, fieldValue interface{}) (interface{}, error) {
	if fieldValue == nil {
		return "", nil
	}
	strSlice, err := cast.ToStringSliceE(fieldValue)
	if err != nil {
		return nil, err
	}
	sort.Slice(strSlice, func(i, j int) bool {
		return strSlice[i] < strSlice[j]
	})
	return strings.Join(strSlice, ","), nil
}

func (CommaSeparatedInt) Output(row map[string]interface{}, fieldValue interface{}) (interface{}, error) {
	parts := strings.Split(cast.ToString(fieldValue), ",")
	var _parts []int
	for _, v := range parts {
		if v == "" {
			continue
		}
		_parts = append(_parts, cast.ToInt(v))
	}

	return _parts, nil
}

type CommaSeparatedString struct {
	CommaSeparatedInt
}

func (CommaSeparatedString) Input(row map[string]interface{}, fieldValue interface{}) (interface{}, error) {
	tmp, err := cast.ToStringSliceE(fieldValue)
	if err != nil {
		return nil, err
	}

	return strings.Join(tmp, ","), nil
}

func (CommaSeparatedString) Output(row map[string]interface{}, fieldValue interface{}) (interface{}, error) {
	var tmp []string
	for _, i := range strings.Split(cast.ToString(fieldValue), ",") {
		if i == "" {
			continue
		}
		tmp = append(tmp, i)
	}
	return tmp, nil
}
