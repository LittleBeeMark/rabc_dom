package router

import (
	"github.com/gin-gonic/gin"

	"demo_casbin/middleware"
)

// RegisterAPI register api group router
func (a *Router) RegisterAPI(app *gin.Engine) {
	g := app.Group("/api")

	g.Use(middleware.CasbinMiddleware(a.CasbinEnforcer,
		middleware.AllowPathPrefixSkipper("/api/v1/pub"),
	))

	v1 := g.Group("/v1")
	{
		pub := v1.Group("/pub")
		{
			gLogin := pub.Group("login")
			{
				gLogin.POST("", a.LoginAPI.Login)
				gLogin.POST("exit", a.LoginAPI.Logout)
			}

		}

		gUser := v1.Group("user")
		{
			gUser.GET(":id", a.UserAPI.Get)
			gUser.POST("", a.UserAPI.Create)
			gUser.PUT(":id", a.UserAPI.Update)
			gUser.DELETE(":id", a.UserAPI.Delete)
		}

		gProject := v1.Group("project")
		{
			gProject.GET(":id", a.ProjectAPI.Get)
			gProject.POST("", a.ProjectAPI.Create)
			gProject.PUT(":id", a.ProjectAPI.Update)
			gProject.DELETE(":id", a.ProjectAPI.Delete)
		}

	}

}
