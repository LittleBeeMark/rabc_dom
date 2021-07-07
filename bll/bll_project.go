package bll

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v2"

	"demo_casbin/model/gorm/model"
	"demo_casbin/pkg/errors"
	"demo_casbin/schema"
)

// Project 用户管理
type Project struct {
	Enforcer         *casbin.SyncedEnforcer
	TransModel       *model.Trans
	ProjectModel     *model.Project
	ProjectUserModel *model.ProjectUser
}

func (a *Project) checkProjectName(ctx context.Context, item schema.Project) error {

	result, err := a.ProjectModel.Query(ctx, schema.ProjectQueryParam{
		PaginationParam: schema.PaginationParam{OnlyCount: true},
		Name:            item.Name,
	})
	if err != nil {
		return err
	} else if result.PageResult.Total > 0 {
		return errors.New400Response("项目名已经存在")
	}
	return nil
}

// Query 查询数据
func (a *Project) Query(ctx context.Context, params schema.ProjectQueryParam, opts ...schema.ProjectQueryOptions) (*schema.ProjectQueryResult, error) {
	return a.ProjectModel.Query(ctx, params, opts...)
}

// Get 查询指定数据
func (a *Project) Get(ctx context.Context, id string, opts ...schema.ProjectQueryOptions) (*schema.Project, error) {
	item, err := a.ProjectModel.Get(ctx, id, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}
	return item, nil
}

// Create 创建数据
func (a *Project) Create(ctx context.Context, item schema.Project) (*schema.IDResult, error) {
	for _, p := range item.ProjectUsers {
		fmt.Printf("userproject: %#v", p)
	}
	err := a.TransModel.Exec(ctx, func(ctx context.Context) error {
		fmt.Println("item. 1", item.ID)
		err := a.ProjectModel.Create(ctx, &item)
		if err != nil {
			return err
		}
		fmt.Println("item. 2", item.ID)

		for _, puItem := range item.ProjectUsers {
			fmt.Println("item. 3", item.ID)
			puItem.ProjectID = item.ID
			err = a.ProjectUserModel.Create(ctx, *puItem)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	LoadCasbinPolicy(ctx, a.Enforcer)
	return schema.NewIDResult(item.ID), nil
}

// Update 更新数据
func (a *Project) Update(ctx context.Context, id string, item schema.Project) error {
	oldItem, err := a.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	} else if oldItem.Name != item.Name {
		err := a.checkProjectName(ctx, item)
		if err != nil {
			return err
		}
	}

	result, err := a.ProjectUserModel.GetProjectUserByProject(ctx, id)
	if err != nil {
		return err
	}

	for _, d := range result {
		oldItem.ProjectUsers = append(oldItem.ProjectUsers, d)
	}

	fmt.Printf("old Item : %#v", item)
	for _, op := range oldItem.ProjectUsers {
		fmt.Printf("old op : %#v", op)
	}

	item.ID = oldItem.ID
	item.CreateAt = oldItem.CreateAt
	err = a.TransModel.Exec(ctx, func(ctx context.Context) error {
		addProjectUsers, delProjectUsers := a.compareProjectUsers(ctx, oldItem.ProjectUsers, item.ProjectUsers)
		for _, rmitem := range addProjectUsers {
			err := a.ProjectUserModel.Create(ctx, *rmitem)
			if err != nil {
				return err
			}
		}

		for _, rmitem := range delProjectUsers {
			err := a.ProjectUserModel.Delete(ctx, rmitem.ID)
			if err != nil {
				return err
			}
		}

		return a.ProjectModel.Update(ctx, id, item)
	})
	if err != nil {
		return err
	}

	LoadCasbinPolicy(ctx, a.Enforcer)
	return nil
}

func (a *Project) compareProjectUsers(ctx context.Context, oldProjectUsers, newProjectUsers schema.ProjectUsers) (addList, delList schema.ProjectUsers) {
	mOldProjectUsers := oldProjectUsers.ToMap()
	mNewProjectUsers := newProjectUsers.ToMap()

	for k, item := range mNewProjectUsers {
		if _, ok := mOldProjectUsers[k]; ok {
			delete(mOldProjectUsers, k)
			continue
		}
		addList = append(addList, item)
	}

	for _, item := range mOldProjectUsers {
		delList = append(delList, item)
	}
	return
}

// Delete 删除数据
func (a *Project) Delete(ctx context.Context, id string) error {
	oldItem, err := a.ProjectModel.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	err = a.TransModel.Exec(ctx, func(ctx context.Context) error {
		err := a.ProjectUserModel.DeleteByProjectID(ctx, id)
		if err != nil {
			return err
		}

		return a.ProjectModel.Delete(ctx, id)
	})
	if err != nil {
		return err
	}

	LoadCasbinPolicy(ctx, a.Enforcer)
	return nil
}
