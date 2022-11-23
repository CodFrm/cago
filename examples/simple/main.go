package main

import (
	"context"
	"log"

	"github.com/codfrm/cago"
	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/examples/simple/internal/api"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/pkg/trace"
	"github.com/codfrm/cago/server/mux"
)

func main() {
	ctx := context.Background()
	cfg, err := configs.NewConfig("simple")
	if err != nil {
		log.Fatalf("load config err: %v", err)
	}
	err = cago.New(ctx, cfg).
		Registry(cago.FuncComponent(logger.Logger)).
		Registry(cago.FuncComponent(trace.Trace)).
		//Registry(cago.FuncComponent(db.Database)).
		RegistryCancel(mux.Http(api.Router)).
		Registry(cago.FuncComponent(func(ctx context.Context, cfg *configs.Config) error {
			logger.Default().Info("cago simple example start")
			go func() {
				<-ctx.Done()
				logger.Default().Info("cago simple example stop")
			}()
			return nil
		})).
		Start()
	if err != nil {
		log.Fatalf("start err: %v", err)
		return
	}
}
