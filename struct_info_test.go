package ggm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_reflectInfo(t *testing.T) {
	info, err := structInfo[User]()
	assert.Equal(t, nil, err)
	JsonDump(info)
}
