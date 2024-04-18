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
tags: [gee,web]
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
func (c *user) Get(ctx context.Context, decode func(point any) (err error)) (data any, err error)
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
2. `req` 和 `res` 的初始类型是结构体，可以直接获取结构体的类型名。


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

那么还剩最后一个问题， `ReqResFunc` 要注册到 `Gin`框架里呢？这里和 day9 遇到的情况一样，还是那句名言：
> All problems in computer science can be solved by another level of indirection.
>
> 计算机科学领域的任何问题都可以通过增加一个中间层来解决。


倘若解决不了，那就再加一个中间层，因此完全可以把 `ReqResFunc` 转换为 `DecodeFunc`，实现也不复杂：
```go
func (f *ReqResFunc) Call(ctx context.Context, decode func(point any) error) (any, error) {
	req := reflect.New(f.req.Elem())
	point := req.Interface()

	if err := decode(point); err != nil {
		return nil, err
	}

	result := f.fn.Call([]reflect.Value{reflect.ValueOf(ctx), req})
	if err := result[1]; !err.IsNil() {
		return nil, err.Interface().(error)
	}
	return result[0].Interface(), nil
}

func (f *ReqResFunc) DecodeFunc() DecodeFunc {
	return f.Call
}
```

而 `GoFrame` 也是这样转化的，相关源码： [`func createRouterFunc(funcInfo handlerFuncInfo) func(r *Request)`](https://github.com/gogf/gf/blob/313d9d138f96b0ed460d47684298a7fb26d3fd75/net/ghttp/ghttp_server_service_handler.go#L264-L306)

## 验证 `ReqResFunc` 类型
增加了一个新特性，做个简单测试，通过 `teamId` 获取团队成员：
```go
// internal/controller/controller.go
func (c *team) GetUsers(ctx context.Context, req *TeamGetUsersReq) (res *TeamGetUsersRes, err error) {
	out, err := service.Team.GetUsers(ctx, &service.TeamGetUsersReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &TeamGetUsersRes{TeamGetUsersRes: out}, nil
}

// internal/service/service.go
func (s *team) GetUsers(ctx context.Context, req *TeamGetUsersReq) (res *TeamGetUsersRes, err error) {
	var users []UserGetRes

	// 查询数据
	// Users 只是一个切片 []User，用于充当数据库
	for _, row := range db.Users {
		if row.TeamId == req.Id {
			users = append(users, UserGetRes{Id: row.Id, Name: row.Name, TeamId: row.TeamId})
		}
	}
	return &TeamGetUsersRes{Users: users}, nil
}
```

测试接口：
```go
func Test_Client(t *testing.T) {
	time.Sleep(time.Second) // 等待服务端启动

	paths := []string{
		"/user/1",       // {"code":200,"msg":"","data":{"id":1,"name":"Alice","teamId":1}}
		"/user/3",       // {"code":400,"msg":"user not found: 3","data":null}
		"/user/1/team",  // {"code":200,"msg":"","data":{"id":1,"name":"Alice","team":{"id":3,"name":"Apple"}}}
		"/team/3",       // {"code":200,"msg":"","data":{"id":3,"name":"Apple"}}
		"/team/5",       // {"code":400,"msg":"team not found: 5","data":null}
		"/team/3/users", // {"code":200,"msg":"","data":{"users":[{"id":1,"name":"Alice","teamId":3}]}}
		"/team/5/users", // {"code":200,"msg":"","data":{"Users":null}}
	}

	for _, path := range paths {
		resp, err := http.Get("http://localhost:8080" + path)
		if err != nil {
			fmt.Println("req err", err)
			continue
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("read resp.Body err", err)
			continue
		}
		fmt.Println(string(data))
	}
}
```

## 通过对象注册路由
事实上，`GoFrame` 还有第三种路由注册方法：[对象注册](https://goframe.org/pages/viewpage.action?pageId=116004922)，向 `(*ghttp.RouterGroup).Bind` 传入一个结构体变量，然后 `GoFrame` 会尝试注册这个结构体上的所有 `ReqResFunc` 类型的方法。这也是通过反射实现的，核心代码也很简短：
```go
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

	for i := 0; i < v.NumMethod(); i++ {
		fn := v.Method(i) // 所有方法都必须为 ReqResFunc 类型
		handles = append(handles, ReqResHandle(fn.Interface()))
	}

	return handles
}
```

但是通过对象注册路由有个缺点，难以为 **`HandlerFunc`** 绑定 `path` 和 `method`。

已知的解决方式：
1. `GoFrame` 是在 `Req`（第二个请求参数）里写 `go tag`，有兴趣的读者，可以查看[「文档：规范参数结构」](https://goframe.org/pages/viewpage.action?pageId=116004922#id-规范路由如何使用-规范参数结构)。
2. `iris` 要求方法名（函数名）的格式为：请求方法+请求路径，如 `GetHelloWorld` 对应 `GET: /hello/world`，示例：[examples/mvc/hello-world/main.go](https://github.com/iris-contrib/examples/blob/master/mvc/hello-world/main.go)

笔者也做了简单的实现：[gee/web/day10/handle/handle_test.go](https://github.com/niluan304/gee/blob/main/web/day10/handle/handle_test.go)，这里就不再赘述了。
