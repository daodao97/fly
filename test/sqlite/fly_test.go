package sqlite

import (
	"github.com/daodao97/fly/interval/util"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"

	"github.com/daodao97/fly"
)

var dsn = "./fly_sqlite.test"

var m fly.Model

func init() {
	_ = os.Remove(dsn)

	err := fly.Init(map[string]*fly.Config{
		"default": {DSN: dsn, Driver: "sqlite3"},
	})
	if err != nil {
		panic(err)
	}

	m = fly.New("user", fly.ColumnHook(fly.Json("profile"), fly.CommaInt("role_ids")), fly.WithFakeDelKey("is_deleted"))

	_, err = m.Exec(`
CREATE TABLE "user" (
    id     integer not null constraint table_name_pk primary key autoincrement,
    name   varchar  default '' not null,
	status integer default 0 not null,
	profile text not null,
	is_deleted integer default 0 not null,
	role_ids text default '' not null,
    ctime  datetime default CURRENT_TIMESTAMP not null,
    mtime  datetime default CURRENT_TIMESTAMP not null
);
`)
	if err != nil {
		panic(err)
	}
}

type User struct {
	ID        int64    `json:"id"`
	Name      string   `json:"name"`
	Status    int64    `json:"status"`
	Profile   *Profile `json:"profile"`
	IsDeleted int      `json:"is_deleted"`
	RoleIds   []int    `json:"role_ids"`
	Score     int      `json:"score"`
}

type Profile struct {
	Hobby string `json:"hobby"`
}

func Test_Insert(t *testing.T) {
	result, err := m.Insert(map[string]interface{}{
		"name": "Seiya",
		"profile": map[string]interface{}{
			"hobby": "Pegasus Ryuseiken",
		},
		"role_ids": []int{1, 2},
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(1), result)

	_, err = m.Insert(1)
	assert.Equal(t, util.ErrParamsType, err)
}

func Test_Select(t *testing.T) {
	var result []*User
	err := m.Select(fly.WhereGe("id", 1)).Binding(&result)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, len(result) > 0)

	var errType int
	err = m.Select(fly.WhereGe("id", 1)).Binding(&errType)
	assert.Equal(t, fly.ErrRowsBindingType, err)
}

func Test_SelectOne(t *testing.T) {
	var result *User
	err := m.SelectOne(fly.WhereEq("id", 1)).Binding(&result)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, result)
	assert.Equal(t, "Seiya", result.Name)
	assert.Equal(t, []int{1, 2}, result.RoleIds)
}

func Test_Update(t *testing.T) {
	result, err := m.Update(map[string]interface{}{
		"id":   1,
		"name": "星矢",
		"profile": map[string]interface{}{
			"hobby": "天马流行拳",
		},
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, true, result)
}

func Test_Count(t *testing.T) {
	count, err := m.Count()
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(1), count)
}

func Test_Delete(t *testing.T) {
	_, err := m.Delete(fly.WhereEq("id", 1))
	assert.Equal(t, nil, err)

	_ = os.Remove(dsn)
}
