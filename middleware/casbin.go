package middleware

import (
	"fmt"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"

	"demo_casbin/config"
	"demo_casbin/gin_util"
	"demo_casbin/pkg/errors"
	"demo_casbin/schema"
)

// CasbinMiddleware casbin中间件
func CasbinMiddleware(enforcer *casbin.SyncedEnforcer, skippers ...SkipperFunc) gin.HandlerFunc {

	return func(c *gin.Context) {
		if !config.C.Casbin.Enable {
			c.Next()
			return
		}

		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		pIDs := strings.Split(c.GetHeader("project_ids"), ",")
		p := c.Request.URL.Path
		m := c.Request.Method

		var userID string
		var ok bool
		if userID, ok = schema.LoginUser["login"]; !ok {
			gin_util.ResError(c, errors.ErrNoLogin)
			return
		}

		fmt.Println("enforce 1:", userID, pIDs, p, m)
		if len(pIDs) == 0 {
			pIDs = []string{
				"null",
			}
		}

		var pass bool
		for _, pID := range pIDs {
			// TODO: 传入项目ID，获取用户ID
			fmt.Println("enforce 2:", userID, pID, p, m)
			fmt.Printf(enforcer.GetModel().ToText())
			for _, e := range enforcer.GetPolicy() {
				fmt.Println("policy : ", e)
			}
			for _, r := range enforcer.GetAllRoles() {
				fmt.Println("roles : ", r)
			}

			if b, err := enforcer.Enforce(userID, pID, p, m); err != nil {
				gin_util.ResError(c, errors.WithStack(err))
				return
			} else if b {
				pass = true
				break
			}
		}

		if !pass {
			gin_util.ResError(c, errors.ErrNoPerm)
			return
		}

		c.Next()
	}
}
