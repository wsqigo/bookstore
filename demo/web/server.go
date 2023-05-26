package web

import (
	"log"
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

// Server 组合 http.Handler 并且增加 Start 方法
// 既可以当成普通的 http.Handler 来使用，又可以一个独立的实体，拥有自己的管理生命周期的能力
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

type HTTPServerOption func(server *HTTPServer)

type HTTPServer struct {
	// addr string 创建的时候传递，而不是 Start 接收。这个都是可以的
	router

	mdls []Middleware
}

func NewHTTPServer(opts ...HTTPServerOption) *HTTPServer {
	res := &HTTPServer{
		router: newRouter(),
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

//// 第一个问题：相对路径还是绝对路径？
//// 你的配置文件格式，json, yaml, xml
//func NewHTTPServerV2(cfgFilePath string) *HTTPServer {
//	// 你在这里加载配置，解析，然后初始化 HTTPSever
//}

func ServerWithMiddleware(mdls ...Middleware) HTTPServerOption {
	return func(server *HTTPServer) {
		server.mdls = mdls
	}
}

// ServeHTTP HTTPServer 处理请求的入口
// web框架的核心入口。我们将整个方法内部完成：
// 1. Context 构建 2. 路由匹配 3. 执行业务逻辑
func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 你的框架代码就在这里
	ctx := &Context{
		Req:  r,
		Resp: w,
	}

	// 最后一个应该是 HTTPServer 执行路由匹配，执行用户代码
	root := s.serve
	// 然后这里就是利用最后一个不断往前回溯组装链条
	// 从后往前
	// 把后一个作为前一个的 next 构造好链条
	for i := len(s.mdls) - 1; i >= 0; i-- {
		root = s.mdls[i](root)
	}

	// 第一个应该是回写响应的
	// 因为它在调用 next 之后才回写响应
	// 所以实际上 flashResp 是最后一个步骤
	var m Middleware = func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			// 设置好了 RespData 和 RespStatusCode
			next(ctx)
			s.flashResp(ctx)
		}
	}
	root = m(root)
	root(ctx)
}

func (s *HTTPServer) flashResp(ctx *Context) {
	if ctx.RespStatusCode != 0 {
		ctx.Resp.WriteHeader(ctx.RespStatusCode)
	}
	_, err := ctx.Resp.Write(ctx.RespData)
	if err != nil {
		log.Fatalln("回写响应失败", err)
	}
}

// 查找路由，执行代码
func (s *HTTPServer) serve(ctx *Context) {
	r := ctx.Req
	info, found := s.findRoute(r.Method, r.URL.Path)
	if !found || info.n == nil || info.n.handler == nil {
		// 路由没有命中，就是404
		ctx.RespStatusCode = http.StatusNotFound
		ctx.RespData = []byte("NOT FOUND")
		return
	}

	ctx.PathParams = info.pathParams
	ctx.MatchedRoute = info.n.route
	info.n.handler(ctx)
}

// Start 启动服务器，用户指定端口
// 这种就是编程接口
func (s *HTTPServer) Start(addr string) error {
	// 也可以自己创建 Server
	// http.Server{}
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
