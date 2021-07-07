package model

import (
	"context"

	"github.com/jinzhu/gorm"

	"demo_casbin/model/gorm/entity"
	"demo_casbin/pkg/errors"
	"demo_casbin/schema"
)

// Project 表结构通信结构
type Project struct {
	DB *gorm.DB
}

func (a *Project) getQueryOption(opts ...schema.ProjectQueryOptions) schema.ProjectQueryOptions {
	var opt schema.ProjectQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query 查询数据
func (a *Project) Query(ctx context.Context, params schema.ProjectQueryParam, opts ...schema.ProjectQueryOptions) (*schema.ProjectQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := entity.GetProjectDB(ctx, a.DB)
	if v := params.Name; v != "" {
		db = db.Where("name=?", v)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByDESC))
	db = db.Order(ParseOrder(opt.OrderFields))

	var list entity.Projects
	pr, err := WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	qr := &schema.ProjectQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaProjects(),
	}
	return qr, nil
}

// Get 查询指定数据
func (a *Project) Get(ctx context.Context, id string, opts ...schema.ProjectQueryOptions) (*schema.Project, error) {
	var item entity.Project
	ok, err := FindOne(ctx, entity.GetProjectDB(ctx, a.DB).Where("id=?", id), &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaProject(), nil
}

// Create 创建数据
func (a *Project) Create(ctx context.Context, item *schema.Project) error {
	sitem := entity.SchemaProject(*item)
	project := sitem.ToProject()
	result := entity.GetProjectDB(ctx, a.DB).Create(project)
	item.ID = project.ID
	return errors.WithStack(result.Error)
}

// Update 更新数据
func (a *Project) Update(ctx context.Context, id string, item schema.Project) error {
	eitem := entity.SchemaProject(item).ToProject()
	result := entity.GetProjectDB(ctx, a.DB).Where("id=?", id).Updates(eitem)
	return errors.WithStack(result.Error)
}

// Delete 删除数据
func (a *Project) Delete(ctx context.Context, id string) error {
	result := entity.GetProjectDB(ctx, a.DB).Where("id=?", id).Delete(entity.Project{})
	return errors.WithStack(result.Error)
}
