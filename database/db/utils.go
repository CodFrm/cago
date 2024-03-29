package db

import (
	"errors"

	mysqlDriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var driver = map[Driver]func(*Config) gorm.Dialector{
	MySQL: func(config *Config) gorm.Dialector {
		return mysqlDriver.New(mysqlDriver.Config{
			DSN: config.Dsn,
		})
	},
}

// RegisterDriver 注册数据库驱动 默认会注册mysql
// t 数据库驱动名
// f 数据库驱动函数 传入配置返回 gorm.Dialector
func RegisterDriver(t Driver, f func(*Config) gorm.Dialector) {
	driver[t] = f
}

// RecordNotFound 判断是否是记录不存在的错误
func RecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
