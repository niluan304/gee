package handle

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int    `json:"code"` // 业务代码，200 表示 OK，其他表示错误
	Msg  string `json:"msg"`  // 错误消息
	Data any    `json:"data"` // 返回的数据
}

type BinderFunc = func(
	ctx context.Context, // 第一个参数：ctx
	bind func(point any) (err error), // 第二个参数：用于反序列化的闭包
) (
	data any, // 返回的数据
	err error, // 错误处理
)

func Handle(binder BinderFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := binder(c, func(point any) (err error) {
			// 实现反序列化
			return c.ShouldBind(point)
		})
		if err != nil {
			c.JSON(http.StatusOK, Response{Code: 400, Msg: err.Error(), Data: nil})
			return
		}

		c.JSON(http.StatusOK, Response{200, "", data})
	}
}
