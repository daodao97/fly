![](https://img.shields.io/badge/build-passing-brightgreen)
![](https://img.shields.io/badge/coverage-%2066.2%25-red)
![](https://img.shields.io/badge/license-MIT-blue)

一个简单易用的 Golang DB 辅助库

- DataHook, 数据处理钩子, 方便的转换数据类型
- SqlBuilder, 无需手写SQL即可构造复杂查询条件
- hasOne/hasMany, 极其方便的获取`一对多`, `多对多`关联数据
- validator, 轻松的进行数据合法性校验
- extensible, 极易扩展功能方法
- cacheEnable, 支持基于自定义存储的缓存管理

# 安装

```go
go get github.com/daodao97/fly
```

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
定义所使用实例的名称, 默认为 `default`

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
var result []*User
err := m.Select(
        fly.Field("id", "name", "profile", "role_ids", "level_id"),
        fly.WhereEq("id", 1)
    ).Binding(&result)
```
查询条件由一系列 `Option` 组成, 这些 `Option` 用来描述查询的字段, 条件, 分页等, 最终自动构造成完整的SQL执行.

更多的查询条件构造, 请查看 [sql_test.go](https://github.com/daodao97/fly/blob/master/sql_test.go)

> tips: 目前仅支持单表查询

> 若模型定义了`伪删除`特性, 查询条件中将自动追加 `WhereEq({is_deleted}, 0)` 条件, 以过滤已删除数据
## 删除
```go
_, err := m.Delete(fly.WhereEq("id", 1))
```
`Delete` 方法的入参跟 `Select` 相同, 有一系列 `Options` 组成.

> 若模型定义了`伪删除`特性,  `Delete` 操作将转换为 `Update set {is_deleted} = 1` 的更新操作

# 数据转换
工作中我们经常遇到这样的场景, 比如在db中某个字段存储的格式是字符串, 但是逻辑意义上一个JSON格式的字符串, 那么我们读取该字段时,
需要把db中的字符串转换为一个 map / struct 对象, 在保存该字段是, 需要将对象转换成JSON字符串, 等等, 类似的转换需求很多,
因此我们抽象了 `ColumnHook` 来快捷的处理这种需求.

```go
type HookData interface {
	Input(row map[string]interface{}, fieldValue interface{}) (interface{}, error)
	Output(row map[string]interface{}, fieldValue interface{}) (interface{}, error)
}
```

任意实现以上接口的都可以作为 `ColumnHook` 使用, 系统内置的`ColumnHook`有,

## Json
```go
{"a":1} <=> map[string]interface{}{"a":1}

[{"a":1}] <=> []interface{}{map[string]interface{}{"a":1}}

m = fly.New(
    "user",
    fly.ColumnHook(fly.Json("profile")),
)
```
如此定义, 便可以在 `Create`, `Update`, `Select` 时自动处理数据类型.

## Object

如果确定数据是 `{k:v}` 型的数据可以直接使用 `fly.Object('xxx')`

## Array

如果确定数据是 `[*]` 型的数据可以直接使用 `fly.Array('xxx')`

## CommaInt

处理 db中 `"1,2,3"` 型数据, 在程序中需要转换成 `[]int`

## CommaStr

处理 db中 `"a,b,c"` 型数据, 在程序中需要转换成 `[]string`

如果有更多的数据类型需求, 可以基于 `HookData` 接口实现自定义hook, 在模型中使用.

# 数据验证

```go
m = fly.New(
    "user",
    fly.ColumnValidator(
        fly.Validate("name", fly.Required(), fly.Unique()),
    ),
)
```
如上便定义了 `name` 字符为必填, 并且唯一 两个约束, 默认支持的校验条件有:

## Required

字段必须存在, 且非零值

## IfRequired

当 `ifField` 字段存在是, 那么, 当前字段是必须的

## Unique

当前字段的值在该表中必须是唯一的

在系统默认的规则不能满足需求是, 也可以自定义验证规则, 具体如下

```go
type Valid = func(v *ValidInfo) error

func CustomeValidate(v *ValidInfo) err {
    return nil
}

fly.Validate("name", CustomeValidate)
```

# 关联数据

## hasOne 一对一

[hasOne](/?id=hasone)

## hasMany 一对多

[hasMany](/?id=hasmany)

# 自定义日志

```go
type Logger interface {
	Log(level Level, keyValues ...interface{}) error
}
```

任何实现以上接口的对象都可以最为内部日志记录的主体, 默认的日志记录到程序的标准输出

当然也可以通过 `fly.SetLogger` 设置自定义的日志记录器.

# 自定义缓存

```go
type Cache interface {
	Get(key string) (string, error)
	Del(key string) error
	Set(key string, data string) error
}
```

实现以上接口的存储对象可以用于查询模型中的缓存控制, 具体的缓存选型可以根据自己的具体情况而定,
样例可以查看 [model_cache_test.go](https://github.com/daodao97/fly/blob/master/model_cache_test.go)

通过 `fly.SetCache` 设置缓存模型后, 那么可以通过 `FindBy`, `UpdateBy`方法查询和更新缓存, 
其他方法不支持缓存.
