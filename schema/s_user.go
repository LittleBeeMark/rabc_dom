package schema

import (
	"time"
)

type Role string

// 目前支持的用户权限
const (
	RoleAdmin          Role = "admin"          // 管理员
	RoleOperatorLeader Role = "project_leader" // 操作员Leader
	RoleOperator       Role = "operator"       // 操作员
	RoleAuditor        Role = "auditor"        // 审计员

)

func (r Role) String() string {
	switch r {
	case RoleAdmin:
		return "管理员"
	case RoleOperatorLeader:
		return "操作员Leader"
	case RoleOperator:
		return "操作员"
	case RoleAuditor:
		return "审计员"
	default:
		return ""
	}
}

func (r Role) ToCode() string {
	return string(r)
}

func (rs RoleList) ToMap() map[string]Role {
	roleListMap := make(map[string]Role, len(rs))
	for _, r := range rs {
		roleListMap[r.String()] = r
	}

	return roleListMap
}

type RoleList []Role

var RoleListData = RoleList{
	RoleAdmin,
	RoleOperatorLeader,
	RoleOperator,
	RoleAuditor,
}

var LoginUser = map[string]string{}

// 返回固定角色列表
func GetRoleList() RoleList {
	return RoleListData
}

// UserQueryResult 查询结果
type UserQueryResult struct {
	Data       Users
	PageResult *PaginationResult
}

// UserQueryParam 查询条件
type UserQueryParam struct {
	PaginationParam

	UserName   string `form:"userName"`   // 用户名
	QueryValue string `form:"queryValue"` // 模糊查询
}

// UserQueryOptions 查询可选参数项
type UserQueryOptions struct {
	OrderFields []*OrderField // 排序字段
}

// CleanSecure 清理安全数据
func (a *User) CleanSecure() *User {
	a.Password = ""
	return a
}

// User 用户
type User struct {
	ID string `json:"id"`
	// 名字
	UserName string `json:"user_name"`
	Password string `json:"password"`

	// 邮箱
	Email string `json:"email"`

	// 角色（管理员/操作员/审计员）
	Role string `json:"role"`

	UpdateAt time.Time `json:"update_at"`
	CreateAt time.Time `json:"create_at"`
}

// Users 用户对象列表
type Users []*User

// ToIDs 转换为唯一标识列表
func (a Users) ToIDs() []string {
	idList := make([]string, len(a))
	for i, item := range a {
		idList[i] = item.ID
	}
	return idList
}

// ToMap 转换为用户ID及信息MAP
func (a Users) ToMap() map[string]*User {
	userMap := make(map[string]*User, len(a))
	for _, item := range a {
		userMap[item.ID] = item
	}
	return userMap
}

// ProjectUser 项目相关人员
type ProjectUser struct {
	ID        int       `json:"id"`
	UserID    string    `json:"user_id"`
	ProjectID string    `json:"project_id"`
	CreateAt  time.Time `json:"create_at"`
}

// ProjectUsers 项目用户
type ProjectUsers []*ProjectUser

// 获取所有的 projectID
func (pus ProjectUsers) ToProjectIDs() []string {
	l := len(pus)
	projectIDs := make([]string, l)
	for i, pu := range pus {
		projectIDs[i] = pu.ProjectID
	}

	return projectIDs
}

// ToMap 转换为map
func (pus ProjectUsers) ToMap() map[string]*ProjectUser {
	m := make(map[string]*ProjectUser)
	for _, item := range pus {
		m[item.UserID] = item
	}
	return m
}
