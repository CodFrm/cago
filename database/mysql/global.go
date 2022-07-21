package mysql

import (
	"context"

	"github.com/codfrm/cago/configs"
	mysqlDriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var db *mysql

type mysql struct {
	*gorm.DB
}

type Config struct {
	Dsn    string `yaml:"dsn" env:"MYSQL_DSN"`
	Prefix string `yaml:"prefix,omitempty" env:"MYSQL_PREFIX"`
}

func Mysql(ctx context.Context, config *configs.Config) error {
	cfg := &Config{}
	if err := config.Scan("mysql", cfg); err != nil {
		return err
	}
	orm, err := gorm.Open(mysqlDriver.New(mysqlDriver.Config{
		DSN: cfg.Dsn,
	}), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   cfg.Prefix,
			SingularTable: true,
		},
	})
	if err != nil {
		return err
	}
	db = &mysql{
		DB: orm,
	}
	return nil
}

func Default() *gorm.DB {
	return db.DB
}
