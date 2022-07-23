# FLY

![](https://img.shields.io/badge/build-passing-brightgreen)
![](https://img.shields.io/badge/coverage-%2066.2%25-red)
![](https://img.shields.io/badge/license-MIT-blue)

一个简单易用的 Golang DB 辅助库, [详细文档](https://daodao97.github.io/fly)

- DataHook, 数据处理钩子, 方便的转换数据类型
- SqlBuilder, 无需手写SQL即可构造复杂查询条件
- hasOne/hasMany, 极其方便的获取`一对多`, `多对多`关联数据 
- validator, 轻松的进行数据合法性校验 
- extensible, 极易扩展功能方法 
- cacheEnable, 支持基于自定义存储的缓存管理

# usages

更多的使用样例请查看 [model_test.go](./model_test.go)

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