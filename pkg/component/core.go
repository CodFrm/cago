package component

import (
	"context"
	"github.com/codfrm/cago"
	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/database/cache"
	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/database/elasticsearch"
	"github.com/codfrm/cago/database/mongo"
	"github.com/codfrm/cago/database/redis"
	"github.com/codfrm/cago/pkg/broker"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/pkg/opentelemetry/metric"
	"github.com/codfrm/cago/pkg/opentelemetry/trace"
)

// Core 核心组件,包括日志组件、链路追踪、指标
func Core() cago.FuncComponent {
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

// Database 数据库组件
func Database() cago.FuncComponent {
	return db.Database
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
func Cache() cago.FuncComponent {
	return cache.Cache
}

// Elasticsearch elasticsearch组件
func Elasticsearch() cago.FuncComponent {
	return elasticsearch.Elasticsearch
}
