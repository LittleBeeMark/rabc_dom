package casbin

import (
	"context"
	"fmt"

	casbinModel "github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"

	"demo_casbin/model/gorm/model"
	"demo_casbin/pkg/logger"
	"demo_casbin/schema"
)

var _ persist.Adapter = (*CasbinAdapter)(nil)

// CasbinAdapter casbin适配器
type CasbinAdapter struct {
	UserModel         *model.User
	ProjectModel      *model.Project
	ProjectUserModel  *model.ProjectUser
	PolicySourceModel *model.PolicySource
}

// LoadPolicy loads all policy rules from the storage.
func (a *CasbinAdapter) LoadPolicy(model casbinModel.Model) error {
	ctx := context.Background()
	err := a.loadRolePolicy(ctx, model)
	if err != nil {
		logger.WithContext(ctx).Errorf("Load casbin role policy error: %s", err.Error())
		return err
	}

	err = a.loadUserPolicy(ctx, model)
	if err != nil {
		logger.WithContext(ctx).Errorf("Load casbin user policy error: %s", err.Error())
		return err
	}

	return nil
}

// 加载角色策略(p,role_id,path,method)
func (a *CasbinAdapter) loadRolePolicy(ctx context.Context, m casbinModel.Model) error {
	// 获取所有的项目ID
	projectResults, err := a.ProjectModel.Query(ctx, schema.ProjectQueryParam{}, schema.ProjectQueryOptions{})
	if err != nil {
		return err
	}
	pjIDs := projectResults.Data.ToIDs()

	for _, rl := range schema.RoleListData {
		// 获取所有的请求资源
		policySources, err := a.PolicySourceModel.Query(ctx, schema.PolicySourceQueryParam{RoleCode: string(rl)})
		if err != nil {
			return err
		}

		mcache := make(map[string]struct{})
		for _, ps := range policySources.Data {
			for _, pjID := range pjIDs {
				if ps.ActionPath == "" || ps.ActionMethod == "" {
					continue
				} else if _, ok := mcache[pjID+ps.ActionPath+ps.ActionMethod]; ok {
					continue
				}

				mcache[pjID+ps.ActionPath+ps.ActionMethod] = struct{}{}
				line := fmt.Sprintf("p,%s,%s,%s,%s", rl.ToCode(), pjID, ps.ActionPath, ps.ActionMethod)
				persist.LoadPolicyLine(line, m)
			}
		}
	}

	return nil
}

// 加载用户策略(g,user_id,role_id)
func (a *CasbinAdapter) loadUserPolicy(ctx context.Context, m casbinModel.Model) error {
	// 获取所有的用户信息
	userResult, err := a.UserModel.Query(ctx, schema.UserQueryParam{})
	if err != nil {
		return err
	} else if len(userResult.Data) > 0 {
		for _, uitem := range userResult.Data {

			var pIDs []string
			if uitem.Role == schema.RoleAdmin.ToCode() {
				projectResults, err := a.ProjectModel.Query(ctx, schema.ProjectQueryParam{}, schema.ProjectQueryOptions{})
				if err != nil {
					return err
				}
				pIDs = projectResults.Data.ToIDs()
			} else {
				var projectUsers schema.ProjectUsers
				projectUsers, err = a.ProjectUserModel.GetProjectUserByUser(ctx, uitem.ID)
				if err != nil {
					return err
				}
				pIDs = projectUsers.ToProjectIDs()
			}

			for _, pID := range pIDs {
				line := fmt.Sprintf("g,%s,%s,%s", uitem.ID, uitem.Role, pID)
				persist.LoadPolicyLine(line, m)

			}
		}
	}

	return nil
}

// SavePolicy saves all policy rules to the storage.
func (a *CasbinAdapter) SavePolicy(model casbinModel.Model) error {
	return nil
}

// AddPolicy adds a policy rule to the storage.
// This is part of the Auto-Save feature.
func (a *CasbinAdapter) AddPolicy(sec string, ptype string, rule []string) error {
	return nil
}

// RemovePolicy removes a policy rule from the storage.
// This is part of the Auto-Save feature.
func (a *CasbinAdapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return nil
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
// This is part of the Auto-Save feature.
func (a *CasbinAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return nil
}
