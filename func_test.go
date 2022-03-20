package ggm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func JsonString(arg interface{}) string {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	var err = encoder.Encode(arg)
	if err != nil {
		return ""
	}
	bs := buffer.Bytes()
	var out bytes.Buffer
	err = json.Indent(&out, bs, "", "  ")
	if err != nil {
		return "error"
	}
	return out.String()
}

func JsonDump(args ...interface{}) {
	fmt.Println("======DEBUG=======")
	for _, v := range args {
		fmt.Println(JsonString(v))
	}
	fmt.Println("======DEBUG=======")
}

func Test_reflectValue(t *testing.T) {
	re, err := structFields[User]()
	assert.Equal(t, nil, err)
	fmt.Println(re)

	re1, err := structFields[*User]()
	assert.Equal(t, nil, err)
	fmt.Println(re1)
}

func Test_reflectNew(t *testing.T) {
	re := reflectNew[User]()
	fmt.Println(re.(TableName).Table())

	re1 := reflectNew[*User]()
	fmt.Println(re1.(TableName).Table())
}
