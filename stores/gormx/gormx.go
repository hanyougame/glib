package gormx

import (
	"github.com/hanyougame/glib/stores/gormx/config"
	"github.com/hanyougame/glib/stores/gormx/database"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Engine *gorm.DB

// Must 创建全局数据库连接
func Must(c config.Config) {
	Engine = NewEngine(c)
}

// NewEngine 创建数据库连接
func NewEngine(c config.Config) *gorm.DB {
	switch c.Mode {
	case config.Mysql:
		return database.NewEngine(c, mysql.Open(c.DSN))
	case config.Postgres:
		return database.NewEngine(c, postgres.Open(c.DSN))
	case config.ClickHouse:
		return database.NewEngine(c, clickhouse.Open(c.DSN))
	}
	panic("unsupported database mode")
}
