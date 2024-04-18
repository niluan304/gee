package handle

import (
	"reflect"

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

func Handle(decode DecodeFunc) gin.HandlerFunc {
	return decode.Handler()
}

// ObjectHandler 通过结构体（对象）注册路由
// 这个结构体的所有方法都必须为 `ReqResFunc` 格式，否则会触发 panic
//
// 缺陷：无法为 handles[i] 绑定 `path` 和 `method`
//
// 已知的解决方法：
//
// 1. 在请求参数 `XXXReq` 里写 tag，
// 参考：[规范参数结构](https://goframe.org/pages/viewpage.action?pageId=116004922)
//
// 2. 要求函数名格式为 请求方法+请求路径，如 `GetHelloWorld` 对应 `GET: /hello/world`，
// 参考：[examples/mvc/hello-world/main.go](https://github.com/iris-contrib/examples/blob/master/mvc/hello-world/main.go)
//
// 实现见 `handle_test.go/TestObjectHandler`
func ObjectHandler(object any, f func(fn *ReqResFunc, methodName string)) {
	v := reflect.ValueOf(object)

	// 如果是结构体, 那么获取这个结构体的指针, 从而遍历到他的所有方法
	if v.Kind() == reflect.Struct {
		newValue := reflect.New(v.Type())
		newValue.Elem().Set(v)
		v = newValue
	}

	if v.Kind() != reflect.Pointer {
		panic("the kind of object must be Struct or *Struct")
	}

	t := v.Type()
	for i := 0; i < t.NumMethod(); i++ {
		fn := NewReqResFunc(v.Method(i).Interface())
		f(fn, t.Method(i).Name)
	}
}
