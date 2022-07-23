package mysql

import (
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"

	"github.com/daodao97/fly"
	"github.com/daodao97/fly/interval/util"
)

var dsn = "root@tcp(127.0.0.1:3306)/fly_test?&parseTime=true"

var m fly.Model

func init() {
	err := fly.Init(map[string]*fly.Config{
		"default": {DSN: dsn},
	})
	if err != nil {
		panic(err)
	}

	m = fly.New(
		"user",
		fly.ColumnHook(fly.Json("profile"), fly.CommaInt("role_ids")),
		fly.ColumnValidator(
			fly.Validate("name", fly.Required()),
		),
		fly.HasOne(fly.HasOpts{Table: "level", OtherKeys: []string{"name as LevelName"}}),
		fly.HasMany(fly.HasOpts{Table: "followed", ForeignKey: "uid", OtherKeys: []string{"follow_to_uid as FollowerId"}}),
		fly.WithFakeDelKey("is_deleted"),
	)

	_, err = m.Exec(`
CREATE TABLE IF NOT EXISTS user (
  id int(11) unsigned NOT NULL AUTO_INCREMENT,
  name varchar(50) NOT NULL,
  status tinyint(4) NOT NULL DEFAULT '0',
  profile varchar(200) NOT NULL,
  is_deleted tinyint(3) unsigned NOT NULL DEFAULT '0',
  role_ids varchar(255) NOT NULL DEFAULT '',
  ctime datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  mtime datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
`)
	if err != nil {
		panic(err)
	}

	_, err = m.Exec("TRUNCATE TABLE `user`;")
	if err != nil {
		panic(err)
	}
}

type User struct {
	ID         int64     `db:"id"`
	Name       string    `db:"name"`
	Status     int64     `db:"status"`
	Profile    *Profile  `db:"profile"`
	RoleIds    []int     `db:"role_ids"`
	CTime      time.Time `db:"ctime"`
	LevelId    int64     `db:"level_id"`
	LevelName  string
	FollowerId []int64
}

type Profile struct {
	Hobby string `json:"hobby"`
}

var userMapMod = map[string]interface{}{
	"name": "Seiya",
	"profile": map[string]interface{}{
		"hobby": "Pegasus Ryuseiken",
	},
	"role_ids": []int{1, 2},
	"level_id": 1,
}

var userMapModPtr = &map[string]interface{}{
	"name": "Seiya",
	"profile": map[string]interface{}{
		"hobby": "Pegasus Ryuseiken",
	},
	"role_ids": []int{1, 2},
	"level_id": 1,
}

var userStructMod = User{
	Name: "Seiya",
	Profile: &Profile{
		Hobby: "Pegasus Ryuseiken",
	},
	RoleIds: []int{1, 2},
	LevelId: 1,
}

var userStructModPtr = &User{
	Name: "Seiya",
	Profile: &Profile{
		Hobby: "Pegasus Ryuseiken",
	},
	RoleIds: []int{1, 2},
	LevelId: 1,
}

func Test_Insert(t *testing.T) {
	_, err := m.Insert(userMapMod)
	assert.Equal(t, nil, err)

	_, err = m.Insert(userMapModPtr)
	assert.Equal(t, nil, err)

	_, err = m.Insert(userStructMod)
	assert.Equal(t, nil, err)

	_, err = m.Insert(userStructModPtr)
	assert.Equal(t, nil, err)

	_, err = m.Insert(1)
	assert.Equal(t, util.ErrParamsType, err)
}

func Test_Select(t *testing.T) {
	var result []*User
	err := m.Select(fly.Field("id", "name", "profile", "role_ids", "level_id"), fly.WhereEq("id", 1)).Binding(&result)
	assert.Equal(t, nil, err)
	userStructModPtr.ID = 1
	userStructModPtr.LevelName = "黄金"
	userStructModPtr.FollowerId = []int64{2, 3}
	assert.Equal(t, []*User{userStructModPtr}, result)

	var errType int
	err = m.Select(fly.WhereGe("id", 1)).Binding(&errType)
	assert.Equal(t, fly.ErrRowsBindingType, err)
}

func Test_SelectOne(t *testing.T) {
	var result *User
	err := m.SelectOne(fly.Field("id", "name", "profile", "role_ids", "level_id"), fly.WhereEq("id", 1)).Binding(&result)
	assert.Equal(t, nil, err)
	userStructModPtr.ID = 1
	userStructModPtr.LevelName = "黄金"
	userStructModPtr.FollowerId = []int64{2, 3}
	assert.Equal(t, userStructModPtr, result)
}

func Test_Update(t *testing.T) {
	result, err := m.Update(User{
		ID:   1,
		Name: "星矢",
		Profile: &Profile{
			Hobby: "天马流行拳",
		},
		RoleIds: []int{2, 3},
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, true, result)
}

func Test_Count(t *testing.T) {
	count, err := m.Count()
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(4), count)
}

func Test_Delete(t *testing.T) {
	_, err := m.Delete(fly.WhereEq("id", 1))
	assert.Equal(t, nil, err)
}
