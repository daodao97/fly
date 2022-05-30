package ggm

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestCache struct {
	cache map[string]string
}

func (t *TestCache) Get(key string) (string, error) {
	c, _ := t.cache[key]
	return c, nil
}

func (t *TestCache) Set(key string, data string) error {
	t.cache[key] = data
	return nil
}

func init() {
	cache = &TestCache{
		cache: map[string]string{},
	}
}

func TestModel_FindBy(t *testing.T) {
	m := New(User{})
	var u User
	err := m.FindBy(10, &u)
	assert.Equal(t, nil, err)
	err = m.FindBy(10, &u)
	assert.Equal(t, nil, err)
	fmt.Println("row", fmt.Sprintf("%+v", u))
}
