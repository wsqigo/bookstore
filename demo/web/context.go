package web

import "net/http"

/*
为什么都已经有了 http 包，还要开发 Web 框架?
高级路由功能，封装HTTP上下文以提供简单API、封装 Server 以提供生命周期控制、设计插件机制以提供无侵入式解决方案，提供如上传下载等默认功能
*/

type Context struct {
	Req        *http.Request
	Resp       http.ResponseWriter
	PathParams map[string]string
}
