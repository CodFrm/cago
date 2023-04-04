package clickhouse

import (
	"context"

	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/pkg/opentelemetry/metric"
	"github.com/codfrm/cago/pkg/opentelemetry/trace"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

var defaultDB *gorm.DB

type Config struct {
	Dsn string `yaml:"dsn"`
}

func Clickhouse(ctx context.Context, config *configs.Config) error {
	cfg := &Config{}
	if err := config.Scan("clickhouse", cfg); err != nil {
		return err
	}
	db, err := gorm.Open(clickhouse.Open(cfg.Dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	tracingPlugin := make([]tracing.Option, 0)
	if tp := trace.Default(); tp != nil {
		tracingPlugin = append(tracingPlugin,
			tracing.WithTracerProvider(tp),
			tracing.WithDBName("clickhouse"),
		)
		if metric.Default() == nil {
			tracingPlugin = append(tracingPlugin,
				tracing.WithoutMetrics(),
			)
		}
	}
	if len(tracingPlugin) != 0 {
		if err := db.Use(tracing.NewPlugin(
			tracingPlugin...,
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
