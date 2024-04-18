package controller

import (
	"context"

	"gee/web/day9/internal/service"
)

type user struct{}

var User = &user{}

func (c *user) Get(ctx context.Context, decode func(point any) (err error)) (data any, err error) {
	var req *struct {
		Id int `uri:"id"`
	}
	err = decode(&req) // 通过闭包反序列化 req
	if err != nil {
		return nil, err
	}

	res, err := service.User.Get(ctx, &service.UserGetReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *user) GetWithTeam(ctx context.Context, decode func(point any) (err error)) (data any, err error) {
	var req *struct {
		Id int `uri:"id"`
	}
	err = decode(&req) // 通过闭包反序列化 req
	if err != nil {
		return nil, err
	}

	userRes, err := service.User.Get(ctx, &service.UserGetReq{Id: req.Id})
	if err != nil {
		return nil, err
	}

	teamRes, err := service.Team.Get(ctx, &service.TeamGetReq{Id: userRes.TeamId})
	if err != nil {
		return nil, err
	}

	type UserWithTeam struct {
		Id   int                `json:"id"`
		Name string             `json:"name"`
		Team service.TeamGetRes `json:"team"`
	}
	return UserWithTeam{
		Id:   userRes.Id,
		Name: userRes.Name,
		Team: *teamRes,
	}, nil
}

type team struct{}

var Team = &team{}

func (c *team) Get(ctx context.Context, decode func(point any) (err error)) (data any, err error) {
	var req *struct {
		Id int `uri:"id"`
	}
	err = decode(&req) // 通过闭包反序列化 req
	if err != nil {
		return nil, err
	}

	res, err := service.Team.Get(ctx, &service.TeamGetReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return res, nil
}
