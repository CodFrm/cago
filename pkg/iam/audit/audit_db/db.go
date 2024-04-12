package audit_db

import (
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuditLog struct {
	ID int64 `gorm:"primaryKey"`
	// 模块
	Module string `gorm:"column:module"`
	// 事件
	Event string `gorm:"column:event"`
	// 字段
	Fields string `gorm:"column:fields"`
}

type Database struct {
	db *gorm.DB
}

func NewDatabaseStorage(db *gorm.DB) (*Database, error) {
	// 创建表结构
	return &Database{
		db: db,
	}, nil
}

func (l *Database) Record(ctx context.Context, module, eventName string, fields ...zap.Field) error {
	return nil
}
