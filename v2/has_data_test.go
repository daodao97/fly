package ggm

import (
	"fmt"
	"testing"
)

func Test_explodeHasStr(t *testing.T) {
	token := []string{
		"pool.db.table:id->lid,a",
		"pool.db.table:id->lid,a as a1",
		"pool.db.table:id->lid,a,b,c",
		"db.table:id->lid,a as a1",
		"db.table:id->lid,a,b,c",
		"table:lid,",
	}

	for _, v := range token {
		fmt.Println(explodeHasStr(v))
	}
}
