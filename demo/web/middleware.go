package web

// Middleware 函数式的责任链模式
// 函数式的洋葱模式
type Middleware func(next HandleFunc) HandleFunc

// AOP 方案在不同的框架，不同的语言里面都有不同的叫法
// Middleware, Handler, Chain, Filter, Filter-Chain
// Interceptor, Wrapper
//type MiddlewareV1 interface {
//	Invoke(next HandleFunc) HandleFunc
//}
//
//type Interceptor interface {
//	Before(ctx *Context)
//	After(ctx *Context)
//	Surround(ctx *Context)
//}
//
//type HandleFuncV1 func(ctx *Context) (next bool)
//
//type FilterChain []HandleFuncV1
//
//type FilterChainV1 struct {
//	handlers []HandleFuncV1
//}
//
//func (f FilterChainV1) Run(ctx *Context) {
//	for _, h := range f.handlers {
//		next := h(ctx)
//		// 这种是中断执行
//		if !next {
//			return
//		}
//	}
//}
