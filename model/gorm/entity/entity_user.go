package entity

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"

	"demo_casbin/schema"
	"demo_casbin/util/structure"
)

// GetUserDB 获取用户存储
func GetUserDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return GetDBWithModel(ctx, defDB, new(User))
}

// SchemaUser 用户对象
type SchemaUser schema.User

// ToUser 转换为用户实体
func (a SchemaUser) ToUser() *User {
	item := new(User)
	structure.Copy(a, item)
	return item
}

// User 用户
type User struct {
	ID string `gorm:"column:id;PRIMARY_KEY"`
	// 名字
	UserName string `gorm:"column:user_name;unique_index:idx_user_name"`
	Password string `gorm:"column:password"`

	// 邮箱
	Email string `gorm:"column:email"`

	// 角色（管理员/操作员/审计员）
	Role string `gorm:"column:role;index:idx_user_role"`

	UpdateAt time.Time `gorm:"column:update_at"`
	CreateAt time.Time `gorm:"column:create_at;default:CURRENT_TIMESTAMP;INDEX:idx_users_create_at"`
}

// ToSchemaUser 转换为用户对象
func (a User) ToSchemaUser() *schema.User {
	item := new(schema.User)
	structure.Copy(a, item)
	return item
}

// Users 用户实体列表
type Users []*User

// ToSchemaUsers 转换为用户对象列表
func (a Users) ToSchemaUsers() []*schema.User {
	list := make([]*schema.User, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaUser()
	}
	return list
}
