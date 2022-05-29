package ggm

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var dsn = "root@tcp(127.0.0.1:3306)/ggm_test?&parseTime=true"

func init() {
	err := Init(map[string]*Config{
		"default": {DSN: dsn},
	})
	if err != nil {
		panic(err)
	}
}

type User struct {
	ID        int64     `db:"id,pk"`
	Name      string    `db:"name"`
	Status    int64     `db:"status"`
	Profile   string    `db:"profile"`
	IsDeleted int64     `db:"is_deleted"`
	RoleIds   string    `db:"role_ids"`
	Score     int       `db:"score"   json:"score" hasOne:"user_score:uid"`
	Logs      []Log     `json:"logs"  hasMany:"user_log:uid"`
	Ctime     time.Time `db:"ctime,ii"`
	Mtime     time.Time `db:"mtime,ii"`
}

func (u User) Table() string {
	return "user"
}

type Log struct {
	Message string `db:"message" json:"message"`
}

func TestModel_Select(t *testing.T) {
	m := New(User{})
	var result []User
	err := m.Select(&result, WhereGt("id", 9))
	assert.Equal(t, nil, err)
	fmt.Println(fmt.Sprintf("%+v", result))
}

func TestModel_SelectOne(t *testing.T) {
	m := New(User{})
	result := &User{}
	err := m.SelectOne(result, WhereEq("id", 1))
	assert.Equal(t, nil, err)
	fmt.Println(fmt.Sprintf("%+v", result))
}

func TestModel_Count(t *testing.T) {
	m := New(User{})
	count, err := m.Count(WhereGt("id", 1))
	assert.Equal(t, nil, err)
	fmt.Println(fmt.Sprintf("%+v", count))
}

func TestModel_Insert(t *testing.T) {
	m := New(User{})
	count, err := m.Insert(User{Name: "okkk"})
	assert.Equal(t, nil, err)
	fmt.Println(fmt.Sprintf("%+v", count))
}

func TestModel_Update(t *testing.T) {
	m := New(User{})
	count, err := m.Update(User{Name: "okiss", ID: 11})
	assert.Equal(t, nil, err)
	fmt.Println(fmt.Sprintf("%+v", count))
}

func TestModel_Delete(t *testing.T) {
	m := New(User{})
	count, err := m.Delete(WhereEq("id", 11))
	assert.Equal(t, nil, err)
	fmt.Println(fmt.Sprintf("%+v", count))
}
