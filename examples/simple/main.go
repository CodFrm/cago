package main

import (
	"context"
	"log"

	"github.com/codfrm/cago"
	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/examples/simple/internal/controller"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/pkg/trace"
	"github.com/codfrm/cago/server/http"
)

// main
// @title    api文档
// @version  1.0
// @BasePath /api/v1
// @Body     x-www-form-urlencoded
func main() {
	ctx := context.Background()
	cfg, err := configs.NewConfig("simple")
	if err != nil {
		log.Fatalf("load config err: %v", err)
	}
	err = cago.New(ctx, cfg).
		Registry(cago.FuncComponent(logger.Logger)).
		Registry(cago.FuncComponent(trace.Trace)).
		Registry(cago.FuncComponent(db.DB)).
		RegistryCancel(http.Http(func(r *http.Router) error {
			return r.Group("/").Bind(
				controller.NewUser(),
			)
		})).
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
