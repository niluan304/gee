---
title: "gee-web-day8 分层设计的必要性"
description: 通过请求分层流转，实现业务代码与 web 框架的解耦
date: 2023-12-20T13:49:08+08:00
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

## 什么是请求分层流转

在阐述为什么需要分层设计之前，笔者想先介绍一下分层流转[^1]：
[^1]:[工程目录设计🔥 - GoFrame (ZH)](https://goframe.org/pages/viewpage.action?pageId=30740166)

![](image.png)
分层之后，可以让每一个函数只负责一件事，这类似于设计模式里单一职责的思想，可以提高项目的可维护性，与接口的可复用性。


## 代码纠缠的困境
这里有一份简单的 `CURD` 代码，目前只有两个功能，`Add, Get`：
```go
type Response struct {
	Code int    // 业务代码，200 表示 OK，其他表示错误
	Msg  string // 错误消息
	database any    // 返回的数据
}

type Row struct {
	Name string
	Age  int
	Job  string
}

// database 只是一个切片 []Row，用于充当数据库
var database = []Row{ {"Alice", 32, "teacher"}, {"Bob", 30, "driver"} }

var User = &user{}

type user struct{}

func (c *user) Add(ctx *gin.Context) {
	var req Row
	err := ctx.ShouldBind(&req) // 请求数据的反序列化
	if err != nil {
		ctx.JSON(http.StatusOK, Response{400, err.Error(), nil})
		return
	}

	// 插入数据，database 只是一个切片 []Row，用于充当数据库
	database = append(database, req)

	// 填写响应内容
	ctx.JSON(http.StatusOK, Response{200, "", nil})
	return
}

func (c *user) Get(ctx *gin.Context) {
	var req Row
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, Response{400, err.Error(), nil})
		return
	}

	// 查询数据，database 只是一个切片 []Row，用于充当数据库
	i := slices.IndexFunc(database, func(row Row) bool { return row.Name == req.Name })
	if i != -1 {
		// 填写响应内容
		ctx.JSON(http.StatusOK, Response{200, "", database[i]})
	}
	return
}
```

看起来没什么问题，但假如随着项目的进展，导致 `Add` 部分业务代码和 `Get` 的完全一致，甚至最好的解决办法是直接调用 `Get` 方式，但这没办法真的去调用 `Get` 方式，因为 `Add` 接收到的 `*gin.Context` 与 `Get` 方法需要的 `*gin.Context` 是有差别的。

这时候就可以进行分层设计，现在比较流行的纯后端 API 模块一般采用下述划分方法[^2]
[^2]:[大型Web项目分层 - Go语言高级编程](https://chai2010.cn/advanced-go-programming-book/ch5-web/ch5-07-layout-of-web-project.html)：
1. Controller，与上述类似，服务入口，负责处理路由，参数校验，请求转发。
2. Logic/Service，逻辑（服务）层，一般是业务逻辑的入口，可以认为从这里开始，所有的请求参数一定是合法的。业务逻辑和业务流程也都在这一层中。常见的设计中会将该层称为 Business Rules。
3. DAO/Repository，这一层主要负责和数据、存储打交道。将下层存储以更简单的函数、接口形式暴露给 Logic 层来使用。负责数据的持久化工作。

## 分层设计
先介绍下分层后的目录结构：
```sh
.
|-- go.mod
|-- go.sum
|-- internal
|   |-- controller
|   |   `-- user.go
|   `-- service
|       |-- user.go
|       `-- user_model.go
`-- main.go
```

当业务代码都放到了 `service` 层时，这一层的代码互相调用是不会被 `controller` 层影响的，这也实现了 `gin` 框架与业务代码的解耦。

`controller` 层的主要代码：
```go
// ./internal/controller/user.go

func (c *user) Add(ctx *gin.Context) {
	// 请求数据的反序列化
	var req *service.UserAddReq
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
	var req *service.UserGetReq
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
```

`service` 层的主要代码：
```go
// ./internal/service/user.go

func (s *user) Add(ctx context.Context, req *UserAddReq) (res *UserAddRes, err error) {
	// 插入数据
	database = append(database, Row{req.Name, req.Age, req.Job})
	return
}

func (s *user) Get(ctx context.Context, req *UserGetReq) (res *UserGetRes, err error) {
	// 查询数据
	i := slices.IndexFunc(database, func(row Row) bool { return row.Name == req.Name })
	if i != -1 {
		// 填写响应内容
		row := database[i]
		return &UserGetRes{Name: row.Name, Age: row.Age, Job: row.Job}, nil
	}
	return
}
```

分层后的总代码行数有所增加，甚至 `controller` 错误处理变得更繁琐了，但是整个项目的布局变得更清晰了，业务代码也不会受到 web框架的干扰，可以集中处理业务。

更可贵的是，`service` 层的方法在调用时，就可以知道所需要的参数，以及返回的值。不过有些读者可能会有疑问，为什么 `service` 层方法的第一个参数都是 `ctx context.Context`，即便代码中未必使用，这算是 `go` 语言在 web 开发中的特色（也可能是技术债），用于并发控制和上下文信息传递的，有兴趣可以自行了解下。



## 小结
本章节介绍下「分层设计」与「单一职责」的联系，并说明如何通过分层设计将业务代码与 web框架解耦。

- 注意：分层设计也会导致一个问题：新增一个业务接口时，需要改动的文件也会变多，不过这可以通过脚本生成代码缓解。


最后让我们来看看程序的运行结果：
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
