package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"gee/web/day8/handle"
	"gee/web/day8/internal/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	{
		user := r.Group("/user")

		user.GET("/", handle.Handle(controller.User.Get))
		user.POST("/", handle.Handle(controller.User.Add))
		user.POST("/upsert", handle.ReqResHandle(controller.User.Upsert))
	}

	go client()

	r.Run()
}

func client() {
	time.Sleep(time.Second) // 等待路由注册

	reqs := []func(host string) (*http.Response, error){
		func(host string) (*http.Response, error) { return http.Get(host + "/user?name=Carol") },
		func(host string) (*http.Response, error) { return http.Get(host + "/user?name=Bob") },
		func(host string) (*http.Response, error) {
			return http.Post(host+"/user", "application/json", bytes.NewBufferString(`{"name":"Carol","age":44,"job":"worker"}`))
		},
		func(host string) (*http.Response, error) { return http.Get(host + "/user?name=Carol") },

		// 测试 upsert 接口
		func(host string) (*http.Response, error) {
			return http.Post(host+"/user/upsert", "application/json", bytes.NewBufferString(`{"name":"Dave","age":32,"job":"nurse"}`))
		},
		func(host string) (*http.Response, error) {
			return http.Post(host+"/user/upsert", "application/json", bytes.NewBufferString(`{"name":"Dave","age":35,"job":"doctor"}`))
		},
		func(host string) (*http.Response, error) { return http.Get(host + "/user?name=Dave") },
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
	// {"code":200,"msg":"","data":null}
	// {"code":200,"msg":"","data":{"name":"Dave","age":32,"job":"nurse"}}
	// {"code":200,"msg":"","data":{"name":"Dave","age":35,"job":"doctor"}}
}
