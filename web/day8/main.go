package main

import (
	"gee/web/day8/internal/controller"

	"github.com/gin-gonic/gin"
)

// output:
//
// curl -X GET "http://localhost:8080/user?name=Carol"
// {"code":200,"msg":"","data":null}
//
// curl -X GET "http://localhost:8080/user?name=Bob"
// {"code":200,"msg":"","data":{"Name":"Bob","Age":30,"Job":"driver"}}
//
// curl -X POST "http://localhost:8080/user?name=Carol&age=44&job=worker"
// {"code":200,"msg":"","data":null}
//
// curl -X GET "http://localhost:8080/user?name=Carol"
// {"code":200,"msg":"","data":{"Name":"Carol","Age":44,"Job":"worker"}}

func main() {
	r := gin.Default()

	{
		user := r.Group("/user")
		user.GET("/", controller.User.Get)
		user.POST("/", controller.User.Add)
	}

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
