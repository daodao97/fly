package fly

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/pkg/errors"
)

type Config struct {
	DSN         string
	ReadDsn     string
	Driver      string
	MaxOpenConn int
	MaxIdleConn int
}

var pool = sync.Map{}

func Init(conns map[string]*Config) error {
	for conn, conf := range conns {
		db, err := newDb(conf)
		if err != nil {
			return err
		}
		pool.Store(conn, db)
		if conf.ReadDsn != "" {
			rdb, err := newDb(&Config{
				DSN:         conf.ReadDsn,
				Driver:      conf.Driver,
				MaxOpenConn: conf.MaxOpenConn,
				MaxIdleConn: conf.MaxIdleConn,
			})
			if err != nil {
				return err
			}
			pool.Store(readConn(conn), rdb)
		}
	}
	return nil
}

func Close() {
	pool.Range(func(key, value interface{}) bool {
		db := value.(*sql.DB)
		_ = db.Close()
		return true
	})
}

func readConn(conn string) string {
	return conn + "_read"
}

func db(conn string) (*sql.DB, error) {
	if _db, ok := pool.Load(conn); ok {
		return _db.(*sql.DB), nil
	}
	return nil, errors.New("connection not found : " + conn)
}

func DB(conn string) (*sql.DB, error) {
	return db(conn)
}

func newDb(conf *Config) (*sql.DB, error) {
	driver := conf.Driver
	if driver == "" {
		driver = "mysql"
	}

	db, err := sql.Open(driver, conf.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed Connection database: %s", err)
	}

	// 设置数据库连接池最大连接数
	MaxOpen := 100
	if conf.MaxOpenConn != 0 {
		MaxOpen = conf.MaxOpenConn
	}
	db.SetMaxOpenConns(MaxOpen)

	// 连接池最大允许的空闲连接数
	// 如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭
	MaxIdle := 20
	if conf.MaxIdleConn != 0 {
		MaxIdle = conf.MaxIdleConn
	}
	db.SetMaxIdleConns(MaxIdle)
	return db, nil
}
