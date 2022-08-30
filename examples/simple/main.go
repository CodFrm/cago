package main

import (
	"context"
	"log"

	"github.com/codfrm/cago"
	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/configs/file"
	"github.com/codfrm/cago/mux"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/server"
)

func main() {
	ctx := context.Background()
	source, err := file.NewSource("config.yaml", file.Yaml())
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := configs.NewConfig("simple", source)
	if err != nil {
		log.Fatalf("load config err: %v", err)
	}
	err = cago.New(ctx, cfg).
		Registry(cago.FuncComponent(logger.Logger)).
		//Registry(cago.FuncComponent(mysql.Mysql)).
		RegistryCancel(server.Http(func(r *mux.RouterGroup) error {
			r.GET("/", func(ctx *mux.WebContext) {
				_, _ = ctx.Writer.Write([]byte("hello world"))
				ctx.Logger().Info("hello world")
			})
			return nil
		})).
		Registry(cago.FuncComponent(func(ctx context.Context, cfg *configs.Config) error {
			logger.Default().Info("hello world")
			return nil
		})).
		Start()
	if err != nil {
		log.Fatalf("start err: %v", err)
		return
	}
}
