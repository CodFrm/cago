package db

import (
	"context"

	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/pkg/trace"
	mysqlDriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/opentelemetry/tracing"
)

var defaultDB *db

type db struct {
	*gorm.DB
}

type Config struct {
	Dsn    string `yaml:"dsn"`
	Prefix string `yaml:"prefix"`
}

func Database(ctx context.Context, config *configs.Config) error {
	cfg := &Config{}
	if err := config.Scan("db", cfg); err != nil {
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
	if tp := trace.Default(); tp != nil {
		if err := orm.Use(tracing.NewPlugin(
			tracing.WithTracerProvider(tp),
			tracing.WithoutMetrics(),
		)); err != nil {
			return err
		}
	}
	defaultDB = &db{
		DB: orm,
	}
	return nil
}

func Default() *gorm.DB {
	return defaultDB.DB
}

func Ctx(ctx context.Context) *gorm.DB {
	return defaultDB.DB.WithContext(ctx)
}
