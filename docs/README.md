![](https://img.shields.io/badge/build-passing-brightgreen)
![](https://img.shields.io/badge/coverage-%2066.2%25-red)
![](https://img.shields.io/badge/license-MIT-blue)

一个多简单易用的 Golang DB 辅助库

- DataHook, 数据处理钩子, 方便的转换数据类型
- SqlBuilder, 无需手写SQL即可构造复杂查询条件
- hasOne/hasMany, 极其方便的获取`一对多`, `多对多`关联数据
- validator, 轻松的进行数据合法性校验
- extensible, 极易扩展功能方法
- cacheEnable, 支持基于自定义存储的缓存管理

# 初始化

```go
fly.Init(map[string]*fly.Config{
    "default": {
        DSN: "root@tcp(127.0.0.1:3306)/fly_test?&parseTime=true",
        ReadDsn: "root@tcp(127.0.0.1:3306)/fly_test?&parseTime=true",
        Driver: "mysql",
        MaxOpenConn: 100,
        MaxIdleConn: 20,
    },
})
```

`map[key]*fly.Config` 
- `key` 为当前链接的名称, 模型默认使用链接为 `default`
- `DSN` 数据库链接, 必填
- `ReadDsn` 只读的数据库链接, 非必填, 若定义, 则查询操作将使用只读实例
- `Driver` 驱动模型, 默认 `mysql`
- `MaxOpenConn` 最大链接数, 默认值 100
- `MaxIdleConn` 最多可闲置的链接数, 默认值 20

# 模型定义

```go
m := fly.New("table_name")
```

## WithConn

```go
m := fly.New(
    "table_name",
    fly.WithConn("conn_name") 
)
```
定义所使用的数量实例的名称, 默认为 `default`

## WithDB
```go
m := fly.New(
    "table_name",
    fly.WithDB(*sql.DB) 
)
```
当所使用的数据库实例不是通过 `Init` 初始化的情况下, 定义自定义DB实例

## WithPrimaryKey
```go
m := fly.New(
    "table_name",
    fly.WithPrimaryKey("id") 
)
```
定义当前模型的主键, 默认为 `id`

## WithFakeDelKey
```go
m := fly.New(
    "table_name",
    fly.WithFakeDelKey("is_deleted") 
)
```
定义当前模型的伪删除字段, 如果配置的伪删除指点, 当执行 `Delete` 操作时, 会自定转为 `set is_deleted = 1` 的 `Update` 操作

## ColumnHook
```go
m := fly.New(
    "table_name",
    fly.ColumnHook(fly.Json("json_str_column")) 
)
```
定义数据库字段的数据转换方式, 比如 mysql 中存储的为 `json 格式字符串`, 代码中的数据模型为 `map` or `struct` 
此时可以使用 `fly.Json` 钩子, 自动处理该字段在 `db` 和 `程序` 之间的相互转换.

更多Hook, 可查看 [数据转换](/?id=数据转换) 章节

## ColumnValidator
```go
m := fly.New(
    "table_name",
    fly.ColumnValidator(
        fly.Validate(
            "name",
            Required(),
            Unique(WithMsg("名称已经存在, 请换一个吧")),
        ),
        fly.Validate("profile", IfRequired("name")),
    )
)
```
如上就定义了
- `name` 字段是必须的 且 是唯一的
- `profile` 字段 在 `name` 存在的情况是必须的

的约束, 在执行 `Insert` 和 `Update` 操作时将进行如上的合法性校验, 
更多的校验规则, 可查看 [数据校验](/?id=数据校验) 章节

## HasOne
```go
m := fly.New(
    "table_name",
    fly.HasOne(HasOpts{Table: "level", OtherKeys: []string{"name as LevelName"}}),
)
```
定义当前表的数据与其他表`一对一`的关联关系, 定义后, 查询模型中将会自动查出对应的关联数据

HasOpts
- `Conn` 关联表所在的db实例, 也就是存在逻辑对应关系的数据可以跨实例, 默认 `default`
- `Database` 关联表所在的库, 默认为空, 也就是当前表所在库 
- `Table` 关联表的表名, 必须填写 
- `LocalKey` 当前表的关联数据字段, 类似`join on a.id = b.id` 中的 `a.id` 
- `ForeignKey` 关联表中的关联数据字段, 类似`join on a.id = b.id` 中的 `b.id`
- `OtherKeys` 关联表中需要连带查询出的其他字段

## HasMany
```go
m := fly.New(
    "table_name",
    fly.HasMany(HasOpts{Table: "followed", ForeignKey: "uid", OtherKeys: []string{"follow_to_uid as FollowerId"}}),
)
```
定义当前表的数据与其他表`一对多`的关联关系, 定义后, 查询模型中将会自动查出对应的关联数据

HasOpts 的定义与 `hasOne` 相同

# 数据操作
## 写入
```go
lastId, err := m.Insert(data)
```

`data` 的数据类型可以是 
- `map[string]interface`
- `*map[string]interface` 
- `struct` 
- `*struct`

中的任意一种

## 更新
```go
ok, err := m.Update(data, opt...)
```

`data` 的数据类型支持范围与 `Insert` 相同

若 `data` 中存在主键, 则会自动追加 `WhereEq(pk, value)` 到 `opt` 的约束条件当中.

## 查询

```go

```


## 删除

# 数据转换

# 数据验证

# 关联数据

# 自定义日志

# 自定义缓存
