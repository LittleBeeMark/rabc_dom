package api

import (
	"github.com/gin-gonic/gin"

	"demo_casbin/bll"
	"demo_casbin/gin_util"
	"demo_casbin/pkg/logger"
	"demo_casbin/schema"
)

// Login 登录管理
type Login struct {
	LoginBll *bll.Login
}

// Login 用户登录
func (a *Login) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.LoginParam
	if err := gin_util.ParseJSON(c, &item); err != nil {
		gin_util.ResError(c, err)
		return
	}

	user, err := a.LoginBll.Verify(ctx, item.UserName, item.Password)
	if err != nil {
		gin_util.ResError(c, err)
		return
	}

	userID := user.ID

	// 将用户ID放入上下文
	gin_util.SetUserID(c, userID)
	schema.LoginUser = make(map[string]string)
	schema.LoginUser["login"] = userID

	ctx = logger.NewUserIDContext(ctx, userID)
	ctx = logger.NewTagContext(ctx, "__login__")
	logger.WithContext(ctx).Infof("登入系统")
	gin_util.ResSuccess(c, user)
}

// Logout 用户登出
func (a *Login) Logout(c *gin.Context) {
	ctx := c.Request.Context()

	// 检查用户是否处于登录状态，如果是则执行销毁
	userID := gin_util.GetUserID(c)
	if userID != "" {
		gin_util.ClearUserID(c)
		ctx = logger.NewTagContext(ctx, "__logout__")
		logger.WithContext(ctx).Infof("登出系统")
	}
	gin_util.ResOK(c)
}
