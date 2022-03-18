package ggm

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
