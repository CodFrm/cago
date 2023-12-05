package component

import (
	"context"

	"github.com/codfrm/cago"
	"github.com/codfrm/cago/configs"
	_ "github.com/codfrm/cago/configs/etcd"
	"github.com/codfrm/cago/database/cache"
	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/database/elasticsearch"
	"github.com/codfrm/cago/database/mongo"
	"github.com/codfrm/cago/database/redis"
	"github.com/codfrm/cago/pkg/broker"
	"github.com/codfrm/cago/pkg/logger"
	_ "github.com/codfrm/cago/pkg/logger/loki"
	"github.com/codfrm/cago/pkg/opentelemetry/metric"
	"github.com/codfrm/cago/pkg/opentelemetry/trace"
	"github.com/codfrm/cago/server/mux"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Core 核心组件,包括日志组件、链路追踪、指标
func Core() cago.FuncComponent {
	mux.RegisterMiddleware(func(cfg *configs.Config, r *gin.Engine) error {
		if cfg.Env != configs.PROD {
			url := ginSwagger.URL("/swagger/doc.json")
			r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
		}
		return nil
	})
	return func(ctx context.Context, cfg *configs.Config) error {
		// 日志组件必须注册
		if err := logger.Logger(ctx, cfg); err != nil {
			return err
		}
		// 判断是否有trace配置
		if ok, err := cfg.Has("trace"); err != nil {
			return err
		} else if ok {
			if err := trace.Trace(ctx, cfg); err != nil {
				return err
			}
		}
		// metrics组件
		if err := metric.Metrics(ctx, cfg); err != nil {
			return err
		}
		return nil
	}
}

func Logger() cago.FuncComponent {
	return logger.Logger
}

func Trace() cago.FuncComponent {
	return trace.Trace
}

func Metrics() cago.FuncComponent {
	return metric.Metrics
}

// Database 数据库组件
func Database() cago.Component {
	return db.Database()
}

// Broker 消息队列组件
func Broker() cago.FuncComponent {
	return broker.Broker
}

// Mongo mongodb组件
func Mongo() cago.FuncComponent {
	return mongo.Mongo
}

// Redis redis组件
func Redis() cago.FuncComponent {
	return redis.Redis
}

// Cache 缓存组件
func Cache() cago.Component {
	return cache.Cache()
}

// Elasticsearch elasticsearch组件
func Elasticsearch() cago.FuncComponent {
	return elasticsearch.Elasticsearch
}
