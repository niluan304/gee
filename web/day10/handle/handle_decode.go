package handle

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DecodeFunc func(
	ctx context.Context, // 第一个参数：ctx
	decode func(point any) (err error), // 第二个参数：用于反序列化的闭包
) (
	data any, // 返回的数据
	err error, // 错误处理
)

func (f DecodeFunc) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := f(c, func(point any) error {
			// 解析动态路由
			if len(c.Params) > 0 {
				err := c.ShouldBindUri(point)
				if err != nil {
					return err
				}
			}

			// 实现反序列化
			return c.ShouldBind(point)
		})
		if err != nil {
			c.JSON(http.StatusOK, Response{Code: CodeBadRequest, Msg: err.Error(), Data: nil})
			return
		}

		c.JSON(http.StatusOK, Response{Code: CodeOK, Msg: "", Data: data})
	}
}
