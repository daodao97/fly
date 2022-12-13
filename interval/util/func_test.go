package util

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func Test_allowTypes(t *testing.T) {
	a := struct {
		A int
	}{}

	assert.Equal(t, true, AllowType(a, []string{"struct"}))

	b := &struct {
		A int
	}{}

	assert.Equal(t, true, AllowType(b, []string{"*struct"}))

	c := map[string]interface{}{
		"a": 1,
	}
	assert.Equal(t, true, AllowType(c, []string{"map[string]interface"}))

	d := &map[string]interface{}{
		"a": 1,
	}
	assert.Equal(t, true, AllowType(d, []string{"*map[string]interface"}))

	d1 := &map[interface{}]interface{}{
		"a": 1,
	}
	assert.Equal(t, false, AllowType(d1, []string{"*map[string]interface"}))

	e := 123
	assert.Equal(t, false, AllowType(e, []string{"map[string]interface"}))

	f := "244"
	assert.Equal(t, false, AllowType(f, []string{"map[string]interface"}))
}

func Test_struct2Map(t *testing.T) {
	type B struct {
		B int
	}
	a := struct {
		A string
		B B
	}{
		A: "a",
		B: B{B: 123},
	}

	m, err := DecodeToMap(a, false)
	assert.Equal(t, nil, err)
	spew.Dump(m)
}

func Test_Decoder(t *testing.T) {
	type A struct {
		F1 [][]int   `json:"f_1,omitempty"`
		T  time.Time `json:"t"`
	}

	data := map[string]interface{}{
		"f_1": [][]int{{1}, {2}},
		"t":   time.Now(),
	}
	var a A
	err := Decoder(data, &a)
	assert.Equal(t, nil, err)
	spew.Dump(a)

	_data, _ := json.Marshal(data)
	_ = json.Unmarshal(_data, &a)
	spew.Dump(a)
}
