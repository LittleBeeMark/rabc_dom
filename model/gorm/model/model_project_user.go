package model

import (
	"context"

	"github.com/jinzhu/gorm"

	"demo_casbin/model/gorm/entity"
	"demo_casbin/pkg/errors"
	"demo_casbin/schema"
)

// ProjectUser 项目相关人员
type ProjectUser struct {
	DB *gorm.DB
}

// Get 查询指定数据
func (a *ProjectUser) Get(ctx context.Context, id string) (*schema.ProjectUser, error) {
	var item entity.ProjectUser
	ok, err := FindOne(ctx, entity.GetProjectUserDB(ctx, a.DB).Where("id=?", id), &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaProjectUser(), nil
}

// Create 创建数据
func (a *ProjectUser) Create(ctx context.Context, item schema.ProjectUser) error {
	sitem := entity.SchemaProjectUser(item)
	result := entity.GetProjectDB(ctx, a.DB).Create(sitem.ToProjectUser())
	return errors.WithStack(result.Error)
}

// Update 更新数据
func (a *ProjectUser) Update(ctx context.Context, id string, item schema.ProjectUser) error {
	eitem := entity.SchemaProjectUser(item).ToProjectUser()
	result := entity.GetProjectUserDB(ctx, a.DB).Where("id=?", id).Updates(eitem)
	return errors.WithStack(result.Error)
}

// Delete 删除数据
func (a *ProjectUser) Delete(ctx context.Context, id int) error {
	result := entity.GetProjectUserDB(ctx, a.DB).Where("id=?", id).Delete(entity.ProjectUser{})
	return errors.WithStack(result.Error)
}

// DeleteByProjectID 根据用户ID删除数据
func (a *ProjectUser) DeleteByProjectID(ctx context.Context, projectID string) error {
	result := entity.GetProjectUserDB(ctx, a.DB).Where("project_id=?", projectID).Delete(entity.ProjectUser{})
	return errors.WithStack(result.Error)
}

// GetProjectUserByProject doc
func (a *ProjectUser) GetProjectUserByProject(ctx context.Context, projectID string) ([]*schema.ProjectUser, error) {
	var item entity.ProjectUsers
	err := entity.GetProjectUserDB(ctx, a.DB).Where("project_id=?", projectID).Find(&item).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return item.ToSchemaProjectUsers(), nil
}

// GetProjectUserByUser doc
func (a *ProjectUser) GetProjectUserByUser(ctx context.Context, userID string) (schema.ProjectUsers, error) {
	var item entity.ProjectUsers
	err := entity.GetProjectUserDB(ctx, a.DB).Where("user_id=?", userID).Find(&item).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return item.ToSchemaProjectUsers(), nil
}
