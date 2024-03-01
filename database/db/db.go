package db

import (
	"context"
	"errors"
	"time"

	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/pkg/opentelemetry/metric"
	"github.com/codfrm/cago/pkg/opentelemetry/trace"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/opentelemetry/tracing"
)

type contextKey int

const (
	dbKey contextKey = iota
)

var defaultDB *DB

type Driver string

const (
	MySQL      Driver = "mysql"
	Clickhouse Driver = "clickhouse"
)

type Config struct {
	Driver Driver `yaml:"driver"`
	Dsn    string `yaml:"dsn"`
	Prefix string `yaml:"prefix"`
	// 读写分离，以后再说吧
	//WriterDsn []string `yaml:"writerDsn,omitempty"` // 写入数据源
	//ReaderDsn []string `yaml:"readerDsn,omitempty"` // 读取数据源
}

type GroupConfig map[string]*Config

type DB struct {
	defaultDb *gorm.DB
	dbs       map[string]*gorm.DB
}

// Database gorm数据库封装
func Database() *DB {
	return &DB{}
}

func (d *DB) newDB(cfg *Config, debug bool) (*gorm.DB, error) {
	if cfg.Driver == "" {
		cfg.Driver = MySQL
	}
	logCfg := logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logger.Warn,
		IgnoreRecordNotFoundError: true,
		Colorful:                  false,
	}
	if debug {
		logCfg.IgnoreRecordNotFoundError = false
		logCfg.Colorful = true
	}
	orm, err := gorm.Open(driver[cfg.Driver](cfg), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   cfg.Prefix,
			SingularTable: true,
		},
		Logger: NewLogger(cfg.Driver, logCfg),
	})
	if err != nil {
		return nil, err
	}
	tracingPlugin := make([]tracing.Option, 0)
	if tp := trace.Default(); tp != nil {
		tracingPlugin = append(tracingPlugin,
			tracing.WithTracerProvider(tp),
		)
		if metric.Default() == nil {
			tracingPlugin = append(tracingPlugin,
				tracing.WithoutMetrics(),
			)
		}
	}
	if len(tracingPlugin) != 0 {
		if err := orm.Use(tracing.NewPlugin(
			tracingPlugin...,
		)); err != nil {
			return nil, err
		}
	}
	return orm, nil
}

func (d *DB) Start(ctx context.Context, config *configs.Config) error {
	cfg := &Config{}
	// 数据库group配置
	cfgGroup := make(GroupConfig)
	if ok, err := config.Has(ctx, "dbs"); err != nil {
		return err
	} else if ok {
		if err := config.Scan(ctx, "dbs", &cfgGroup); err != nil {
			return err
		}
	} else {
		if err := config.Scan(ctx, "db", cfg); err != nil {
			return err
		}
		cfgGroup["default"] = cfg
	}
	if _, ok := cfgGroup["default"]; !ok {
		return errors.New("no default db config")
	}
	dbs := make(map[string]*gorm.DB)
	orm, err := d.newDB(cfgGroup["default"], config.Debug)
	if err != nil {
		return err
	}
	delete(cfgGroup, "default")
	if len(cfgGroup) > 0 {
		for name, v := range cfgGroup {
			db, err := d.newDB(v, config.Debug)
			if err != nil {
				return err
			}
			dbs[name] = db
		}
	}
	d.defaultDb = orm
	d.dbs = dbs
	defaultDB = d
	return nil
}

func (d *DB) CloseHandle() {
	db, _ := d.defaultDb.DB()
	_ = db.Close()
	for _, v := range d.dbs {
		db, _ := v.DB()
		_ = db.Close()
	}
}

func Default() *gorm.DB {
	return defaultDB.defaultDb
}

func With(key string) *gorm.DB {
	if key == "default" {
		return defaultDB.defaultDb
	}
	return defaultDB.dbs[key]
}

func ContextWith(ctx context.Context, key string) context.Context {
	if key == "default" {
		return ContextWithDB(ctx, defaultDB.defaultDb)
	}
	return ContextWithDB(ctx, defaultDB.dbs[key])
}

func ContextWithDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, dbKey, db)
}

func Ctx(ctx context.Context) *gorm.DB {
	return CtxWith(ctx, "default")
}

func CtxWith(ctx context.Context, key string) *gorm.DB {
	if db, ok := ctx.Value(dbKey).(*gorm.DB); ok {
		return db.WithContext(ctx)
	}
	if key == "default" {
		return defaultDB.defaultDb.WithContext(ctx)
	}
	return defaultDB.dbs[key].WithContext(ctx)
}
