package bll

import (
	"context"
	"os"
	"strings"

	"github.com/lib/pq"

	"demo_casbin/model/gorm/model"
	"demo_casbin/schema"
	"demo_casbin/util/yaml"
)

type PolicySource struct {
	PolicySourceModel model.PolicySource
}

// InitData 初始化菜单数据
func (a *PolicySource) InitData(ctx context.Context, dataFile string) error {
	data, err := a.readData(dataFile)
	if err != nil {
		return err
	}

	return a.createMenus(ctx, data)
}

func (a *PolicySource) readData(name string) ([]*schema.PolicySource, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data []*schema.PolicySource
	d := yaml.NewDecoder(file)
	d.SetStrict(true)
	err = d.Decode(&data)
	return data, err
}

func ToRoleCode(role string) pq.StringArray {
	roles := strings.Split(role, ",")
	if len(roles) == 0 {
		return []string{}
	}

	var roleCodes []string
	for _, r := range roles {
		if v, ok := schema.RoleListData.ToMap()[r]; ok {
			roleCodes = append(roleCodes, string(v))
		}
	}

	return roleCodes
}

func (a *PolicySource) createMenus(ctx context.Context, list []*schema.PolicySource) error {

	for _, item := range list {
		sitem := schema.PolicySource{
			ModuleCode:   item.ModuleCode,
			ModuleName:   item.ModuleName,
			ActionCode:   item.ActionCode,
			ActionName:   item.ActionName,
			ActionMethod: item.ActionMethod,
			ActionPath:   item.ActionPath,
			Status:       item.Status,
			RoleCode:     ToRoleCode(item.Role),
			Role:         item.Role,
			Explain:      item.Explain,
		}

		_, _, err := a.Create(ctx, sitem)
		if err != nil {
			return err
		}

	}

	return nil
}

func (a *PolicySource) checkCode(ctx context.Context, item schema.PolicySource) (string, bool, error) {
	var itemID string
	result, err := a.PolicySourceModel.Query(ctx, schema.PolicySourceQueryParam{
		PaginationParam: schema.PaginationParam{
			OnlyCount:  false,
			Pagination: false,
		},
		ModuleCode: item.ModuleCode,
		ActionCode: item.ActionCode,
	})
	if err != nil {
		return itemID, false, err
	} else if len(result.Data) > 0 {
		itemID = result.Data[0].ID
		return itemID, true, nil
	}
	return itemID, false, nil
}

// Create 创建数据
func (a *PolicySource) Create(ctx context.Context, item schema.PolicySource) (*schema.IDResult, bool, error) {
	itemID, exist, err := a.checkCode(ctx, item)
	if err != nil {
		return nil, false, err
	}
	if exist {
		err = a.PolicySourceModel.Update(ctx, itemID, item)
		if err != nil {
			return nil, exist, err
		}

		return nil, exist, err
	}

	err = a.PolicySourceModel.Create(ctx, item)
	if err != nil {
		return nil, false, err
	}
	return schema.NewIDResult(item.ID), false, nil
}
