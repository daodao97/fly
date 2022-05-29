package ggm

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test_reflectInfo(t *testing.T) {
	u := &User{}
	info, err := structInfo(u)
	assert.Equal(t, nil, err)
	spew.Dump(info)
}

func Test_realType(t *testing.T) {
	fmt.Println(getRealType(reflect.TypeOf(&User{})))
	fmt.Println(getRealType(reflect.TypeOf(User{})))
}
