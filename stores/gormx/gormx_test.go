package gormx

import (
	"github.com/hanyougame/glib/stores/gormx/config"
	"github.com/zeromicro/go-zero/core/jsonx"
	"gorm.io/plugin/dbresolver"
	"testing"
)

func TestPostgresEngine(t *testing.T) {
	dsn := "host=192.168.6.218 user=postgresql password=bingtangMySQL dbname=k1-dev port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	c := config.Config{
		Mode:  config.Postgres,
		DSN:   dsn,
		Debug: true,
	}

	Must(c)
	result := map[string]interface{}{}
	err := Engine.Postgres.Table("roles").Model(&Model{}).FirstOrInit(&result).Error
	if err != nil {
		t.Error(err)
	}
	t.Log("result", result)
}

type User struct {
	Name string
	Age  int
}

func (u *User) TableName() string {
	return "users"
}

func TestClickHouseEngine(t *testing.T) {
	dsn := "clickhouse://:948Q3x7K@192.168.6.218:9000/default?dial_timeout=10s&read_timeout=20s"
	c := config.Config{
		Mode:  config.ClickHouse,
		DSN:   dsn,
		Debug: true,
	}

	Must(c)

	var users = []User{
		{Name: "小明", Age: 15},
		{Name: "小红", Age: 16},
		{Name: "小王", Age: 17},
	}
	err := Engine.ClickHouse.Create(&users).Error
	if err != nil {
		t.Error(err)
	}

	var list []User
	err = Engine.ClickHouse.Find(&list).Error
	if err != nil {
		t.Error(err)
	}
	t.Log(jsonx.MarshalToString(list))
}

//
//func TestMysqlEngine(t *testing.T) {
//	mysqlDSN := "root:123456@tcp(127.0.0.1:3306)/video_live?charset=utf8mb4&parseTime=True&loc=Local"
//	c := config.Config{
//		Mode: config.Mysql,
//		DSN:  mysqlDSN,
//	}
//	db := NewEngine(c)
//
//	t.Log(db.Error)
//}
//
//func TestPagination(t *testing.T) {
//	mysqlDSN := "root:123456@tcp(127.0.0.1:3306)/video_live?charset=utf8mb4&parseTime=True&loc=Local"
//	c := config.Config{
//		Mode:  config.Mysql,
//		DSN:   mysqlDSN,
//		Debug: true,
//	}
//	db := NewEngine(c)
//
//	var page = &paginate.Pagination[[]model.Alert]{
//		Page:     1,
//		PageSize: 10,
//	}
//
//	db.
//		Table("users").
//		Where("user_id = 1").
//		Scopes(
//			paginate.Paginate(page),
//		)
//}

func TestWriteRead(t *testing.T) {
	dsn := "host=192.168.6.218 user=postgresql password=bingtangMySQL dbname=k1-dev port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	dsn2 := "host=192.168.6.218 user=postgresql password=bingtangMySQL dbname=k1-dev port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	c := config.Config{
		Mode:       config.Postgres,
		Trace:      true,
		DSN:        dsn,
		Debug:      true,
		Separation: true,
		Sources:    []string{dsn},
		Replicas:   []string{dsn2},
	}

	Must(c)
	result := map[string]interface{}{}
	err := Engine.Postgres.Clauses(dbresolver.Write).Table("roles").Model(&Model{}).FirstOrInit(&result).Error
	if err != nil {
		t.Error(err)
	}
	t.Log("result", result)
}
