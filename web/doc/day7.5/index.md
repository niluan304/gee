---
title: "gee-web-day7.5 「7天 web框架」读后感"
description: "An Idea: 将 gin 改造为 GoFrame"
date: 2023-12-19T11:27:57+08:00
image: 
license: 
hidden: false
comments: true
draft: false
math: false
tags: [gee,web,gin,goframe]
categories: go
---

跟着这篇教程 [7天用Go从零实现Web框架Gee教程 | 极客兔兔](https://geektutu.com/post/gee.html)，笔者实现了一个简单的 web框架，从而明白了 web框架所需要的特性，还迸发出封装 `gin` 的想法。

## web框架 所需要的特性

### 分组控制

分组控制 (Group Control)是 Web 框架应当提供的基础功能，有了分组与分组嵌套，系统才可以更方便的管理不同的业务，搭配中间件，就可以处理同组下的公共逻辑。

### 中间件

中间件 (Middleware)是 web 框架的灵魂，为 web 框架提供无限的扩展能力。有了中间件，可以只对 `admin` 和 `api` 分组进行鉴权，对特定接口限流。对最顶层目录 "/"，也就是整个系统，可以配置日志，请求耗时等功能。

### 动态路由

动态路由 (Dynamic Route)，本质上只是将参数给映射为路径的一部分了，比如 `https://github.com/niluan304` 和 `https://github.com/?user=niluan304` 对于后端来讲并没有什么不同。不过动态路由还是有优势的：
- 提高SEO：搜索引擎解析静态URL更为轻松，动态路由将参数内嵌至URL中，可以提高SEO效果。
- 可读性和可维护性：现阶段的动态路由可以有多个参数，如 `/user/:name/article/:articleId`，可读性和可维护性明显高于 `user/article?name=zwei&?articleId=1234`。
- 可拓展性：动态路由可以使用前缀进行分组，这样可以很容易添加和修改同一分组下的中间件。

## `gin` 与 `GoFrame` 的差异
实现简易web框架之后，才算理解 `gin` 的中间件与业务代码的格式为什么必须为：
```go
func(c *gin .Context) {
    // 中间件 或 业务代码逻辑
}
```

但是这和笔者最熟悉的 `GoFrame` 框架有一些差异，`GoFrame` 还额外支持另一种写法：
```go
// 中间件式写法，类似 gin  的 func(c *gin .Context)
func (r *ghttp.Request) {
    // 中间件 或 业务代码逻辑
}

// 写法二
func (ctx context.Context, req *{Prefix}Req) (res *{Prefix}Res, err error){
    // 业务代码逻辑
}
```

带着疑问，笔者探究了两个框架之间的差异，并尝试将 `gin` 的接口改造为 `GoFrame` 格式。
总结后可以概括为两个部分，这算是「7天教程」的读后感，笔者也仿照命名为：
- day8 分层设计的必要性
- day9 反序列化与解耦