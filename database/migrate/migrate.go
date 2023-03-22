package migrate

import (
	"context"
	"errors"
	"reflect"

	mongo2 "github.com/codfrm/cago/database/migrate/mongomigrate"
	mongo3 "github.com/codfrm/cago/database/mongo"
	gormigrate "github.com/go-gormigrate/gormigrate/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// MigrateFunc 数据库迁移函数
type MigrateFunc[T any] func(ctx context.Context, db T) error

// RollbackFunc 数据库回滚函数
type RollbackFunc[T any] func(ctx context.Context, db T) error

type Migration[T any] struct {
	ID       string
	Migrate  MigrateFunc[T]
	Rollback RollbackFunc[T]
}

// RunMigrations 数据库迁移操作
func RunMigrations[T *gorm.DB | *mongo.Client](
	ctx context.Context,
	db T, fs ...func() *Migration[T],
) error {
	v := reflect.ValueOf(db)

	switch v.Type().String() {
	case "*gorm.DB":
		ms := make([]*gormigrate.Migration, 0, len(fs))
		for _, f := range fs {
			m := f()
			ms = append(ms, &gormigrate.Migration{
				ID: m.ID,
				Migrate: func(db *gorm.DB) error {
					v := reflect.ValueOf(db)
					return m.Migrate(ctx, v.Interface())
				},
				Rollback: func(db *gorm.DB) error {
					v := reflect.ValueOf(db)
					return m.Rollback(ctx, v.Interface())
				},
			})
		}
		m := gormigrate.New(v.Interface().(*gorm.DB), &gormigrate.Options{
			TableName:                 "migrations",
			IDColumnName:              "id",
			IDColumnSize:              200,
			UseTransaction:            true,
			ValidateUnknownMigrations: true,
		}, ms)
		return m.Migrate()
	case "*mongo.Client":
		ms := make([]*mongo2.Migration, 0, len(fs))
		for _, f := range fs {
			m := f()
			ms = append(ms, &mongo2.Migration{
				ID: m.ID,
				Migrate: func(ctx context.Context, db *mongo3.Client) error {
					v := reflect.ValueOf(db)
					return m.Migrate(ctx, v.Interface())
				},
				Rollback: func(ctx context.Context, db *mongo3.Client) error {
					v := reflect.ValueOf(db)
					return m.Rollback(ctx, v.Interface())
				},
			})
		}
		m := mongo2.New(ctx, v.Interface().(*mongo3.Client), ms)
		return m.Migrate()
	default:
		return errors.New("unsupported database type")
	}
}
