package mysql

import (
	"context"

	"github.com/codfrm/cago"
	"github.com/codfrm/cago/config"
	"github.com/jinzhu/gorm"
)

var db *mysql

type mysql struct {
	*gorm.DB
}

type mysqlConfig struct {
	Dsn    string `yaml:"dsn" env:"MYSQL_DSN"`
	Prefix string `yaml:"prefix" env:"MYSQL_PREFIX"`
}

func Mysql(ctx context.Context, config *config.Config) error {
	cfg := &mysqlConfig{}
	if err := config.Scan("mysql", cfg); err != nil {
		return err
	}
	orm, err := gorm.Open(cfg.Dsn)
	if err != nil {
		return err
	}
	db = &mysql{
		DB: orm,
	}
	return nil
}

func Ctx(ctx cago.Context) *gorm.DB {
	return db.DB.New()
}
