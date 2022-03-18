# ggm

golang generic orm, base on [sqlx](https://github.com/jmoiron/sqlx)

![](https://img.shields.io/badge/build-passing-brightgreen)
![](https://img.shields.io/badge/coverage-%2066.2%25-red)
![](https://img.shields.io/badge/license-MIT-blue)

## install

```shll
go get github.com/daodao97/ggm
```

## usage

Below is an example which shows some common use cases for sqlx. Check [model_test.go](./model_test.go) for more usage.

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/daodao97/ggm"
)

func init() {
	err := ggm.Init(map[string]*ggm.Config{
		"default": {
			DSN: "root@tcp(127.0.0.1:3306)/ggm_test?&parseTime=true",
		},
	})

	if err != nil {
		panic(err)
	}
}

func main() {
	m := ggm.New[*User]()

	user1 := &User{
		Name: "Seiya",
		Profile: ggm.NewJson(&Profile{
			Skill: "Pegasus Ryuseiken",
		}),
	}

	user2 := &User{
		Name: "Shun",
		Profile: ggm.NewJson(&Profile{
			Skill: "Nebula Chain",
		}),
	}

	_, err := m.Insert(user1, user2)
	if err != nil {
		log.Fatalln("Insert err", err)
		return
	}

	list, err := m.Select(ggm.WhereEq("id", 1))
	if err != nil {
		log.Fatalln("Select err", err)
		return
	}

	listJson, _ := json.Marshal(list)
	fmt.Printf("user list: %+v", string(listJson))

	user1.Id = 1
	user1.Name = "Seiya!"
	_, err = m.Update(user1)
	if err != nil {
		log.Fatalln("Update 1 err", err)
		return
	}

	user2.Name = "Shun!"
	_, err = m.Update(user2, ggm.WhereEq("name", "Shun"))
	if err != nil {
		log.Fatalln("Update2 err", err)
		return
	}

	_, err = m.Delete(ggm.WhereEq("id", 1))
	if err != nil {
		log.Fatalln("Delete err", err)
		return
	}
}

type User struct {
	Id      int                 `db:"id,pk"   json:"id"`
	Name    string              `db:"name"    json:"name"`
	Profile *ggm.Json[*Profile] `db:"profile" json:"profile"`
}

func (u User) Table() string {
	return "user"
}

func (u User) FakeDeleteKey() string {
	return "is_deleted"
}

type Profile struct {
	Skill string `json:"skill"`
}

```