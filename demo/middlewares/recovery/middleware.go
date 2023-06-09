package recovery

import (
	"bookstore/demo/web"
)

type MiddlewareBuilder struct {
	StatusCode int
	ErrMsg     string
	LogFunc    func(ctx *web.Context)
}

func (b MiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			defer func() {
				if err := recover(); err != nil {
					ctx.RespStatusCode = b.StatusCode
					ctx.RespData = []byte(b.ErrMsg)
					// 万一 LogFunc 也panic，那我们也无能为力了
					b.LogFunc(ctx)
				}
			}()
			next(ctx)
		}
	}
}