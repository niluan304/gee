package service

import (
	"context"
	"fmt"
	"slices"

	"gee/web/day9/db"
)

type user struct{}

var User = &user{}

func (s *user) Get(ctx context.Context, req *UserGetReq) (res *UserGetRes, err error) {
	// 查询数据
	i := slices.IndexFunc(db.Users, func(row db.User) bool { return row.Id == req.Id })
	if i == -1 { // 数据库未找到数据
		return nil, fmt.Errorf("user not found: %d", req.Id)
	}

	// 返回数据库内容
	row := db.Users[i] // Users 只是一个切片 []User，用于充当数据库
	return &UserGetRes{Id: row.Id, Name: row.Name, TeamId: row.TeamId}, nil
}

type team struct{}

var Team = &team{}

func (s *team) Get(ctx context.Context, req *TeamGetReq) (res *TeamGetRes, err error) {
	// 查询数据
	i := slices.IndexFunc(db.Teams, func(row db.Team) bool { return row.Id == req.Id })
	if i == -1 { // 数据库未找到数据
		return nil, fmt.Errorf("team not found: %d", req.Id)
	}

	// 返回数据库内容
	row := db.Teams[i] // Teams 只是一个切片 []Team，用于充当数据库
	return &TeamGetRes{Id: row.Id, Name: row.Name}, nil
}
