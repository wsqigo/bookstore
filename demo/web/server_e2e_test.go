//go:build e2e

package web

import (
	"bytes"
	"fmt"
	"html/template"
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

	s.Get("/login", func(ctx *Context) {
		// 返回登录页面
		tpl := template.New("login")
		tpl, err := tpl.Parse(`
<html>
	<body>
		<form>
			// 在这里继续写页面
		</form>
	</body>
</html>
`)
		if err != nil {
			t.Fatal(err)
		}

		page := &bytes.Buffer{}
		err = tpl.Execute(page, nil)
		if err != nil {
			t.Fatal(err)
		}

		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = page.Bytes()
	})

	// 用法二
	s.Start(":8081")
}
