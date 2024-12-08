package paginate

import (
	"gorm.io/gorm"
)

// Paginate 分页
func Paginate[T any](pagination *Pagination[T]) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var totalRows int64
		db.Session(&gorm.Session{}).Model(db.Statement.Model).Count(&totalRows)
		pagination.Total = totalRows
		if totalRows == 0 {
			return db
		}
		return db.Offset(pagination.Offset()).Limit(pagination.Limit())
	}
}
