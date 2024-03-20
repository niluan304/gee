package controller

import (
	"context"

	"gee/web/day8/internal/service"
)

var User = &user{}

type user struct{}

func (c *user) Add(ctx context.Context, bind func(point any) (err error)) (data any, err error) {
	var req service.UserAddReq
	err = bind(&req) // 请求数据的反序列化
	if err != nil {
		return nil, err
	}
	res, err := service.User.Add(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *user) Get(
	ctx context.Context, // 第一个入参一定为 ctx
	bind func(point any) (err error), // 用于反序列化的闭包
) (
	data any, // 返回的数据
	err error, // 错误处理
) {
	var req service.UserGetReq
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
