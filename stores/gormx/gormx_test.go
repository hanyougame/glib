package gormx

import (
	"github.com/hanyougame/glib/stores/gormx/config"
	"github.com/hanyougame/glib/stores/gormx/scopes/paginate"
	"github.com/prometheus/common/model"
	"testing"
)

func TestMysqlEngine(t *testing.T) {
	mysqlDSN := "root:123456@tcp(127.0.0.1:3306)/video_live?charset=utf8mb4&parseTime=True&loc=Local"
	c := config.Config{
		Mode: config.Mysql,
		DSN:  mysqlDSN,
	}
	db := NewEngine(c)

	t.Log(db.Error)
}

func TestPagination(t *testing.T) {
	mysqlDSN := "root:123456@tcp(127.0.0.1:3306)/video_live?charset=utf8mb4&parseTime=True&loc=Local"
	c := config.Config{
		Mode:  config.Mysql,
		DSN:   mysqlDSN,
		Debug: true,
	}
	db := NewEngine(c)

	var page = &paginate.Pagination[[]model.Alert]{
		Page:     1,
		PageSize: 10,
	}

	db.
		Table("users").
		Where("user_id = 1").
		Scopes(
			paginate.Paginate(page),
		)
}
