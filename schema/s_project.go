package schema

import (
	"time"
)

// ProjectQueryResult  查询结果
type ProjectQueryResult struct {
	Data       Projects
	PageResult *PaginationResult
}

// ProjectQueryParam  查询条件
type ProjectQueryParam struct {
	PaginationParam

	Name       string `form:"name"`       // 用户名
	QueryValue string `form:"queryValue"` // 模糊查询
}

// ProjectQueryOptions 查询可选参数项
type ProjectQueryOptions struct {
	OrderFields []*OrderField // 排序字段
}

type Rols []string

// Project 项目
type Project struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Describe string `json:"describe"`

	Ro           Rols         `json:"ro"`
	ProjectUsers ProjectUsers `json:"project_users"`

	UpdateAt time.Time `json:"update_at"`
	CreateAt time.Time `json:"create_at"`
}

// Projects 项目对象列表
type Projects []*Project

// ToIDs 转换为唯一标识
func (a Projects) ToIDs() []string {
	idList := make([]string, len(a))
	for i, item := range a {
		idList[i] = item.ID
	}
	return idList
}

// ToMap 转换为用户ID及信息MAP
func (a Projects) ToMap() map[string]*Project {
	projectMap := make(map[string]*Project, len(a))
	for _, item := range a {
		projectMap[item.ID] = item
	}
	return projectMap
}
