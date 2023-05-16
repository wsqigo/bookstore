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
	s.Get("/user", func(ctx *Context) {
		handler1(ctx)
		handler2(*ctx)
	})

	s.Get("/order/detail", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, order detail"))
	})

	s.Post("/form", func(ctx *Context) {
		ctx.Req.ParseForm()
		ctx.Resp.Write([]byte(fmt.Sprintf("hello, %s", ctx.Req.URL.Path)))
	})

	s.Post("/values/:id", func(ctx *Context) {
		id, err := ctx.PathValueV2("id").String()
		if err != nil {
			ctx.Resp.WriteHeader(http.StatusBadRequest)
			ctx.Resp.Write([]byte("id 输入不对"))
		}
		ctx.Resp.Write([]byte(fmt.Sprintf("hello, %s", id)))
	})

	// 用法二
	s.Start(":8081")
}
