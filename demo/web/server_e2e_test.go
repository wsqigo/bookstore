//go:build e2e

package web

import (
	"fmt"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	// 直接创建router为nil
	//var s Server = &HTTPServer{}
	s := NewHTTPServer()

	handler1 := func(ctx *Context) {
		fmt.Println("处理第一件事")
	}
	handler2 := func(ctx Context) {
		fmt.Println("处理第二件事")
	}

	// 用户自己去管这种
	s.addRoute(http.MethodGet, "/user", func(ctx *Context) {
		handler1(ctx)
		handler2(*ctx)
	})

	s.addRoute(http.MethodGet, "/order/detail", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, order detail"))
	})

	// 用法二
	s.Start(":8081")
}
