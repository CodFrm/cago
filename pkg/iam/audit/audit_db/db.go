package audit_db

import (
	"context"
	"encoding/json"

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
	// 创建时间
	Createtime int64 `gorm:"column:createtime"`
}

type Database struct {
	db *gorm.DB
}

func NewDatabaseStorage(db *gorm.DB) (*Database, error) {
	// 创建表结构
	if err := db.Migrator().AutoMigrate(&AuditLog{}); err != nil {
		return nil, err
	}
	return &Database{
		db: db,
	}, nil
}

func (l *Database) Record(ctx context.Context, module, eventName string, fields ...zap.Field) error {
	data, err := json.Marshal(fields)
	if err != nil {
		return err
	}
	if err := l.db.WithContext(ctx).Create(&AuditLog{
		ID:     0,
		Module: module,
		Event:  eventName,
		Fields: string(data),
	}).Error; err != nil {
		return err
	}
	return nil
}
