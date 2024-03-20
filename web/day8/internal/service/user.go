package service

import (
	"context"
	"slices"

	"gee/web/day8/internal/model"
)

type Row struct {
	Name string
	Age  int
	Job  string
}

// database 只是一个切片 []Row，用于充当数据库
var database = []Row{
	{"Alice", 32, "teacher"},
	{"Bob", 30, "driver"},
}

var User = &user{}

type user struct{}

func (s *user) Add(ctx context.Context, req model.UserAddReq) (res *model.UserAddRes, err error) {
	// 插入数据
	database = append(database, Row{req.Name, req.Age, req.Job})
	return
}

func (s *user) Get(ctx context.Context, req model.UserGetReq) (res *model.UserGetRes, err error) {
	// 查询数据
	i := slices.IndexFunc(database, func(row Row) bool { return row.Name == req.Name })
	if i != -1 {
		// 填写响应内容
		row := database[i]
		return &model.UserGetRes{Name: row.Name, Age: row.Age, Job: row.Job}, nil
	}
	return
}
