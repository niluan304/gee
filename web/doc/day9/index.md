---
title: "gee-web-day9 反序列化与解耦"
description: "更优雅的反序列化，并实现支持多种协议"
date: 2023-12-21T18:15:11+08:00
image: 
license: 
hidden: false
comments: true
draft: false
math: false
tags: [gee,web,gin]
categories: go
---

> 源代码/数据集已上传到：[GitHub - follow gee to learn go](https://github.com/niluan304/gee)

经过分层处理后，项目布局有了很大改善，但是仍然存在问题。
1. `controller` 层的错误处理代码特别繁琐，有太多的：`ctx.JSON(http.StatusOK, Response{400, err.Error(), nil})`
2. `controller` 层只支持 `gin` 框架，更不支持其他协议


## 让调用者帮忙反序列化
观察 `controller` 层可以得出一个结论：`controller` 的需求其实很简单：反序列化为 `service` 层所需要 `go` 类型，并在  `err != nil` 时做控制的流转。

那么该怎么实现呢？这其实并不难。

不知道是否有读者会好奇 `Get` 方法里的搜索函数：
```go
// ./day8/internal/service/user.go

i := slices.IndexFunc(database, func(row Row) bool { 
    return row.Name == req.Name
})
```
只需要往 `slices.IndexFunc` 传入了一个闭包，就可以等待函数返回所在值的索引，`slices.IndexFunc`的实现，也同样简单：
```go
// std: slices/slices.go

// IndexFunc returns the first index i satisfying f(s[i]),
// or -1 if none do.
func IndexFunc[S ~[]E, E any](s S, f func(E) bool) int {
	for i := range s {
		if f(s[i]) {
			return i
		}
	}
	return -1
}
```

但让我们换个角度思考，对于 `slices.IndexFunc` 来说，它并不关心 `s` 里的元素，只关心 `f(s[i])` 的结果，因此只需要调用传入的闭包，就完成了任务。同样地，对于 `controller` 层来说，只关心反序列化 `var req service.UserGetReq`。那么能否仿照 `slices.IndexFunc` 让外部传入一个闭包，让闭包帮 `controller` 层反序列化 `req`呢？答案是可以的，只需要修改 `controller` 的入参为：
```go
func (c *user) Get(
	ctx context.Context,               // 第一个参数：ctx
	bind func(point any) (err error),  // 第二个参数：用于反序列化的闭包
) (
	data any,  // 返回的数据
	err error, // 错误处理
){
	var req service.UserGetReq
	err = bind(&req) // 通过闭包反序列化 req
	if err != nil {
		return nil, err
	}

	res, err := service.User.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
```

经过改动的 `controller` 和 web框架彻底解耦了，完全看不到 `gin` 框架的代码，毕竟反序列化的工作也并不是 `controller` 的任务，错误处理和数据返回也变得非常简单，只需要抛给上层处理（要么 `return nil, err`，要么 `return res, nil`），解耦之后也为 `controller` 层兼容多种协议带来了可能。

不过但也带来了一个问题：这样的函数，该如何注册到 `gin`框架里呢？


## 统一错误处理和数据返回
阐述这部分内容之前，笔者想简单的介绍一下「设计模式」里的「适配器模式」[^1]：
[^1]:[适配器模式](https://www.yuque.com/aceld/lfhu8y/vnhf4b#gVTIW)

简单来讲，就是通过接口转换，让两个不兼容的接口，能够一起工作，现实中的经典例子：
![](image.png)

和上面的图片类似，修改后的函数类型已经和框架要求的 `gin.HandlerFunc` 截然不同，但借鉴适配器模式的思想，通过中间函数转化，就可以了：
```go
// ./handle/handle.go

// 设置为类型，用于优化参数显示
type BinderFunc = func(
	ctx context.Context,               // 第一个参数：ctx
	bind func(point any) (err error), // 第二个参数：用于反序列化的闭包
) (
	data any,  // 返回的数据
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

```

新的函数作为「适配器」，也会被其他路由调用，也不是业务相关的内容，不适合放到 `interna`l 包，应当放到新的包（文件夹）里，笔者将之保存至 `/handle/handle.go`。

相应地，路由注册也有些变化：
```go
func main() {
	r := gin.Default()
	{
		user := r.Group("/user")
		user.GET("/", handle.Handle(controller.User.Get))  
		user.POST("/", handle.Handle(controller.User.Add))
	}

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
```

对比 day8 的注册模式：
```go
user.GET("/", controller.User.Get)  // 函数签名：func (c *user) Get(ctx *gin.Context)
user.POST("/", controller.User.Add) // 函数签名：func (c *user) Add(ctx *gin.Context)
```

虽然注册路由时，必须得借用 `handle.Handle` 才能转化为 `gin.HandlerFunc`，但是可以不用在 `controller`层里写 `c.JSON(http.StatusOK, Response{Code: 400, Msg: err.Error(), Data: nil})` 和 `c.JSON(http.StatusOK, Response{200, "", data})` 了。至此，我们就完成了错误处理和数据返回的统一。



## 小结
本章节主要做了两件事：
1. `controller` 层通过传入 `binder` 闭包，调用闭包实现反序列化出 `service` 层所需数据，也完成了 `controller` 与框架的解耦，日后可以兼容其他框架（如 echo）和其他协议（如 rpc）。
2. 借鉴适配器模式，将 `binderFunc` 函数转化为框架所需要的类型，并实现错误处理和数据返回的统一。

运行结果也没有变化：
```go
func client() {
	time.Sleep(time.Second) // 等待路由注册

	reqs := []func(host string) (*http.Response, error){
		func(host string) (*http.Response, error) { return http.Get(host + "/user?name=Carol") },
		func(host string) (*http.Response, error) { return http.Get(host + "/user?name=Bob") },
		func(host string) (*http.Response, error) {
			return http.Post(host+"/user", "application/json", bytes.NewBufferString(`{"name":"Carol","age":44,"job":"worker"}`))
		},
		func(host string) (*http.Response, error) { return http.Get(host + "/user?name=Carol") },
	}

	for _, req := range reqs {
		resp, err := req("http://localhost:8080")
		if err != nil {
			fmt.Println("req err", err)
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("read resp.Body err", err)
		}
		fmt.Println(string(data))
	}

	// Output:
	//
	// {"code":200,"msg":"","data":null}
	// {"code":200,"msg":"","data":{"name":"Bob","age":30,"job":"driver"}}
	// {"code":200,"msg":"","data":null}
	// {"code":200,"msg":"","data":{"name":"Carol","age":44,"job":"worker"}}
}
```
