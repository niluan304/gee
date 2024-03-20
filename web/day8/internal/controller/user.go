package controller

import (
	"net/http"

	"gee/web/day8/internal/model"
	"gee/web/day8/internal/service"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int    `json:"code"`           // 业务代码，200 表示 OK，其他表示错误
	Msg  string `json:"msg"`            // 错误消息
	Data any    `json:"data,omitempty"` // 返回的数据
}

var User = &user{}

type user struct{}

func (c *user) Add(ctx *gin.Context) {
	// 请求数据的反序列化
	var req model.UserAddReq
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, Response{400, err.Error(), nil})
		return
	}
	res, err := service.User.Add(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusOK, Response{400, err.Error(), nil})
		return
	}

	// 填写响应内容
	ctx.JSON(http.StatusOK, Response{200, "", res})
	return
}

func (c *user) Get(ctx *gin.Context) {
	var req model.UserGetReq
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, Response{400, err.Error(), nil})
		return
	}

	res, err := service.User.Get(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusOK, Response{400, err.Error(), nil})
		return
	}

	// 填写响应内容
	ctx.JSON(http.StatusOK, Response{200, "", res})
	return
}
