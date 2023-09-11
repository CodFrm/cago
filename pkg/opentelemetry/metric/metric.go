package metric

import (
	"context"

	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/server/mux"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric"
)

func init() {
	mux.RegisterMiddleware(func(cfg *configs.Config, r *gin.Engine) error {
		// 加入metrics中间件
		if Default() != nil {
			r.GET("/metrics", gin.WrapH(promhttp.Handler()))
		}
		return nil
	})
}

var provider *metric.MeterProvider

func Metrics(ctx context.Context, cfg *configs.Config) error {
	// 初始化全局Meter实例并绑定Prometheus Exporter
	exporter, err := prometheus.New()
	if err != nil {
		return err
	}
	provider = metric.NewMeterProvider(metric.WithReader(exporter))
	global.SetMeterProvider(provider)
	return nil
}

func Default() *metric.MeterProvider {
	return provider
}
