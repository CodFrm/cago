package migrations

import (
	"github.com/codfrm/cago/examples/simple/internal/model/entity/user_entity"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func T20230611() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20230611",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				&user_entity.User{},
			)
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	}
}
