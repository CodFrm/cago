package user_repo

import (
	"context"

	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/examples/simple/internal/model/entity/user_entity"
	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/utils/httputils"
)

type UserRepo interface {
	Find(ctx context.Context, id int64) (*user_entity.User, error)
	FindPage(ctx context.Context, page httputils.PageRequest) ([]*user_entity.User, int64, error)
	Create(ctx context.Context, user *user_entity.User) error
	Update(ctx context.Context, user *user_entity.User) error
	Delete(ctx context.Context, id int64) error

	FindByUsername(ctx context.Context, username string) (*user_entity.User, error)
}

var defaultUser UserRepo

func User() UserRepo {
	return defaultUser
}

func RegisterUser(i UserRepo) {
	defaultUser = i
}

type userRepo struct {
}

func NewUser() UserRepo {
	return &userRepo{}
}

func (u *userRepo) Find(ctx context.Context, id int64) (*user_entity.User, error) {
	ret := &user_entity.User{}
	if err := db.Ctx(ctx).Where("id=? and status=?", id, consts.ACTIVE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *userRepo) Create(ctx context.Context, user *user_entity.User) error {
	return db.Ctx(ctx).Create(user).Error
}

func (u *userRepo) Update(ctx context.Context, user *user_entity.User) error {
	return db.Ctx(ctx).Updates(user).Error
}

func (u *userRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Model(&user_entity.User{}).Where("id=?", id).Update("status", consts.DELETE).Error
}

func (u *userRepo) FindPage(ctx context.Context, page httputils.PageRequest) ([]*user_entity.User, int64, error) {
	var list []*user_entity.User
	var count int64
	find := db.Ctx(ctx).Model(&user_entity.User{}).Where("status=?", consts.ACTIVE)
	if err := find.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := find.Order("createtime desc").Offset(page.GetOffset()).Limit(page.GetLimit()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (u *userRepo) FindByUsername(ctx context.Context, username string) (*user_entity.User, error) {
	ret := &user_entity.User{}
	if err := db.Ctx(ctx).Where("username=?", username).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}
