package gormx

import (
	"errors"
	"gorm.io/gorm"
)

// NotFound 判断是否未未找到错误
func NotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// IsUniqueError 是否是唯一索引错误
func IsUniqueError(err error) bool {
	return errors.Is(err, gorm.ErrDuplicatedKey)
}
