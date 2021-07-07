package entity

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"

	"demo_casbin/schema"
	"demo_casbin/util/structure"
)

// GetProjectUserDB 获取用户存储
func GetProjectUserDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return GetDBWithModel(ctx, defDB, new(ProjectUser))
}

// SchemaUser 用户对象
type SchemaProjectUser schema.ProjectUser

// ToUser 转换为用户实体
func (a SchemaProjectUser) ToProjectUser() *ProjectUser {
	item := new(ProjectUser)
	structure.Copy(a, item)
	return item
}

// ProjectUser 项目相关人员
type ProjectUser struct {
	ID        int       `gorm:"column:id;type:SERIAL;PRIMARY_KEY"`
	UserID    string    `gorm:"column:user_id;unique_index:idx_project_user"`
	ProjectID string    `gorm:"column:project_id;unique_index:idx_project_user"`
	CreateAt  time.Time `gorm:"column:create_at;default:CURRENT_TIMESTAMP"`
}

// ToSchemaProjectUser 换为用户对象
func (a ProjectUser) ToSchemaProjectUser() *schema.ProjectUser {
	item := new(schema.ProjectUser)
	structure.Copy(a, item)
	return item
}

// Users 用户实体列表
type ProjectUsers []*ProjectUser

// ToSchemaUsers 转换为用户对象列表
func (a ProjectUsers) ToSchemaProjectUsers() schema.ProjectUsers {
	list := make([]*schema.ProjectUser, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaProjectUser()
	}
	return list
}
