package handle_test

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"
	"unicode"

	"gee/web/day8/handle"

	"github.com/gin-gonic/gin"
)

func TestObjectHandler(t *testing.T) {
	r := gin.Default()

	// 创建一个 HTTP 服务器
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	defer server.ListenAndServe()

	go func() {
		// 等待路由注册
		time.Sleep(time.Second * 2)

		http.Get("http://localhost:8080/gf/hello-world")
		http.Post("http://localhost:8080/gf/hello-world", "application/json", nil)

		http.Get("http://localhost:8080/iris/hello/world")
		http.Post("http://localhost:8080/iris/hello/world", "application/json", nil)

		// 等待打印请求响应
		time.Sleep(time.Second * 2)
		server.Shutdown(context.TODO())
	}()

	// 根据结构体的 tag 注册路由
	gf := r.Group("/gf")
	handle.ObjectHandler(Hello{}, func(f *handle.ReqResFunc, methodName string) {
		field, find := f.Req().Elem().FieldByName("meta")
		if !find {
			panic("req must be contain meta filed")
		}

		tag := field.Tag
		gf.Handle(strings.ToUpper(tag.Get("method")), tag.Get("path"), f.Handler())
	})

	// 根据方法名注册路由
	iris := r.Group("/iris")
	handle.ObjectHandler(Hello{}, func(f *handle.ReqResFunc, methodName string) {
		var name []rune
		for _, r := range methodName {
			if unicode.IsUpper(r) {
				name = append(name, '/')
			}
			name = append(name, unicode.ToLower(r))
		}
		path := string(name)[1:]
		i := strings.Index(path, "/")
		iris.Handle(strings.ToUpper(path[:i]), path[i:], f.Handler())
	})
}

type (
	HelloGetReq struct {
		meta struct{} `method:"GET" path:"/hello-world" `
	}
	HelloGetRes struct{}

	HelloPostReq struct {
		meta struct{} `method:"POST" path:"/hello-world" `
	}
	HelloPostRes struct{}
)

type Hello struct{}

func (c *Hello) GetHelloWorld(ctx context.Context, req *HelloGetReq) (res *HelloGetRes, err error) {
	return &HelloGetRes{}, nil
}

func (c *Hello) PostHelloWorld(ctx context.Context, req *HelloPostReq) (res *HelloPostRes, err error) {
	return &HelloPostRes{}, nil
}
