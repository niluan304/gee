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
