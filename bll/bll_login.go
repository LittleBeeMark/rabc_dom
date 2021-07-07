package bll

import (
	"context"

	"demo_casbin/model/gorm/model"
	"demo_casbin/pkg/errors"
	"demo_casbin/schema"
)

// Login 登录管理
type Login struct {
	UserModel *model.User
}

func (a *Login) checkAndGetUser(ctx context.Context, userID string) (*schema.User, error) {
	user, err := a.UserModel.Get(ctx, userID)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, errors.ErrInvalidUser
	}
	return user, nil
}

// Verify 登录验证
func (a *Login) Verify(ctx context.Context, userName, password string) (*schema.User, error) {

	result, err := a.UserModel.Query(ctx, schema.UserQueryParam{
		UserName: userName,
	})
	if err != nil {
		return nil, err
	} else if len(result.Data) == 0 {
		return nil, errors.ErrInvalidUserName
	}

	item := result.Data[0]
	if item.Password != password {
		return nil, errors.ErrInvalidPassword
	}

	return item, nil
}
