package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

func init() {
	once := sync.Once{}
	once.Do(func() {
		extra.RegisterFuzzyDecoders()
	})
}

var ErrParamsType = errors.New("param record type must be map[string]interface, *map[string]interface, struct, *struct")

func DecodeToMap(s interface{}, saveZero bool) (map[string]interface{}, error) {
	tmp := map[string]interface{}{}
	t := reflect.TypeOf(s)
	if isMapStrInterface(t) {
		return s.(map[string]interface{}), nil
	}

	if isPtrMapStrInterface(t) {
		return *s.(*map[string]interface{}), nil
	}

	v := reflect.Indirect(reflect.ValueOf(s))
	if isStruct(t) || isPtrStruct(t) {
		t = Deref(t)
		for i := 0; i < v.NumField(); i++ {
			f := t.Field(i)
			name := f.Tag.Get("db")
			_v := v.Field(i)
			if !saveZero && _v.IsZero() {
				continue
			}
			tmp[name] = _v.Interface()
		}
		return tmp, nil
	}

	return nil, ErrParamsType
}

func Binding(from interface{}, to interface{}) error {
	switch from := from.(type) {
	case []byte:
		return jsoniter.Unmarshal(from, to)
	case string:
		if from == "" {
			return fmt.Errorf("the source data is empty string")
		}
		return jsoniter.UnmarshalFromString(from, to)
	case io.ReadCloser:
		body, err := ioutil.ReadAll(from)
		if err != nil {
			return err
		}
		return Binding(body, to)
	default:
		tmp, err := jsoniter.Marshal(from)
		if err != nil {
			return err
		}
		err = jsoniter.Unmarshal(tmp, to)
		if err != nil {
			return err
		}
		return nil
	}
}

func Decoder(source, dest interface{}) error {
	_decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           dest,
		WeaklyTypedInput: true,
		TagName:          "db",
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeHookFunc(time.RFC3339),
		),
	})
	if err != nil {
		return err
	}

	err = _decoder.Decode(source)
	if err != nil {
		return err
	}

	return nil
}
