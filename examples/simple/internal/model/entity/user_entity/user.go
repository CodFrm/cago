package user_entity

import (
	"context"
	"github.com/codfrm/cago/examples/simple/internal/pkg/code"
	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/i18n"
)

type User struct {
	ID         int64  `gorm:"column:id;type:bigint(20);not null;primary_key"`
	Username   string `gorm:"column:username;type:varchar(255);index:username,unique;not null"`
	Status     int    `gorm:"column:status;type:int(11);not null"`
	Createtime int64  `gorm:"column:createtime;type:bigint(20)"`
	Updatetime int64  `gorm:"column:updatetime;type:bigint(20)"`
}

func (u *User) Check(ctx context.Context) error {
	if u == nil {
		return i18n.NewError(code.UserNotFound)
	}
	if u.Status != consts.ACTIVE {
		return i18n.NewError(code.UserIsBanned)
	}
	return nil
}
