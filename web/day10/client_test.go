package main

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

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
