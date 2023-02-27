package clickhouse

import (
	"context"

	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/pkg/trace"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

var defaultDB *gorm.DB

type Config struct {
	Dsn    string `yaml:"dsn"`
	Prefix string `yaml:"prefix"`
}

func Clickhouse(ctx context.Context, config *configs.Config) error {
	cfg := &Config{}
	if err := config.Scan("db", cfg); err != nil {
		return err
	}
	db, err := gorm.Open(clickhouse.Open(cfg.Dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	if tp := trace.Default(); tp != nil {
		if err := db.Use(tracing.NewPlugin(
			tracing.WithTracerProvider(tp),
			tracing.WithoutMetrics(),
		)); err != nil {
			return err
		}
	}
	defaultDB = db
	return nil
}

func Default() *gorm.DB {
	return defaultDB
}

func Ctx(ctx context.Context) *gorm.DB {
	return defaultDB.WithContext(ctx)
}