package internal

import (
	"fmt"
	"net/http"
	"slices"

	"gee/web/day7.5/db"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int    `json:"code"` // 业务代码，200 表示 OK，其他表示错误
	Msg  string `json:"msg"`  // 错误消息
	Data any    `json:"data"` // 返回的数据
}

const (
	CodeOK         = 200
	CodeBadRequest = 400
)

type user struct{}

var User = &user{}

func (c *user) Get(ctx *gin.Context) {
	var req *struct {
		Id int `uri:"id"`
	}
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, Response{CodeBadRequest, err.Error(), nil})
		return
	}

	i := slices.IndexFunc(db.Users, func(row db.User) bool { return row.Id == req.Id })
	if i == -1 { // 数据库未找到数据
		ctx.JSON(http.StatusOK, Response{CodeBadRequest, fmt.Sprintf("user not found: %d", req.Id), nil})
		return
	}

	// 返回数据库内容
	row := db.Users[i] // Teams 只是一个切片 []Team，用于充当数据库表
	ctx.JSON(http.StatusOK, Response{CodeOK, "", row})
	return
}

type team struct{}

var Team = &team{}

func (c *team) Get(ctx *gin.Context) {
	var req *struct {
		Id int `uri:"id"`
	}
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, Response{CodeBadRequest, err.Error(), nil})
		return
	}

	// 查询数据
	i := slices.IndexFunc(db.Teams, func(row db.Team) bool { return row.Id == req.Id })
	if i == -1 { // 数据库未找到数据
		ctx.JSON(http.StatusOK, Response{CodeBadRequest, fmt.Sprintf("team not found: %d", req.Id), nil})
		return
	}
	// 返回数据库内容
	row := db.Teams[i] // Teams 只是一个切片 []Team，用于充当数据库表
	ctx.JSON(http.StatusOK, Response{CodeOK, "", row})
	return
}
