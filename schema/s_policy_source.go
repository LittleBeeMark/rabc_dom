package schema

import (
	"time"

	"github.com/lib/pq"
)

// PolicySourceQueryResult  查询结果
type PolicySourceQueryResult struct {
	Data       PolicySources
	PageResult *PaginationResult
}

// PolicySourcetQueryParam  查询条件
type PolicySourceQueryParam struct {
	PaginationParam

	Role       string `form:"role"`       // 用户名
	RoleCode   string `form:"roleCode"`   // 用户名
	ModuleCode string `form:"moduleCode"` // 模块吗
	ActionCode string `form:"actionCode"` // 行为码
	QueryValue string `form:"queryValue"` // 模糊查询
}

// PolicySourceQueryOptions 查询可选参数项
type PolicySourceQueryOptions struct {
	OrderFields []*OrderField // 排序字段
}

// PolicySource 资源
type PolicySource struct {
	ID         string `json:"id"`
	ModuleCode string `yaml:"module_code" json:"module_code"` // 编码 （唯一）
	ModuleName string `yaml:"module_name" json:"module_name"` //模块名

	Status  string `yaml:"status" json:"status"`   //(启用/禁用）
	Explain string `yaml:"explain" json:"explain"` // 说明

	ActionCode   string `yaml:"action_code" json:"action_code"`     //行为名 （唯一）
	ActionName   string `yaml:"action_name" json:"action_name"`     //行为名 （唯一）
	ActionPath   string `yaml:"action_path" json:"action_path"`     //    行为路径
	ActionMethod string `yaml:"action_method" json:"action_method"` //行为方法

	Role     string         `yaml:"role" json:"role"`           //  可用角色 (管理员，审计员）
	RoleCode pq.StringArray `yaml:"role_code" json:"role_code"` // 可用角色编码

	CreateAt time.Time `yaml:"create_at" json:"create_at"` // 创建时间
	UpdateAt time.Time `yaml:"update_at" json:"update_at"` // 更新时间
}

// PolicySources 项目对象列表
type PolicySources []*PolicySource

// ToIDs 转换为唯一标识
func (a PolicySources) ToIDs() []string {
	idList := make([]string, len(a))
	for i, item := range a {
		idList[i] = item.ID
	}
	return idList
}

// ToMap 转换为用户ID及信息MAP
func (a PolicySources) ToMap() map[string]*PolicySource {
	policySourceMap := make(map[string]*PolicySource, len(a))
	for _, item := range a {
		policySourceMap[item.ID] = item
	}
	return policySourceMap
}
