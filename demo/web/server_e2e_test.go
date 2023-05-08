//go:build e2e

package web

import (
	"fmt"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	var s Server = &HTTPServer{}

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

	// 用法二
	s.Start(":8081")
}
