package gormx

import "gorm.io/gorm"

type BaseDAO interface {
	Model() *gorm.DB
}
