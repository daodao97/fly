package ggm

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"sync"
)

var pool = sync.Map{}

func Init(conns map[string]*Config) error {
	for conn, conf := range conns {
		db, err := newDb(conf)
		if err != nil {
			return err
		}
		pool.Store(conn, db)
	}
	return nil
}

func Close() {
	pool.Range(func(key, value interface{}) bool {
		db := value.(*sqlx.DB)
		_ = db.Close()
		return true
	})
}

func xdb(conn string) (*sqlx.DB, bool) {
	if _db, ok := pool.Load(conn); ok {
		return _db.(*sqlx.DB), ok
	}
	return nil, false
}

type Config struct {
	DSN         string
	Driver      string
	MaxOpenConn int
	MaxIdleConn int
}

// newDb new db and retry connection when has error.
func newDb(c *Config) (db *sqlx.DB, err error) {
	db, err = makeDb(c)
	if err != nil {
		return nil, errors.Wrap(err, "open mysql error")
	}
	return
}

// 生成 原生 xdb 对象
func makeDb(conf *Config) (*sqlx.DB, error) {
	driver := conf.Driver
	if driver == "" {
		driver = "mysql"
	}

	db, err := sqlx.Open(driver, conf.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed Connection database: %s", err)
	}

	// 设置数据库连接池最大连接数
	MaxOpenConns := 100
	if conf.MaxOpenConn != 0 {
		MaxOpenConns = conf.MaxOpenConn
	}
	db.SetMaxOpenConns(MaxOpenConns)

	// 连接池最大允许的空闲连接数
	// 如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭
	MaxIdleConns := 20
	if conf.MaxIdleConn != 0 {
		MaxIdleConns = conf.MaxIdleConn
	}
	db.SetMaxIdleConns(MaxIdleConns)
	return db, nil
}
