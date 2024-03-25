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
// 参数 reqResFunc 必须是 func(context.Context, *XXXReq) (*XXXRes, error) 格式，否则会触发 panic
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
// 参数 reqResFunc 必须是 func(context.Context, *XXXReq) (*XXXRes, error) 格式，否则会触发 panic
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

// ObjectHandler 通过结构体（对象）注册路由
// 这个结构体的所有方法都必须为 `ReqResFunc` 格式，否则会触发 panic
//
// TODO 解决缺陷：无法为 handles[i] 绑定 `path` 和 `method`
//
// 已知的解决方法：
//
// 1. 在请求参数 `XXXReq` 里写 tag，
// 参考：[规范参数结构](https://goframe.org/pages/viewpage.action?pageId=116004922)
//
// 2. 要求函数名格式为 请求方法+请求路径，如 `GetHelloWorld` 对应 `GET: /hello/world`，
// 参考：[examples/mvc/hello-world/main.go](https://github.com/iris-contrib/examples/blob/master/mvc/hello-world/main.go)
func ObjectHandler(object any) (handles []gin.HandlerFunc) {
	v := reflect.ValueOf(object)

	// 如果是结构体, 那么获取这个结构体的指针, 从而遍历到他的所有方法
	if v.Kind() == reflect.Struct {
		newValue := reflect.New(v.Type())
		newValue.Elem().Set(v)
		v = newValue
	}

	if v.Kind() != reflect.Pointer {
		panic("v.Kind() must be reflect.Pointer")
	}

	t := v.Type()
	for i := 0; i < t.NumMethod(); i++ {
		fn := v.MethodByName(t.Method(i).Name) // 所有方法都必须为 ReqResFunc 类型
		handles = append(handles, ReqResHandle(fn.Interface()))
	}

	return handles
}
