package gormx

import "gorm.io/plugin/soft_delete"

// Model 基础模型
type Model struct {
	CreatedAt int64                 `json:"created_at"`
	UpdatedAt int64                 `json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `json:"deleted_at" gorm:"default:0" `
}
