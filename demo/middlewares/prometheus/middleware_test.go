package prometheus

import (
	"bookstore/demo/web"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

// 启动之后，访问一下 localhost:8081/user
// 然后再访问一下 localhost:2112/metrics
// 就能看到类似的输出，注意找一下

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := MiddlewareBuilder{
		Namespace: "bookstore",
		Subsystem: "web",
		Name:      "http_response",
	}
	s := web.NewHTTPServer(web.ServerWithMiddleware(builder.Build()))
	s.Get("/", func(ctx *web.Context) {
		ctx.Resp.Write([]byte("hello, world"))
	})
	s.Get("/user", func(ctx *web.Context) {
		val := rand.Intn(1000) + 1
		time.Sleep(time.Duration(val) * time.Millisecond)
		ctx.RespJSON(http.StatusAccepted, User{Name: "Tom"})
	})

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		// 一般来说，在实际中我们都会单独准备一个端口给这种监控
		http.ListenAndServe(":2112", nil)
	}()
	s.Start(":8081")
}

type User struct {
	Name string
}
