package handle

import (
	"context"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int    `json:"code"` // 业务代码，200 表示 OK，其他表示错误
	Msg  string `json:"msg"`  // 错误消息
	Data any    `json:"data"` // 返回的数据
}

type BinderFunc func(
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

// ReqResHandle 返回 gin.HandlerFunc
// 参数 reqResFunc 必须是 func(context.Context, *XXXReq) (*XXXRes, error) 格式
func ReqResHandle(reqResFunc any) gin.HandlerFunc {
	f := NewReqResFunc(reqResFunc)
	return func(c *gin.Context) {
		req := reflect.New(f.req).Interface() // 使用 reflect.New 初始化变量

		if err := c.ShouldBindJSON(req); err != nil {
			c.JSON(http.StatusOK, Response{Code: 400, Msg: err.Error(), Data: nil})
			return
		}

		result := f.fn.Call([]reflect.Value{reflect.ValueOf(c), reflect.ValueOf(req).Elem()})
		if err := result[1]; !err.IsNil() {
			c.JSON(http.StatusOK, Response{Code: 400, Msg: err.Interface().(error).Error(), Data: nil})
		}

		c.JSON(http.StatusOK, Response{200, "", result[0].Interface()})
	}
}

type ReqResFunc struct {
	fn reflect.Value // 函数调用入口

	ctx reflect.Type // 第一个请求参数：context.Context
	req reflect.Type // 第二个请求参数：XXXReq

	res reflect.Type // 第一个返回参数：XXXRes
	err reflect.Type // 第二个返回参数：error
}

// NewReqResFunc 返回 ReqResFunc
// 参数 reqResFunc 必须是 func(context.Context, *XXXReq) (*XXXRes, error) 格式
func NewReqResFunc(reqRes any) *ReqResFunc {
	fn := reflect.ValueOf(reqRes)
	fnType := fn.Type()

	if fnType.NumIn() != 2 {
		panic("parameter must be context.Context and XXXReq")
	}
	if fnType.NumOut() != 2 {
		panic("return value must be XXXRes and error")
	}

	ctx, req := fnType.In(0), fnType.In(1)
	res, err := fnType.Out(0), fnType.Out(1)

	if !ctx.Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
		panic("the first parameter must be context.Context")
	}
	if !err.Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		panic("the second return value must be error")
	}

	if !strings.HasSuffix(req.String(), "Req") {
		panic("the name of second parameter must be XXXReq")
	}
	if !strings.HasSuffix(res.String(), "Res") {
		panic("the name of first return value must be XXXRes")
	}

	return &ReqResFunc{
		fn:  fn,
		ctx: ctx,
		req: req,
		res: res,
		err: err,
	}
}
