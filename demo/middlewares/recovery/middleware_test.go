package recovery

import (
	"bookstore/demo/web"
	"log"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := MiddlewareBuilder{
		StatusCode: 500,
		ErrMsg:     "你 Panic 了",
		LogFunc: func(ctx *web.Context) {
			log.Println("panic 路径:", ctx.Req.URL.Path)
		},
	}

	s := web.NewHTTPServer(web.ServerWithMiddleware(builder.Build()))
	s.Get("/user", func(ctx *web.Context) {
		ctx.RespData = []byte("hello, world")
	})
	s.Get("/panic", func(ctx *web.Context) {
		panic("发生 panic 了")
	})
	s.Start(":8081")
}
