package database

import (
	"github.com/hanyougame/glib/stores/gormx/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
	"time"
)

func NewEngine(c config.Config, dialector gorm.Dialector, opt ...gorm.Option) *gorm.DB {
	cfg := gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}

	engine, err := gorm.Open(dialector, append(opt, &cfg)...)

	if err != nil {
		panic(err)
	}

	if c.Debug {
		engine = engine.Debug()
	}

	if c.Trace {
		registerTraceHook(engine)
	}

	if c.Separation { //读写分离
		return newRWDBEngine(c, engine)
	}

	sqlDB, err := engine.DB()

	if err != nil {
		panic(err)
	}

	if c.MaxIdleConn > 0 {
		sqlDB.SetMaxIdleConns(c.MaxIdleConn)
	}

	if c.MaxOpenConn > 0 {
		sqlDB.SetMaxOpenConns(c.MaxOpenConn)
	}

	if c.MaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Minute * time.Duration(c.MaxLifetime))
	}
	return engine
}

func newRWDBEngine(c config.Config, engine *gorm.DB) *gorm.DB {
	var replicas, sources []gorm.Dialector
	for _, v := range c.Replicas {
		replicas = append(replicas, mysql.New(mysql.Config{
			DSN: v,
		}))
	}
	for _, v := range c.Sources {
		sources = append(sources, mysql.New(mysql.Config{
			DSN: v,
		}))
	}
	resolver := dbresolver.Register(dbresolver.Config{
		Sources:  sources,
		Replicas: replicas,
		Policy:   dbresolver.RandomPolicy{},
	})
	if c.MaxIdleConn > 0 {
		resolver.SetMaxIdleConns(c.MaxIdleConn)
	}
	if c.MaxOpenConn > 0 {
		resolver.SetMaxOpenConns(c.MaxOpenConn)
	}
	if c.MaxLifetime > 0 {
		resolver.SetConnMaxLifetime(time.Minute * time.Duration(c.MaxLifetime))
	}
	if err := engine.Use(resolver); err != nil {
		panic(err)
	}
	return engine
}
