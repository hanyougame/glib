package gormx

import (
	"errors"
	"gorm.io/gorm"
	"strings"
)

const uniqueErrorKey = "duplicate key"

// NotFound 判断是否未未找到错误
func NotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// IsUniqueError 是否是唯一索引错误
func IsUniqueError(err error) bool {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	if err != nil && strings.Contains(err.Error(), uniqueErrorKey) {
		return true
	}

	return false
}
