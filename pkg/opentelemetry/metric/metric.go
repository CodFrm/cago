package metric

import (
	"context"

	"github.com/codfrm/cago"
	"github.com/codfrm/cago/configs"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric"
)

var provider *metric.MeterProvider

func Metrics() cago.FuncComponent {
	return func(ctx context.Context, cfg *configs.Config) error {
		// 初始化全局Meter实例并绑定Prometheus Exporter
		exporter, err := prometheus.New()
		if err != nil {
			return err
		}
		provider = metric.NewMeterProvider(metric.WithReader(exporter))
		global.SetMeterProvider(provider)
		return nil
	}
}

func Default() *metric.MeterProvider {
	return provider
}
