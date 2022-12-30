package repository

import (
	"context"

	"github.com/codfrm/cago/examples/simple/internal/model/entity"
)

type UserRepo interface {
	Find(ctx context.Context, id int64) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id int64) error
}

var defaultUser UserRepo

func User() UserRepo {
	return defaultUser
}

func RegisterUser(i UserRepo) {
	defaultUser = i
}
