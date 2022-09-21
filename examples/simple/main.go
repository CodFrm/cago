package main

import (
	"context"
	"log"

	"github.com/codfrm/cago"
	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/mux"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/pkg/trace"
	"github.com/codfrm/cago/server"
)

// @title    api文档
// @version  1.0
// @BasePath /api/v1
func main() {
	ctx := context.Background()
	cfg, err := configs.NewConfig("simple")
	if err != nil {
		log.Fatalf("load config err: %v", err)
	}
	err = cago.New(ctx, cfg).
		Registry(cago.FuncComponent(logger.Logger)).
		Registry(cago.FuncComponent(trace.Trace)).
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

// Api
// @Author      CodFrm
// @Summary     一个测试API
// @Description 一个测试API描述
// @ID          example
// @Tags        example
// @Accept      json
// @Produce     json
// @Param       Request body     ApiRequest true "请求信息"
// @Success     200     {object} ApiResponse
// @Failure     400     {object} ApiFailResponse
// @Router      /api [POST]
