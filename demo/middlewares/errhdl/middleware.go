package errhdl

import "bookstore/demo/web"

type MiddlewareBuilder struct {
	// 这种设计只能返回固定的值
	// 不能做到动态渲染
	resp map[int][]byte
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		// 这里可以非常大方，因为在设计中用户会关心的错误码不可能超过 64
		resp: make(map[int][]byte, 64),
	}
}

// RegisterError 注册一个错误码，并且返回特定的错误数据
// 这个错误可以是一个字符串，也可以是一个页面
func (b *MiddlewareBuilder) RegisterError(status int, data []byte) *MiddlewareBuilder {
	b.resp[status] = data
	return b
}

func (b MiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			next(ctx)
			resp, ok := b.resp[ctx.RespStatusCode]
			if ok {
				ctx.RespData = resp
			}
		}
	}
}
