package cron

import (
	"context"

	"github.com/codfrm/cago"
	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/pkg/trace"
	"github.com/robfig/cron/v3"
	trace2 "go.opentelemetry.io/otel/trace"
)

const (
	tracerName = "github.com/codfrm/cago/server/cron"
)

type Callback func(r Crontab) error

type server struct {
	//ctx context.Context
	//cancel   context.CancelFunc
	cron     *cron.Cron
	callback Callback
}

// Cron 定时任务组件,需要先注册logger和redis组件
func Cron(callback Callback) cago.Component {

	return &server{
		cron:     cron.New(),
		callback: callback,
	}
}

func (s *server) Start(ctx context.Context, cfg *configs.Config) error {
	var tracer trace2.Tracer
	if trace.Default() != nil {
		tracer = trace.Default().Tracer(
			tracerName,
			trace2.WithInstrumentationVersion("0.1.0"),
		)
	}
	if err := s.callback(&crontab{tracer: tracer, cron: s.cron}); err != nil {
		return err
	}
	s.cron.Start()
	return nil
}

func (s *server) CloseHandle() {
	s.cron.Stop()
}
