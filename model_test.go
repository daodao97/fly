package ggm

/*

CREATE TABLE `user` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `status` tinyint(4) NOT NULL DEFAULT '0',
  `profile` varchar(200) NOT NULL, -- json string
  `is_deleted` tinyint(3) unsigned NOT NULL DEFAULT '0',
  `ctime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `mtime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `user_score` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uid` int(11) unsigned NOT NULL, -- user id
  `score` int(11) unsigned NOT NULL DEFAULT '0',
  `score2` int(11) unsigned DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `user_log` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uid` int(11) unsigned NOT NULL DEFAULT '0', -- user id
  `message` varchar(50) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

*/

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type User struct {
	Id      int              `db:"id,pk"   json:"id"`
	Name    string           `db:"name"    json:"name"`
	Profile *Json[*Profile]  `db:"profile" json:"profile"`
	RoleIds *CommaSlice[int] `db:"role_ids" json:"role_ids"`
	CTime   Time             `db:"ctime"   json:"ctime"`
	Score   int              `db:"score"   json:"score" hasOne:"user_score:uid"`
	Score2  int              `db:"score2"  json:"score2" hasOne:"user_score:uid"`
	Logs    []*Log           `json:"logs"  hasMany:"user_log:uid"`
}

func (u User) Table() string {
	return "user"
}

func (u User) FakeDeleteKey() string {
	return "is_deleted"
}

type Profile struct {
	Hobby string `json:"hobby"`
}

type Log struct {
	Message string `db:"message" json:"message"`
}

var dsn = "root@tcp(127.0.0.1:3306)/ggm_test?&parseTime=true"

func init() {
	err := Init(map[string]*Config{
		"default": {DSN: dsn},
	})
	if err != nil {
		panic(err)
	}
}

func mysqlConf() *Config {
	return &Config{
		DSN: dsn,
	}
}

func TestModel_InsertSingle(t *testing.T) {

	// T is ptr struct
	m := New[*User]()

	p := NewJson(&Profile{
		Hobby: "will_ok",
	})

	user := &User{
		Name:    "ok",
		Profile: p,
	}
	_, err := m.Insert(user)
	assert.Equal(t, nil, err)

	// T is struct
	m1 := New[User]()

	p1 := NewJson(&Profile{
		Hobby: "will_ok_1",
	})

	user1 := User{
		Name:    "ok_1",
		Profile: p1,
	}
	_, err1 := m1.Insert(user1)
	assert.Equal(t, nil, err1)
}

func TestModel_InsertMulti(t *testing.T) {
	m := New[*User]()

	users := []*User{
		{
			Name: "ok",
			Profile: NewJson(&Profile{
				Hobby: "will_ok",
			}),
		},
		{
			Name: "ok2",
			Profile: NewJson(&Profile{
				Hobby: "will_ok2",
			}),
		},
	}
	_, err := m.Insert(users...)
	assert.Equal(t, nil, err)
}

func TestModel_Select(t *testing.T) {
	m := New[*User]()
	list, err := m.Select(Where("id", "<", 5))
	assert.Equal(t, nil, err)
	jsonStr, err := json.Marshal(list)
	assert.Equal(t, nil, err)
	fmt.Println("list: ", string(jsonStr))
}

func TestModel_Update(t *testing.T) {
	m := New[*User]()

	p := NewJson(&Profile{
		Hobby: "you",
	})

	user := &User{
		Id:      1,
		Name:    "ok_man",
		Profile: p,
	}
	_, err := m.Update(user)
	assert.Equal(t, nil, err)
}

func TestModel_Delete(t *testing.T) {
	m := New[*User]()

	user := &User{
		Id: 1,
	}
	_, err := m.Delete(WhereEq("id", user.Id))
	assert.Equal(t, nil, err)
}
