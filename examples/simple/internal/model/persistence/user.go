package persistence

import (
	"context"

	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/examples/simple/internal/model/entity"
	"github.com/codfrm/cago/examples/simple/internal/repository"
)

type user struct {
}

func NewUser() repository.IUser {
	return &user{}
}

func (u *user) Find(ctx context.Context, id int64) (*entity.User, error) {
	ret := &entity.User{ID: id}
	if err := db.Ctx(ctx).First(ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (u *user) Create(ctx context.Context, user *entity.User) error {
	return db.Ctx(ctx).Create(user).Error
}

func (u *user) Update(ctx context.Context, user *entity.User) error {
	return db.Ctx(ctx).Updates(user).Error
}

func (u *user) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Delete(&entity.User{ID: id}).Error
}
