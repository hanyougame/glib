package paginate

import (
	"gorm.io/gorm"
)

// Paginate 分页
func Paginate(pagination *Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pagination == nil {
			return db
		}

		var (
			tx        = db.Session(&gorm.Session{})
			totalRows int64
		)

		tx.Model(db.Statement.Model).Count(&totalRows)
		pagination.Total = totalRows
		return db.Offset(pagination.Offset()).Limit(pagination.Limit())
	}
}
