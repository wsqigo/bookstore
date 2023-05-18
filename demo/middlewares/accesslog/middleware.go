package accesslog

import (
	"bookstore/demo/web"
	"encoding/json"
)

type MiddlewareBuilder struct {
	logFunc func(accessLog string)
}

func NewBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{}
}

func (b *MiddlewareBuilder) LogFunc(fn func(log string)) *MiddlewareBuilder {
	b.logFunc = fn
	return b
}

func (b MiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			// 我们在 defer 里面才最终输出日志，因为
			// 确保即便 next 里面发生了 panic，也能将请求记录下来
			// 获得 MatchedRoute: 它只有在执行了 next 之后才能获得，因为依赖于最终的路由树匹配 (HTTPServer.serve)
			defer func() {
				l := &accessLog{
					Host:       ctx.Req.Host,
					Route:      ctx.MatchedRoute,
					HTTPMethod: ctx.Req.Method,
					Path:       ctx.Req.URL.Path,
				}
				val, _ := json.Marshal(l)
				b.logFunc(string(val))
			}()
			next(ctx)
		}
	}
}

type accessLog struct {
	Host string `json:"host"`
	// 命中的路由
	Route      string `json:"route"`
	HTTPMethod string `json:"http_method"`
	Path       string `json:"path"`
}
