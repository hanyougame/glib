package scopes

import (
	"fmt"
	"gorm.io/gorm"
)

// Equal 等于
func Equal(field string, value any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field+" = ?", value)
	}
}

// NotEqual 不等于
func NotEqual(field string, value any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field+" != ?", value)
	}
}

// Like 模糊查询
func Like(field string, value any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field+" LIKE ?", fmt.Sprintf("%%%v%%", value))
	}
}

// In in查询
func In[T comparable](field string, value []T) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field+" IN ?", value)
	}
}

// NotIn not in查询
func NotIn[T comparable](field string, value []T) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field+" NOT IN ?", value)
	}
}

// GT 大于
func GT(field string, value any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field+" > ?", value)
	}
}

// GTE 大于等于
func GTE(field string, value any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field+" >= ?", value)
	}
}

// LT 小于
func LT(field string, value any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field+" < ?", value)
	}
}

// LTE 小于等于
func LTE(field string, value any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field+" <= ?", value)
	}
}

// Between between查询
func Between[T comparable](field string, start, end T) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field+" BETWEEN ? AND ?", start, end)
	}
}

// NotBetween not between查询
func NotBetween[T comparable](field string, start, end T) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field+" NOT BETWEEN ? AND ?", start, end)
	}
}
