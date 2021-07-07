package router

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"

	"demo_casbin/api"
)

var _ IRouter = (*Router)(nil)

// IRouter 注册路由
type IRouter interface {
	Register(app *gin.Engine) error
	Prefixes() []string
}

// Router 路由管理器
type Router struct {
	CasbinEnforcer *casbin.SyncedEnforcer
	LoginAPI       *api.Login
	UserAPI        *api.User
	ProjectAPI     *api.Project
}

// Register 注册路由
func (a *Router) Register(app *gin.Engine) error {
	a.RegisterAPI(app)
	return nil
}

// Prefixes 路由前缀列表
func (a *Router) Prefixes() []string {
	return []string{
		"/api/",
	}
}
