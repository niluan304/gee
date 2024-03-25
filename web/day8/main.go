package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"gee/web/day8/internal/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	{
		user := r.Group("/user")
		user.GET("/", controller.User.Get)
		user.POST("/", controller.User.Add)
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
