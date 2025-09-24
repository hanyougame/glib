package database

import (
	"fmt"
	"time"

	"github.com/hanyougame/glib/stores/gormx/config"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

func NewEngine(c config.Config, dialector gorm.Dialector, opt ...gorm.Option) (*gorm.DB, error) {
	cfg := &gorm.Config{
		PrepareStmt:            c.PrepareStmt,
		SkipDefaultTransaction: c.SkipDefaultTransaction,
		Logger:                 New(),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}

	engine, err := gorm.Open(dialector, append(opt, cfg)...)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	if c.Debug {
		engine = engine.Debug()
	}

	if c.Trace {
		registerTraceHook(engine)
	}

	// 设置连接池参数
	sqlDB, err := engine.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %v", err)
	}

	if c.MaxIdleConn > 0 {
		sqlDB.SetMaxIdleConns(c.MaxIdleConn)
	}
	if c.MaxOpenConn > 0 {
		sqlDB.SetMaxOpenConns(c.MaxOpenConn)
	}
	if c.MaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(c.MaxLifetime) * time.Minute)
	}

	// 读写分离
	if c.Separation {
		engine, err = newRWDBEngine(c, engine)
		if err != nil {
			return nil, fmt.Errorf("failed to setup read-write separation for database: %v", err)
		}
	}

	return engine, nil
}

// 读写分离
func newRWDBEngine(c config.Config, engine *gorm.DB) (*gorm.DB, error) {
	var (
		replicas, sources []gorm.Dialector
		dialectorFactory  func(string) gorm.Dialector
	)

	switch c.Mode {
	case config.Mysql:
		dialectorFactory = mysql.Open
	case config.Postgres:
		dialectorFactory = postgres.Open
	case config.ClickHouse:
		dialectorFactory = clickhouse.Open
	default:
		return nil, fmt.Errorf("unsupported database mode: %d", c.Mode)
	}

	for _, v := range c.Replicas {
		replicas = append(replicas, dialectorFactory(v))
	}
	for _, v := range c.Sources {
		sources = append(sources, dialectorFactory(v))
	}

	resolver := dbresolver.Register(dbresolver.Config{
		Sources:           sources,
		Replicas:          replicas,
		Policy:            dbresolver.RandomPolicy{},
		TraceResolverMode: c.Trace,
	})

	if c.MaxIdleConn > 0 {
		resolver.SetMaxIdleConns(c.MaxIdleConn)
	}
	if c.MaxOpenConn > 0 {
		resolver.SetMaxOpenConns(c.MaxOpenConn)
	}
	if c.MaxLifetime > 0 {
		resolver.SetConnMaxLifetime(time.Duration(c.MaxLifetime) * time.Minute)
	}

	if err := engine.Use(resolver); err != nil {
		return nil, fmt.Errorf("failed to use resolver: %v", err)
	}

	return engine, nil
}
