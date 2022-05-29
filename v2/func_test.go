package ggm

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_structFields(t *testing.T) {
	u := &User{}
	f, err := structFields(u)
	assert.Equal(t, nil, err)
	spew.Dump(f)
}
