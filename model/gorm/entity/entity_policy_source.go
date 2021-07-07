package entity

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"

	"demo_casbin/schema"
	"demo_casbin/util/structure"
)

// GetPolicySourceDB 获取资源DB
func GetPolicySourceDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return GetDBWithModel(ctx, defDB, new(PolicySource))
}

// SchemaPolicySource 项目对象
type SchemaPolicySource schema.PolicySource

// ToUser 转换为用户实体
func (a SchemaPolicySource) ToPolicySource() *PolicySource {
	item := new(PolicySource)
	structure.Copy(a, item)
	return item
}

type PolicySource struct {
	ID         string `gorm:"column:id;PRIMARY_KEY"`
	ModuleCode string `gorm:"column:module_code;unique_index:idx_module_action" json:"module_code"` // 编码 （唯一）
	//ModuleName string `gorm:"column:module_name" json:"module_name"`                                //模块名

	Status  string `gorm:"column:status"`  //(启用/禁用）
	Explain string `gorm:"column:explain"` // 说明

	ActionCode string `gorm:"column:action_code;unique_index:idx_module_action" json:"action_code"` //行为名 （唯一）
	//ActionName   string `gorm:"column:action_name" json:"action_name"`                                //行为名 （唯一）
	ActionPath   string `gorm:"column:action_path" json:"action_path"`     //    行为路径
	ActionMethod string `gorm:"column:action_method" json:"action_method"` //行为方法

	RoleCode pq.StringArray `gorm:"column:role_code;type:TEXT[]"` // 可用角色编码

	CreateAt time.Time `gorm:"column:create_at" json:"create_at"` // 创建时间
	UpdateAt time.Time `gorm:"cloumn:update_at" json:"update_at"` // 更新时间
}

// ToSchemaUser 转换为用户对象
func (a PolicySource) ToSchemaPolicySource() *schema.PolicySource {
	item := new(schema.PolicySource)
	structure.Copy(a, item)
	return item
}

// PolicySources 资源实体列表
type PolicySources []*PolicySource

// ToSchemaPolicySources 转换为策略资源列表对象列表
func (a PolicySources) ToSchemaPolicySources() []*schema.PolicySource {
	list := make([]*schema.PolicySource, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaPolicySource()
	}
	return list
}
