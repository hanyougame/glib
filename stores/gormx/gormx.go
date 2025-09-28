package gormx

import (
	"fmt"
	"github.com/hanyougame/glib/stores/gormx/config"
	"github.com/hanyougame/glib/stores/gormx/database"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
)

var (
	once   sync.Once
	Engine *DBManager
)

// DBManager gorm db manager
type DBManager struct {
	sync.RWMutex
	Mysql      *gorm.DB
	Postgres   *gorm.DB
	ClickHouse *gorm.DB
}

func NewDBManager() *DBManager {
	return &DBManager{}
}

// Initialize initialize the database
func (dm *DBManager) Initialize(cs ...config.Config) error {
	for _, c := range cs {
		if err := dm.newEngine(c); err != nil {
			return err
		}
	}
	return nil
}

// Must initialize the database
func Must(cs ...config.Config) {
	once.Do(func() {
		Engine = NewDBManager()
		if err := Engine.Initialize(cs...); err != nil {
			panic(fmt.Sprintf("failed to initialize databases: %v", err))
		}
	})
}

func (dm *DBManager) newEngine(c config.Config) error {
	var (
		dialector gorm.Dialector
		dbKey     string
	)

	switch c.Mode {
	case config.Mysql:
		dialector = mysql.Open(c.DSN)
		dbKey = "Mysql"
	case config.Postgres:
		dialector = postgres.Open(c.DSN)
		//dialector = postgres.New(postgres.Config{
		//	DSN:                  c.DSN,
		//	PreferSimpleProtocol: true, // disables implicit prepared statement usage
		//})
		dbKey = "Postgres"
	case config.ClickHouse:
		dialector = clickhouse.Open(c.DSN)
		dbKey = "ClickHouse"
	default:
		return fmt.Errorf("unsupported database mode: %d", c.Mode)
	}

	engine, err := database.NewEngine(c, dialector)
	if err != nil {
		return fmt.Errorf("failed to initialize %s database: %v", dbKey, err)
	}

	dm.Lock()
	defer dm.Unlock()

	switch c.Mode {
	case config.Mysql:
		if dm.Mysql == nil {
			dm.Mysql = engine
		} else {
			return fmt.Errorf("mysql connection already exists")
		}
	case config.Postgres:
		if dm.Postgres == nil {
			dm.Postgres = engine
		} else {
			return fmt.Errorf("postgres connection already exists")
		}
	case config.ClickHouse:
		if dm.ClickHouse == nil {
			dm.ClickHouse = engine
		} else {
			return fmt.Errorf("clickhouse connection already exists")
		}
	}

	return nil
}

// CloseAll close all connections
func (dm *DBManager) CloseAll() error {
	var errs []error
	dm.Lock()
	defer dm.Unlock()

	closeDB := func(dbName string, db *gorm.DB) {
		if db != nil {
			sqlDB, err := db.DB()
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to get underlying sql.DB for %s: %v", dbName, err))
				return
			}

			if e := sqlDB.Close(); e != nil {
				errs = append(errs, fmt.Errorf("failed to close %s connection: %v", dbName, e))
			}
		}
	}

	closeDB("MySQL", dm.Mysql)
	closeDB("Postgres", dm.Postgres)
	closeDB("ClickHouse", dm.ClickHouse)

	if len(errs) > 0 {
		return fmt.Errorf("errors occurred while closing connections: %v", errs)
	}
	return nil
}
