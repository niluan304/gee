package service

import (
	"context"
	"slices"
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

func (s *user) Add(ctx context.Context, req *UserAddReq) (res *UserAddRes, err error) {
	// 插入数据，database 只是一个切片 []Row，用于充当数据库
	database = append(database, Row{req.Name, req.Age, req.Job})
	return
}

func (s *user) Get(ctx context.Context, req *UserGetReq) (res *UserGetRes, err error) {
	// 查询数据，database 只是一个切片 []Row，用于充当数据库
	i := slices.IndexFunc(database, func(row Row) bool { return row.Name == req.Name })
	if i != -1 {
		row := database[i]
		return &UserGetRes{Name: row.Name, Age: row.Age, Job: row.Job}, nil
	}
	return
}

func (s *user) Update(ctx context.Context, req *UserUpdateReq) (res *UserUpdateRes, err error) {
	// 更新数据，database 只是一个切片 []Row，用于充当数据库
	i := slices.IndexFunc(database, func(row Row) bool { return row.Name == req.Name })
	if i != -1 {
		row := database[i] // 用于返回旧数据
		database[i] = Row{Name: req.Name, Age: req.Age, Job: req.Job}
		return &UserUpdateRes{Name: row.Name, Age: row.Age, Job: row.Job}, nil
	}
	return nil, nil
}