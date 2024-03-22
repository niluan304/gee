---
title: "gee-web-day10 通过反射构造规范路由"
description: 第二种路由注册方法：func(context.Context, *XXXReq) (*XXXRes, error)
date: 2023-12-22T15:43:59+08:00
image: 
license: 
hidden: false
comments: true
draft: false
math: false
tags: [gee,web,goframe]
categories: go
---

> 源代码/数据集已上传到：[GitHub - follow gee to learn go](https://github.com/niluan304/gee)

## `GoFrame` 的 `ReqResFunc` 类型
在 day7.5 开篇的时候，笔者提到 `GoFrame` 支持第二种路由注册方法，这里笔者称之为 `ReqResFunc` 类型（下文同）：
```go
// 写法二
func (ctx context.Context, req *{Prefix}Req) (res *{Prefix}Res, err error){
    // 业务代码逻辑
}
```

但是 day9 实现的函数签名：
```go
func (c *user) Get(ctx context.Context, bind func(point any) (err error)) (data any, err error)
```

对比可以发现，和 `ReqResFunc` 类型有明显不同，我们可以在 `GoFrame` 的源码里一探究竟。

从 `GoFrame` 的 [文档「路由注册-函数注册」中](https://goframe.org/pages/viewpage.action?pageId=1114240)，可以找到 [入口函数](https://github.com/gogf/gf/blob/313d9d138f96b0ed460d47684298a7fb26d3fd75/net/ghttp/ghttp_server_service_handler.go#L21-L39)：
```go
// https://github.com/gogf/gf/blob/313d9d138f96b0ed460d47684298a7fb26d3fd75/net/ghttp/ghttp_server_service_handler.go#L21-L39

// BindHandler registers a handler function to server with a given pattern.
//
// Note that the parameter `handler` can be type of:
// 1. func(*ghttp.Request)
// 2. func(context.Context, BizRequest)(BizResponse, error)
func (s *Server) BindHandler(pattern string, handler interface{}) {
	var ctx = context.TODO()
	funcInfo, err := s.checkAndCreateFuncInfo(handler, "", "", "")
	if err != nil {
		s.Logger().Fatalf(ctx, `%+v`, err)
	}
	s.doBindHandler(ctx, doBindHandlerInput{
		Prefix:     "",
		Pattern:    pattern,
		FuncInfo:   funcInfo,
		Middleware: nil,
		Source:     "",
	})
}
```

在源码里，可以发现关键代码在 [`checkAndCreateFuncInfo` 方法](https://github.com/gogf/gf/blob/313d9d138f96b0ed460d47684298a7fb26d3fd75/net/ghttp/ghttp_server_service_handler.go#L148)，继续前行，就能够发现端倪：
```go
// https://github.com/gogf/gf/blob/313d9d138f96b0ed460d47684298a7fb26d3fd75/net/ghttp/ghttp_server_service_handler.go#L148

func (s *Server) checkAndCreateFuncInfo(f interface{}, pkgPath, structName,methodName string,) (funcInfo handlerFuncInfo, err error) {
	funcInfo = handlerFuncInfo{ // 根据传入的 f，初始化返回值
		Type:  reflect.TypeOf(f),
		Value: reflect.ValueOf(f),
	}
}
```

`GoFrame` 通过反射 `reflect`，获取了传入的函数的参数信息，并做了相应的校验，关键代码有 5 行：
```go
// 校验请求和返回的参数数量
if reflectType.NumIn() != 2 || reflectType.NumOut() != 2 

// 第一个请求参数必须为 context.Context 类型
if !reflectType.In(0).Implements(reflect.TypeOf((*context.Context)(nil)).Elem())


// 第二个返回参数必须为 error 类型
if !reflectType.Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem())


// 第二个请求参数必须为以 `Req` 结尾
if !strings.HasSuffix(reflectType.In(1).String(), `Req`)

// 第一个返回参数必须为以 `Res` 结尾
if !strings.HasSuffix(reflectType.Out(0).String(), `Res`) 
```

通过这些校验，`GoFrame` 就实现了规范路由函数必须是 `ResReqFunc` 类型的约束。校验过程中，有一些细节：
1. `ctx` 和 `error` 是接口类型，只能调用 `func (Type) Implements(u Type) bool` 确认是否实现了对应的接口，`(*error)(nil)` 和 `(*context.Context)(nil)` 则是声明了对应接口的空值 `nil`。
2. `req` 和 `res` 的初始类型是结构体，可以直接获取结构体的。


## 实现 `ReqResFunc` 类型的约束
接下来，我们就可以仿照 `GoFrame`，实现 `ResReqFunc` 类型的约束。首先需要创建一个结构体，用于保存反射解析出来的值：
```go
type ReqResFunc struct {
	fn reflect.Value // 函数调用入口

	ctx reflect.Type // 第一个请求参数：context.Context
	req reflect.Type // 第二个请求参数：XXXReq

	res reflect.Type // 第一个返回参数：XXXRes
	err reflect.Type // 第二个返回参数：error
}
```

具体的解析代码，可以全部仿照 `GoFrame` 的流程，获取入参 `reqRes` 的反射对象，然后逐个校验，最后再构造 `ReqResFunc`。

那么还剩最后一个问题， `ReqResFunc` 要注册到 `gin`框架里呢？这里和 day9 遇到的情况一样，想让两个不兼容的接口，能够一起工作，就需要一个中间函数：
```go
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

		c.JSON(http.StatusOK, result[0].Interface())
	}
}
```

而 `GoFrame` 也是这样转化的，相关源码： [`func createRouterFunc(funcInfo handlerFuncInfo) func(r *Request)`](https://github.com/gogf/gf/blob/313d9d138f96b0ed460d47684298a7fb26d3fd75/net/ghttp/ghttp_server_service_handler.go#L264-L306)

## 验证 `ReqResFunc` 类型
增加了一个新特性，那必然是需要测试的，这里可以新建 `Upsert` 接口（如果不存在就新增，存在就更新）用于测试：
```go
func (c *user) Upsert(ctx context.Context, req *UserUpsertReq) (res *UserUpsertRes, err error) {
	// 尝试更新数据
	update, err := service.User.Update(ctx, &service.UserUpdateReq{Name: req.Name, Age: req.Age, Job: req.Job})
	if err != nil {
		return nil, err
	}
	if update != nil {
		return &UserUpsertRes{Name: update.Name, Age: update.Age, Job: update.Job}, nil
	}

	// 数据不存在则新增
	// TODO 更新数据不存在时，应当返回自定义类型的错误，而不是通过 nil 判断
	_, err = service.User.Add(ctx, &service.UserAddReq{Name: req.Name, Age: req.Age, Job: req.Job})
	if err != nil {
		return nil, err
	}
	return nil, nil
}
```

代码实现很简单，先调用 `service.User.Update`，如果 `update == nil` 就表示更新失败，数据库未找到这条数据。然后再执行数据插入 `service.User.Add` 的操作。当然，主流的数据库里都有类似的语法实现 `Upsert`，如 `MySQL` 和 `Postgres`。

但这并不是重点，我们应该聚焦于新接口 `*user.Upsert` 本身，他只调用两个已有的业务方法就完成自身的工作，这解决了 day8 提到的「代码难以复用」的问题：

> 假如随着项目的进展，导致 `Add` 部分业务代码和 `Get` 的完全一致，甚至最好的解决办法是直接调用 `Get` 方式，但这没办法真的去调用 `Get` 方式，因为 `Add` 接收到的 `*gin.Context` 与 `Get` 方法需要的 `*gin.Context` 是有差别的。

测试接口：
```go
func client() {
	time.Sleep(time.Second) // 等待路由注册

	resp1, _ := http.Get("http://localhost:8080/user?name=Carol")
	resp2, _ := http.Get("http://localhost:8080/user?name=Bob")
	resp3, _ := http.Post("http://localhost:8080/user", "application/json", bytes.NewBufferString(`{"name":"Carol","age":44,"job":"worker"}`))
	resp4, _ := http.Get("http://localhost:8080/user?name=Carol")
	resp5, _ := http.Post("http://localhost:8080/user/upsert", "application/json", bytes.NewBufferString(`{"name":"Dave","age":32,"job":"nurse"}`))
	resp6, _ := http.Post("http://localhost:8080/user/upsert", "application/json", bytes.NewBufferString(`{"name":"Dave","age":35,"job":"doctor"}`))
	resp7, _ := http.Get("http://localhost:8080/user?name=Dave")

	for _, resp := range []*http.Response{resp1, resp2, resp3, resp4, resp5, resp6, resp7} {
		data, _ := io.ReadAll(resp.Body)
		fmt.Println(string(data))
	}

	// Output:
	// {"code":200,"msg":"","data":null}
	// {"code":200,"msg":"","data":{"name":"Bob","age":30,"job":"driver"}}
	// {"code":200,"msg":"","data":null}
	// {"code":200,"msg":"","data":{"name":"Carol","age":44,"job":"worker"}}
	// {"code":200,"msg":"","data":null}
	// {"code":200,"msg":"","data":{"name":"Dave","age":32,"job":"nurse"}}
	// {"code":200,"msg":"","data":{"name":"Dave","age":35,"job":"doctor"}}
}
```

从 `resp567` 的结果来看，先插入了 `Dave` 为 `nurse` 的数据，后面将 `Dave` 修改了 `doctor`，符合预期，说明 `ReqResFunc` 类型测试通过。

