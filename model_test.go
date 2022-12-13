package fly

import (
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

var dsn = "root@tcp(127.0.0.1:3306)/fly_test?&parseTime=true"

var m Model

func init() {
	err := Init(map[string]*Config{
		"default": {DSN: dsn},
	})
	if err != nil {
		panic(err)
	}

	m = New(
		"user",
		ColumnHook(CommaInt("role_ids"), Json("profile")),
		ColumnValidator(
			Validate(
				"name",
				func(v *ValidInfo) error {
					_ = logger.Log(LevelDebug, v.Field+"字段的一个自定义校验")
					return nil
				},
				Required(),
				Unique(WithMsg("名称已经存在, 请换一个吧")),
			),
			Validate("profile", IfRequired("name")),
		),
		HasOne(HasOpts{Table: "level", OtherKeys: []string{"name as LevelName"}}),
		HasMany(HasOpts{Table: "followed", ForeignKey: "uid", OtherKeys: []string{"follow_to_uid as FollowerId"}}),
		WithFakeDelKey("is_deleted"),
	)
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
	Hobby string `db:"hobby"`
}

func Test_Insert(t *testing.T) {
	u := map[string]interface{}{
		"name": "Seiya",
		"profile": map[string]interface{}{
			"hobby": "Pegasus Ryuseiken",
		},
		"role_ids": []int{1, 2},
		"level_id": 1,
	}
	result, err := m.Insert(u)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, result > 0)

	_, err = m.Insert(u)
	assert.NotEqual(t, nil, err)
}

func Test_Select(t *testing.T) {
	var result []*User
	err := m.Select(WhereGe("id", 1)).Binding(&result)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, len(result) > 0)

	var errType int
	err = m.Select(WhereGe("id", 1)).Binding(&errType)
	assert.Equal(t, ErrRowsBindingType, err)
}

func Test_SelectOne(t *testing.T) {
	var result *User
	row := m.SelectOne(WhereEq("id", 1))
	err := row.Binding(&result)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, result)

	//tmp, _ := json.Marshal(result)
	//spew.Dump(string(tmp))

	spew.Dump(row.Data, result)
}

func Test_Update(t *testing.T) {
	//result, err := m.Update(User{
	//	ID:   1,
	//	Name: "星矢",
	//	Profile: &Profile{
	//		Hobby: "天马流行拳",
	//	},
	//	RoleIds: []int{2, 3},
	//	CTime:   time.Now().UTC(),
	//})
	//assert.Equal(t, nil, err)
	//assert.Equal(t, true, result)

	result, err := m.Update(map[string]interface{}{
		"id":   1,
		"name": "星矢1",
		"profile": &Profile{
			Hobby: "天马流行拳1",
		},
		"role_ids": []int{1, 2, 3},
		"ctime":    time.Now().UTC(),
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, true, result)

	result, err = m.Update(map[string]interface{}{
		"id":   1,
		"name": "星矢2",
		"profile": map[string]interface{}{
			"hobby": "天马流行拳2",
		},
		"role_ids": []int{1, 2, 3, 4},
		"ctime":    time.Now().UTC(),
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, true, result)
}

func Test_Delete(t *testing.T) {
	_, err := m.Delete(WhereEq("id", 1))
	assert.Equal(t, nil, err)
}

func Test_Count(t *testing.T) {
	count, err := m.Count()
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(0), count)
}
