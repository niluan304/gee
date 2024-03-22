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

	resp1, _ := http.Get("http://localhost:8080/user?name=Carol")
	resp2, _ := http.Get("http://localhost:8080/user?name=Bob")
	resp3, _ := http.Post("http://localhost:8080/user", "application/json", bytes.NewBufferString(`{"name":"Carol","age":44,"job":"worker"}`))
	resp4, _ := http.Get("http://localhost:8080/user?name=Carol")

	for _, resp := range []*http.Response{resp1, resp2, resp3, resp4} {
		data, _ := io.ReadAll(resp.Body)
		fmt.Println(string(data))
	}

	// Output:
	// {"code":200,"msg":"","data":null}
	// {"code":200,"msg":"","data":{"Name":"Bob","Age":30,"Job":"driver"}}
	// {"code":200,"msg":"","data":null}
	// {"code":200,"msg":"","data":{"Name":"Carol","Age":44,"Job":"worker"}}
}
