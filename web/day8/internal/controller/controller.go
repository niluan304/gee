package controller

import (
	"net/http"

	"gee/web/day8/internal/service"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int    `json:"code"` // 业务代码，200 表示 OK，其他表示错误
	Msg  string `json:"msg"`  // 错误消息
	Data any    `json:"data"` // 返回的数据
}

const (
	CodeOK         = 200 // 业务正常
	CodeBadRequest = 400 // 请求参数异常
)

type user struct{}

var User = &user{}

func (c *user) Get(ctx *gin.Context) {
	var req *service.UserGetReq
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, Response{CodeBadRequest, err.Error(), nil})
		return
	}

	res, err := service.User.Get(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusOK, Response{CodeBadRequest, err.Error(), nil})
		return
	}

	// 填写响应内容
	ctx.JSON(http.StatusOK, Response{CodeOK, "", res})
	return
}

func (c *user) GetWithTeam(ctx *gin.Context) {
	var req *service.UserGetReq
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, Response{CodeBadRequest, err.Error(), nil})
		return
	}

	userRes, err := service.User.Get(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusOK, Response{CodeBadRequest, err.Error(), nil})
		return
	}

	teamRes, err := service.Team.Get(ctx, &service.TeamGetReq{Id: userRes.TeamId})
	if err != nil {
		ctx.JSON(http.StatusOK, Response{CodeBadRequest, err.Error(), nil})
		return
	}

	type UserWithTeam struct {
		Id   int                `json:"id"`
		Name string             `json:"name"`
		Team service.TeamGetRes `json:"team"`
	}
	// 填写响应内容
	ctx.JSON(http.StatusOK, Response{CodeOK, "", UserWithTeam{
		Id:   userRes.Id,
		Name: userRes.Name,
		Team: *teamRes,
	}})
	return
}

type team struct{}

var Team = &team{}

func (c *team) Get(ctx *gin.Context) {
	var req *service.TeamGetReq
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, Response{CodeBadRequest, err.Error(), nil})
		return
	}

	res, err := service.Team.Get(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusOK, Response{CodeBadRequest, err.Error(), nil})
		return
	}

	// 填写响应内容
	ctx.JSON(http.StatusOK, Response{CodeOK, "", res})
	return
}
