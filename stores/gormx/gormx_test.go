package gormx

import (
	"github.com/hanyougame/glib/stores/gormx/config"
	"github.com/zeromicro/go-zero/core/jsonx"
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
	dsn := "clickhouse://127.0.0.1:9000/default?dial_timeout=10s&read_timeout=20s"
	//db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{
	//	Logger: logger.Default.LogMode(logger.Info),
	//})
	//if err != nil {
	//	panic(err)
	//}

	//db.AutoMigrate(&User{})
	//db.Create(&User{Name: "Tom", Age: 18})
	var user User
	//db.First(&user)
	//t.Log(jsonx.MarshalToString(user))
	//t.Log(db)
	c := config.Config{
		Mode: config.ClickHouse,
		DSN:  dsn,
	}

	Must(c)
	//t.Log(Engine.ClickHouse.AutoMigrate(&User{}).Error())
	//Engine.ClickHouse.Create(&User{Name: "Jerry", Age: 15})
	err := Engine.ClickHouse.First(&user, "name=?", "Tom").Error
	if err != nil {
		t.Error(err)
	}

	t.Log(jsonx.MarshalToString(user))
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
