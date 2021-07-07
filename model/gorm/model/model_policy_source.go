package model

import (
	"context"

	"github.com/jinzhu/gorm"

	"demo_casbin/model/gorm/entity"
	"demo_casbin/pkg/errors"
	"demo_casbin/schema"
)

// PolicySource 表结构通信结构
type PolicySource struct {
	DB *gorm.DB
}

func (a *PolicySource) getQueryOption(opts ...schema.PolicySourceQueryOptions) schema.PolicySourceQueryOptions {
	var opt schema.PolicySourceQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query 查询数据
func (a *PolicySource) Query(ctx context.Context, params schema.PolicySourceQueryParam, opts ...schema.PolicySourceQueryOptions) (*schema.PolicySourceQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := entity.GetPolicySourceDB(ctx, a.DB)
	if v := params.ModuleCode; v != "" {
		db = db.Where("module_code=?", v)
	}
	if v := params.ActionCode; v != "" {
		db = db.Where("action_code=?", v)
	}
	if v := params.Role; v != "" {
		v = "%" + v + "%"
		db = db.Where("role LIKE ?", v)
	}
	if v := params.RoleCode; v != "" {
		db = db.Where("? = ANY(role_code)", v)
	}

	if v := params.QueryValue; v != "" {
		v = "%" + v + "%"
		db = db.Where("module_name LIKE ? OR action_name LIKE ? OR action_path LIKE ? OR action_method LIKE ?", v, v, v, v)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByDESC))
	db = db.Order(ParseOrder(opt.OrderFields))

	var list entity.PolicySources
	pr, err := WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	qr := &schema.PolicySourceQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaPolicySources(),
	}
	return qr, nil
}

// Get 查询指定数据
func (a *PolicySource) Get(ctx context.Context, id string, opts ...schema.PolicySourceQueryOptions) (*schema.PolicySource, error) {
	var item entity.PolicySource
	ok, err := FindOne(ctx, entity.GetPolicySourceDB(ctx, a.DB).Where("id=?", id), &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaPolicySource(), nil
}

// Create 创建数据
func (a *PolicySource) Create(ctx context.Context, item schema.PolicySource) error {
	sitem := entity.SchemaPolicySource(item)
	result := entity.GetPolicySourceDB(ctx, a.DB).Create(sitem.ToPolicySource())
	return errors.WithStack(result.Error)
}

// Update 更新数据
func (a *PolicySource) Update(ctx context.Context, id string, item schema.PolicySource) error {
	eitem := entity.SchemaPolicySource(item).ToPolicySource()
	result := entity.GetPolicySourceDB(ctx, a.DB).Where("id=?", id).Updates(eitem)
	return errors.WithStack(result.Error)
}

// Delete 删除数据
func (a *PolicySource) Delete(ctx context.Context, id string) error {
	result := entity.GetPolicySourceDB(ctx, a.DB).Where("id=?", id).Delete(entity.PolicySource{})
	return errors.WithStack(result.Error)
}
