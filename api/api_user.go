package api

import (
	"github.com/gin-gonic/gin"

	"demo_casbin/bll"
	"demo_casbin/gin_util"
	"demo_casbin/pkg/errors"
	"demo_casbin/schema"
)

// User 用户管理
type User struct {
	UserBll *bll.User
}

// Get 查询指定数据
func (a *User) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.UserBll.Get(ctx, c.Param("id"))
	if err != nil {
		gin_util.ResError(c, err)
		return
	}
	gin_util.ResSuccess(c, item.CleanSecure())
}

// Create 创建数据
func (a *User) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.User
	if err := gin_util.ParseJSON(c, &item); err != nil {
		gin_util.ResError(c, err)
		return
	} else if item.Password == "" {
		gin_util.ResError(c, errors.New400Response("密码不能为空"))
		return
	}

	result, err := a.UserBll.Create(ctx, item)
	if err != nil {
		gin_util.ResError(c, err)
		return
	}
	gin_util.ResSuccess(c, result)
}

// Update 更新数据
func (a *User) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.User
	if err := gin_util.ParseJSON(c, &item); err != nil {
		gin_util.ResError(c, err)
		return
	}

	err := a.UserBll.Update(ctx, c.Param("id"), item)
	if err != nil {
		gin_util.ResError(c, err)
		return
	}
	gin_util.ResOK(c)
}

// Delete 删除数据
func (a *User) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.UserBll.Delete(ctx, c.Param("id"))
	if err != nil {
		gin_util.ResError(c, err)
		return
	}
	gin_util.ResOK(c)
}
