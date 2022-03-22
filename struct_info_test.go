package ggm

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test_reflectInfo(t *testing.T) {
	info, err := structInfo[User]()
	assert.Equal(t, nil, err)
	JsonDump(info)
}

func Test_realType(t *testing.T) {
	fmt.Println(getRealType(reflect.TypeOf(&User{})))
	fmt.Println(getRealType(reflect.TypeOf(User{})))
}
