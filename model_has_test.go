package fly

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_otherKey(t *testing.T) {
	k, err := otherKeys([]string{"name", "  id as aid  ", "  ok AS OKS  "})
	assert.Equal(t, nil, err)
	assert.Equal(t, []string{"name", "aid", "OKS"}, k)

	k, err = otherKeys([]string{`

name as id

`, "  id as aid  ", "  ok AS OKS  "})
	assert.Equal(t, nil, err)
	assert.Equal(t, []string{"id", "aid", "OKS"}, k)
}
