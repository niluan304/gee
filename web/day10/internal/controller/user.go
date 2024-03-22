package controller

import (
	"context"

	"gee/web/day8/internal/service"
)

var User = &user{}

type user struct{}

func (c *user) Add(ctx context.Context, bind func(point any) (err error)) (data any, err error) {
	var req *service.UserAddReq
	err = bind(&req) // 通过闭包反序列化 req
	if err != nil {
		return nil, err
	}
	res, err := service.User.Add(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *user) Get(ctx context.Context, bind func(point any) (err error)) (data any, err error) {
	var req *service.UserGetReq
	err = bind(&req) // 通过闭包反序列化 req
	if err != nil {
		return nil, err
	}

	res, err := service.User.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *user) Upsert(ctx context.Context, req *UserUpsertReq) (res *UserUpsertRes, err error) {
	// 尝试更新数据
	update, err := service.User.Update(ctx, &service.UserUpdateReq{Name: req.Name, Age: req.Age, Job: req.Job})
	if err != nil {
		return nil, err
	}
	if update != nil {
		return &UserUpsertRes{Name: update.Name, Age: update.Age, Job: update.Job}, nil
	}

	// 数据不存在则新增
	// TODO 更新数据不存在时，应当返回自定义类型的错误，而不是通过 nil 判断
	_, err = service.User.Add(ctx, &service.UserAddReq{Name: req.Name, Age: req.Age, Job: req.Job})
	if err != nil {
		return nil, err
	}
	return nil, nil
}
