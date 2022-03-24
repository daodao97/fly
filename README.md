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

Below is an example which shows some common use cases for ggm. Check [model_test.go](./model_test.go) for more usage.

### init db

We can initialize some db resources commonly used by programs, like this

```go
// map[conn_name]db_config
ggm.Init(map[string]*ggm.Config{
    "default": {
        DSN: "root@tcp(127.0.0.1:3306)/ggm_test?&parseTime=true",
    },
})
```

Of course, we can also instantiate some temporary DB resources, like this

```go
m := ggm.NewConn(&ggm.Config{
	DSN: "root@tcp(127.0.0.1:3306)/ggm_test?&parseTime=true" 
})
```

### data model

For example, we have a table with the following structure

```sql
CREATE TABLE `user` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `status` tinyint(4) NOT NULL DEFAULT '0',
  `profile` varchar(200) NOT NULL,
  `is_deleted` tinyint(3) unsigned NOT NULL DEFAULT '0',
  `ctime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `mtime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
```

#### struct define
The structure model for this table is

```go
type User struct {
	Id      int             `db:"id,pk"   json:"id"`
	Name    string          `db:"name"    json:"name"`
	Profile string          `db:"profile" json:"profile"`
	CTime   time.Time       `db:"ctime"   json:"ctime"`
}

func (u User) Table() string {
	return "user"
}

```

`interface: Table() string` the struct must implement.

struct field must have `db` tag, value is db field name.

#### set conn

If you are using a db resource that is not the `default`

```go
func (u User) Conn() string {
	return "conn_name"
}
```

#### fake delete

If you have a field to mark fake delete

```go
func (u User) FakeDeleteKey() string {
	return "is_deleted"
}
```

then, the `delete sql` will converted to `update ${fakeDeleteKey} = 1` when we delete data

`select sql` will auto add `${fackDeleteKey} = 0`, to filter deleted data.

### select

```go
m := ggm.New[*User]() // or ggm.New[User]() 

m.Select(ggm.WhereEq("id", 1))
```

detail of `where condition` see [where condition](https://github.com/daodao97/ggm#where-condition)

### insert

```go
user := &User{Name: "Seiya"}

// single insert
m.Insert(user)

// batch insert
m.Insert(user, user2, ...)
```

### update

#### use primary key update

```go
user := &User{Id: 1, Name: "Seiya!!!"}
m.Update(user)
```

#### with where condition

```go
user := &User{Name: "Seiya!!!"}
m.Update(user, ggm.WhereEq("id", 1))
```

### delete

```go
m.Delete(ggm.WhereEq("id", 1))
```

### where condition

```go
m.Select(
    WhereEq("id", 1),
    WhereGt("age", 20),
    WhereLike("name", "dd"),
    WhereGroup(
        WhereEq("sex", 1),
        WhereOrEq("class", 2),
        WhereGroup(
            WhereEq("sex1", 1),
            WhereEq("class2", 2),
        ),
    ),
    OrderBy("id", DESC),
)
```

more example, checkout [sql_test.go](/sql_test.go)

### data type

#### Json

If the value of field `user.profile` is `json_string` like `{"skill":"Pegasus Ryuseiken"}`

```go
type User struct {
	Id      int                 `db:"id,pk"   json:"id"`
	Name    string              `db:"name"    json:"name"`
	Profile *ggm.Json[*Profile] `db:"profile" json:"profile"`
	CTime   time.Time           `db:"ctime"   json:"ctime"`
}

type Profile struct {
    Skill string `json:"skill"`
}
```

`Profile{Skill: "xxx"}`  <==> '{"skill":"xxx"}'

Data can be automatically converted into struct for use by programs.

#### Time

```go
type User struct {
	Id      int                 `db:"id,pk"   json:"id"`
	Name    string              `db:"name"    json:"name"`
	Profile *ggm.Json[*Profile] `db:"profile" json:"profile"`
	CTime   ggm.Time            `db:"ctime"   json:"ctime"`
}
```

when api response or json.Marshal

ctime : `2022-03-19T11:52:19Z` => `2022-03-19 11:52:19`

#### define yourself data type

Implement the following interfaces

```go
type DataType[T any] interface {
	Value() (driver.Value, error)
	Scan(value any) error
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(b []byte) error
	Get() T
}
```

### Linked data

#### hasOne

`one to one`

```go
type User struct {
	Id      int             `db:"id,pk"   json:"id"`
	Name    string          `db:"name"    json:"name"`
	Profile *Json[*Profile] `db:"profile" json:"profile"`
	Score   int             `db:"score"   json:"score" hasOne:"user_score:uid"`
	Score2  int             `db:"score2"  json:"score2" hasOne:"user_score:uid"`
}
```

#### hasMany

`one to N`

```go
type User struct {
	Id      int             `db:"id,pk"   json:"id"`
	Name    string          `db:"name"    json:"name"`
	Profile *Json[*Profile] `db:"profile" json:"profile"`
	Logs    []*Log          `json:"logs"  hasMany:"user_log:uid"`
}

type Log struct {
	Message string `db:"message" json:"message"`
}
```

hasOne or hasMany tag token:

`[conn.][database.]table:[local_key->]foreign_key`

`m.Select()` will auto query the linked data into the struct.

Check [model_test.go](./model_test.go) for detail.