package xtype

import (
	"encoding/json"
	"math"
	"reflect"
	"strconv"
	"strings"
)

type Type int

const (
	STRING Type = iota
	NUMBER
	BOOL
	MAP
	ARRAY
	NULL
	UNKNOWN
)

func (t Type) String() string {
	switch t {
	case STRING:
		return "STRING"
	case NUMBER:
		return "NUMBER"
	case BOOL:
		return "BOOL"
	case MAP:
		return "MAP"
	case ARRAY:
		return "ARRAY"
	case NULL:
		return "NULL"
	default:
		return "UNKNOWN"
	}
}

//Go's integer types are: uint8 , uint16 , uint32 , uint64 , int8 , int16 , int32 and int64. 8, 16, 32 and 64 tell us how many bits each of the types use. uint means “unsigned integer” while int means “signed integer”. Unsigned integers only contain positive numbers (or zero).
var _numberTypes = map[string]bool{
	"uint8":   true,
	"uint16":  true,
	"uint32":  true,
	"uint64":  true,
	"int8":    true,
	"int16":   true,
	"int32":   true,
	"int64":   true,
	"uint":    true,
	"int":     true,
	"float32": true,
	"float64": true,
}

// GetType get xtype.Type of given object
func GetType(obj interface{}) Type {
	if obj == nil {
		return NULL
	}

	typeName := reflect.TypeOf(obj).String()
	if strings.HasPrefix(typeName, "[") {
		return ARRAY
	}

	if strings.HasPrefix(typeName, "map") {
		return MAP
	}

	if _, ok := _numberTypes[typeName]; ok {
		return NUMBER
	}

	if typeName == "string" {
		return STRING
	} else if typeName == "bool" {
		return BOOL
	}

	return UNKNOWN
}

// Float try to get float value of given object.
// number: type cast
// bool: true => 1.0 false => 0.0
// string: empty => 0, otherwise parse float
// nil && otherwise: 0.0
func Float(obj interface{}) float64 {
	if obj == nil {
		return 0.0
	}

	switch v := obj.(type) {
	case bool:
		if v {
			return 1.0
		} else {
			return 0.0
		}
	case float32:
		return float64(v)
	case float64:
		return v
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint8:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	case json.Number:
		vv, _ := json.Number(v).Float64()
		return vv
	case string:
		if v == "" {
			return 0.0
		}
		if n, err := strconv.ParseFloat(v, 64); err != nil {
			return 0.0
		} else {
			return n
		}
	default:
		return 0.0
	}
}

// Int try to get int64 value of given object.
// number : type cast
// bool: true => 1, false => 0,
// string: empty => 0, otherwise parse int
// nil & otherwise : 0
func Int(obj interface{}) int64 {
	if obj == nil {
		return 0
	}

	switch v := obj.(type) {
	case bool:
		if v {
			return 1
		} else {
			return 0
		}
	case float32:
		return int64(v)
	case float64:
		return int64(v)
	case int:
		return int64(v)
	case int8:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return v
	case uint:
		return int64(v)
	case uint8:
		return int64(v)
	case uint16:
		return int64(v)
	case uint32:
		return int64(v)
	case uint64:
		return int64(v)
	case json.Number:
		vv, _ := json.Number(v).Int64()
		return int64(vv)
	case string:
		if v == "" {
			return 0
		}
		if n, err := strconv.ParseInt(v, 10, 64); err != nil {
			return 0
		} else {
			return n
		}
	default:
		return 0
	}
}

// Uint try to get uint64 value of given object.
// number : type cast
// bool: true => 1, false => 0,
// string: empty => 0, otherwise parse int
// nil & otherwise : 0
func Uint(obj interface{}) uint64 {
	if obj == nil {
		return 0
	}

	switch v := obj.(type) {
	case bool:
		if v {
			return 1
		} else {
			return 0
		}
	case float32:
		return uint64(v)
	case float64:
		return uint64(v)
	case int:
		return uint64(v)
	case int8:
		return uint64(v)
	case int16:
		return uint64(v)
	case int32:
		return uint64(v)
	case int64:
		return uint64(v)
	case uint64:
		return v
	case uint:
		return uint64(v)
	case uint8:
		return uint64(v)
	case uint16:
		return uint64(v)
	case uint32:
		return uint64(v)
	case json.Number:
		vv, _ := json.Number(v).Int64()
		return uint64(vv)
	case string:
		if v == "" {
			return 0
		}
		if n, err := strconv.ParseUint(v, 10, 64); err != nil {
			return 0
		} else {
			return n
		}
	default:
		return 0
	}
}

// String try to get string value of given object.
// string/[]byte : original string
// bool : true => "1", false => ""
// number : number format
// nil & otherwise : ""
func String(obj interface{}) string {
	if obj == nil {
		return ""
	}

	switch v := obj.(type) {
	case bool:
		if v {
			return "1"
		} else {
			return ""
		}
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatInt(int64(v), 10)
	case uint8:
		return strconv.FormatInt(int64(v), 10)
	case uint16:
		return strconv.FormatInt(int64(v), 10)
	case uint32:
		return strconv.FormatInt(int64(v), 10)
	case uint64:
		return strconv.FormatInt(int64(v), 10)
	case []byte:
		return string(v)
	case string:
		return v
	default:
		return ""
	}

}

// Bytes try to get []byte value of given object
// see String()
func Bytes(obj interface{}) []byte {
	return []byte(String(obj))
}

const EPSILON float64 = 1e-9
const FALSE_STRINGS = "no,false,off,0,"

// Bool try to get bool value of given object.
// number: 0 => false, otherwise => true
// string: ("", "false", "off", "no", "0") => false (case insensitive), otherwise => true
// nil: false
// otherwise: true
func Bool(obj interface{}) bool {
	if obj == nil {
		return false
	}

	switch v := obj.(type) {
	case bool:
		return v
	case float32:
		return math.Abs(float64(v)) > EPSILON
	case float64:
		return math.Abs(v) > EPSILON
	case int:
		return v != 0
	case int8:
		return v != 0
	case int16:
		return v != 0
	case int32:
		return v != 0
	case int64:
		return v != 0
	case uint:
		return v != 0
	case uint8:
		return v != 0
	case uint16:
		return v != 0
	case uint32:
		return v != 0
	case uint64:
		return v != 0
	case []byte:
		s := strings.ToLower(string(v))
		return s != "" && strings.Index(FALSE_STRINGS, s+",") == -1
	case string:
		s := strings.ToLower(v)
		return s != "" && strings.Index(FALSE_STRINGS, s+",") == -1
	default:
		return true
	}
}
