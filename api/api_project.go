package api

import (
	"github.com/gin-gonic/gin"

	"demo_casbin/bll"
	"demo_casbin/gin_util"
	"demo_casbin/schema"
)

// Project 项目管理
type Project struct {
	ProjectBll *bll.Project
}

// Get 查询指定数据
func (a *Project) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.ProjectBll.Get(ctx, c.Param("id"))
	if err != nil {
		gin_util.ResError(c, err)
		return
	}
	gin_util.ResSuccess(c, item)
}

// Create 创建数据
func (a *Project) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Project
	if err := gin_util.ParseJSON(c, &item); err != nil {
		gin_util.ResError(c, err)
		return
	}

	result, err := a.ProjectBll.Create(ctx, item)
	if err != nil {
		gin_util.ResError(c, err)
		return
	}
	gin_util.ResSuccess(c, result)
}

// Update 更新数据
func (a *Project) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Project
	if err := gin_util.ParseJSON(c, &item); err != nil {
		gin_util.ResError(c, err)
		return
	}

	err := a.ProjectBll.Update(ctx, c.Param("id"), item)
	if err != nil {
		gin_util.ResError(c, err)
		return
	}
	gin_util.ResOK(c)
}

// Delete 删除数据
func (a *Project) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.ProjectBll.Delete(ctx, c.Param("id"))
	if err != nil {
		gin_util.ResError(c, err)
		return
	}
	gin_util.ResOK(c)
}
