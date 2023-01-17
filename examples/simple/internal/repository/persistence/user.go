package persistence

import (
	"context"

	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/examples/simple/internal/model/entity"
	"github.com/codfrm/cago/examples/simple/internal/repository"
	"github.com/codfrm/cago/pkg/utils/httputils"
)

type userRepo struct {
}

func NewUser() repository.UserRepo {
	return &userRepo{}
}

func (u *userRepo) Find(ctx context.Context, id int64) (*entity.User, error) {
	ret := &entity.User{ID: id}
	if err := db.Ctx(ctx).First(ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (u *userRepo) FindPage(ctx context.Context, page httputils.PageRequest) ([]*entity.User, int64, error) {
	var users []*entity.User
	var count int64
	if err := db.Ctx(ctx).Model(&entity.User{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Ctx(ctx).Offset(page.GetOffset()).Limit(page.GetLimit()).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func (u *userRepo) Create(ctx context.Context, user *entity.User) error {
	return db.Ctx(ctx).Create(user).Error
}

func (u *userRepo) Update(ctx context.Context, user *entity.User) error {
	return db.Ctx(ctx).Updates(user).Error
}

func (u *userRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Delete(&entity.User{ID: id}).Error
}
