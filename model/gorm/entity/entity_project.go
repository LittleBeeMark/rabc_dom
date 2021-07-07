package entity

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"

	"demo_casbin/schema"
	"demo_casbin/util/structure"
)

// GetUserDB 获取用户存储
func GetProjectDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return GetDBWithModel(ctx, defDB, new(Project))
}

// SchemaProject 项目对象
type SchemaProject schema.Project

// ToUser 转换为用户实体
func (a SchemaProject) ToProject() *Project {
	item := new(Project)
	structure.Copy(a, item)
	return item
}

// Project 项目
type Project struct {
	ID       string `gorm:"column:id;PRIMARY_KEY"`
	Name     string `gorm:"column:name;unique_index:idx_project_name"`
	Describe string `gorm:"column:describe"`

	UpdateAt time.Time `gorm:"column:update_at"`
	CreateAt time.Time `gorm:"column:create_at;default:CURRENT_TIMESTAMP;INDEX:idx_projects_create_at"`
}

// ToSchemaUser 转换为用户对象
func (a Project) ToSchemaProject() *schema.Project {
	item := new(schema.Project)
	structure.Copy(a, item)
	return item
}

// Users 用户实体列表
type Projects []*Project

// ToSchemaUsers 转换为用户对象列表
func (a Projects) ToSchemaProjects() []*schema.Project {
	list := make([]*schema.Project, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaProject()
	}
	return list
}
