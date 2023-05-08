package web

import (
	"net"
	"net/http"
)

/*
Server从特性上来说，至少要提供三部分功能：
1. 生命周期控制：即启动、关闭。如果在后期，我们还要考虑增加生命周期回调特性
2. 路由注册接口：提供路由注册功能
3. 作为http包到Web框架的桥梁 http包暴露了一个接口Handler，它是我们引入自定义Web框架相关的连接点
*/

type HandleFunc func(ctx *Context)

// 确保 HTTPServer 一定实现了 Server 接口
var _ Server = &HTTPServer{}

// Server 组合http.Handler并且增加Start方法
// 既可以当成普通的http.Handler来使用，又可以一个独立的实体，拥有自己的管理生命周期的能力
type Server interface {
	http.Handler
	// Start 启动服务器
	// addr 是监听地址，如果只指定端口，可以使用 ":8081"
	// 或者 "localhost:8081"
	Start(addr string) error
	// AddRoute 注册一个路由
	// method 是 HTTP 方法
	// path 是路径，必须以 / 为开头
	// handleFunc 是你的业务逻辑
	addRoute(method string, path string, handleFunc HandleFunc)
	// 我们并不采取这种设计方案
	// addRoute(method string, path string, handles ...HandleFunc)

}

type HTTPServer struct {
	// addr string 创建的时候传递，而不是 Start 接收。这个都是可以的
	*router
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{
		router: newRouter(),
	}
}

// ServeHTTP HTTPServer 处理请求的入口
// web框架的核心入口。我们将整个方法内部完成：
// 1. Context构建 2. 路由匹配 3. 执行业务逻辑
func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 你的框架代码就在这里
	ctx := &Context{
		Req:  r,
		Resp: w,
	}

	s.serve(ctx)
}

// 查找路由，执行代码
func (s *HTTPServer) serve(ctx *Context) {
	//
}

func (s *HTTPServer) Start(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	// 在这里，可以让用户注册所谓的 after start 回调
	// 比如说往你的 admin 注册一下自己这个实例
	// 在这里执行一些你业务所需的前置条件

	return http.Serve(l, s)
}

//func (s *HTTPServer) addRoute(method string, path string, handleFunc HandleFunc) {
//	// 这里注册到路由树里面
//}

func (s *HTTPServer) Get(path string, handleFunc HandleFunc) {
	s.addRoute(http.MethodGet, path, handleFunc)
}

func (s *HTTPServer) Post(path string, handleFunc HandleFunc) {
	s.addRoute(http.MethodPost, path, handleFunc)
}
