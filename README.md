# FLY

![](https://img.shields.io/badge/build-passing-brightgreen)
![](https://img.shields.io/badge/coverage-%2066.2%25-red)
![](https://img.shields.io/badge/license-MIT-blue)

One easy-to-use go db library  [简体中文](./README_zh.md), [Document](https://daodao97.github.io/fly)

- data hook, Easy transition db data
- sql builder, Not need handwritten SQL
- hasOne/hasMany, Convenient get linked data
- validator, Flexible verification policies
- extensible, Easy extend custom hook/sql/validator
- cacheEnable, Support for custom cache implementations

# usages

more example please check out [model_test.go](./model_test.go)

```go
package main

import (
    "fmt"

    "github.com/daodao97/fly"
    _ "github.com/go-sql-driver/mysql"
)

func init() {
    err := fly.Init(map[string]*fly.Config{
        "default": {DSN: "root@tcp(127.0.0.1:3306)/fly_test?&parseTime=true"},
    })
    if err != nil {
        panic(err)
    }
}

func main() {
    m := fly.New(
        "user",
        fly.ColumnHook(fly.CommaInt("role_id"), fly.Json("profile")),
    )

    var err error

    _, err = m.Insert(map[string]interface{}{
        "name": "Seiya",
        "profile": map[string]interface{}{
            "hobby": "Pegasus Ryuseiken",
        },
    })

    var result []*User
    err = m.Select(fly.WhereGt("id", 1)).Binding(&result)

    var result1 *User
    err = m.SelectOne(fly.WhereEq("id", 1)).Binding(&result1)

    count, err := m.Count()
    fmt.Println("count", count)

    _, err = m.Update(User{
        ID:   1,
        Name: "星矢",
        Profile: &Profile{
            Hobby: "天马流行拳",
        },
        RoleIds: []int{2, 3},
    })

    _, err = m.Delete(fly.WhereEq("id", 1))

    fmt.Println(err)
}

type User struct {
    ID        int64    `db:"id"`
    Name      string   `db:"name"`
    Status    int64    `db:"status"`
    Profile   *Profile `db:"profile"`
    IsDeleted int      `db:"is_deleted"`
    RoleIds   []int    `db:"role_ids"`
    Score     int      `db:"score"`
}

type Profile struct {
    Hobby string `json:"hobby"`
}
```