package main

import (
	"context"
	"github.com/codfrm/cago/examples/simple/internal/repository/user_repo"
	"log"

	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/examples/simple/internal/task/consumer"
	"github.com/codfrm/cago/examples/simple/migrations"
	"github.com/codfrm/cago/pkg/component"

	"github.com/codfrm/cago"
	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/examples/simple/internal/api"
	"github.com/codfrm/cago/server/mux"
)

func main() {
	ctx := context.Background()
	cfg, err := configs.NewConfig("simple")
	if err != nil {
		log.Fatalf("load config err: %v", err)
	}

	// 注册储存实例
	user_repo.RegisterUser(user_repo.NewUser())

	err = cago.New(ctx, cfg).
		Registry(component.Core()).
		Registry(component.Database()).
		Registry(component.Broker()).
		Registry(component.Redis()).
		Registry(component.Cache()).
		Registry(consumer.Consumer()).
		Registry(cago.FuncComponent(func(ctx context.Context, cfg *configs.Config) error {
			return migrations.RunMigrations(db.Default())
		})).
		RegistryCancel(mux.HTTP(api.Router)).
		Start()
	if err != nil {
		log.Fatalf("start err: %v", err)
		return
	}
}
