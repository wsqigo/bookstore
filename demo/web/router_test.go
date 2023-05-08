package web

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_router_AddRoute(t *testing.T) {
	// 第一个步骤是构造路由树
	// 第二个步骤是验证路由树
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodPost,
			path:   "/login",
		},
	}

	mockHandler := func(ctx *Context) {}
	r := newRouter()
	for _, route := range testRoutes {
		r.addRoute(route.method, route.path, mockHandler)
	}

	// 在这里断言路由树和你预期的一模一样
	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: {path: "/", children: map[string]*node{
				"user": {path: "user", children: map[string]*node{
					"home": {path: "home", handler: mockHandler},
				}, handler: mockHandler},
				"order": {path: "order", children: map[string]*node{
					"detail": {path: "detail", handler: mockHandler},
				}},
			}, handler: mockHandler},
			http.MethodPost: {path: "/", children: map[string]*node{
				"order": {path: "order", children: map[string]*node{
					"create": {path: "create", handler: mockHandler},
				}},
				"login": {path: "login", handler: mockHandler},
			}},
		},
	}

	// 断言两者相等
	msg, ok := wantRouter.equal(r)
	assert.True(t, ok, msg)

	// 这个是不行的，因为 HandleFunc 是不可比的
	//assert.Equal(t, wantRouter, r)

	// 非法用例
	r = newRouter()

	// 空字符串
	assert.PanicsWithValue(t, "web: 路由是空字符串", func() {
		r.addRoute(http.MethodGet, "", mockHandler)
	})

	// 前导没有 /
	assert.PanicsWithValue(t, "web: 路径必须以 / 开头", func() {
		r.addRoute(http.MethodGet, "a/b/c", mockHandler)
	})

	// 后缀有 /
	assert.PanicsWithValue(t, "web: 路径不能以 / 结尾", func() {
		r.addRoute(http.MethodGet, "/a/b/c/", mockHandler)
	})

	// 多个 /
	assert.PanicsWithValue(t, "web: 非法路由。不允许使用 //a/b, /a//b 之类的路由,[/a/b//c]", func() {
		r.addRoute(http.MethodGet, "/a/b//c", mockHandler)
	})

	r = newRouter()
	r.addRoute(http.MethodGet, "/", mockHandler)
	assert.PanicsWithValue(t, "web: 路由冲突[/]", func() {
		r.addRoute(http.MethodGet, "/", mockHandler)
	})

	r.addRoute(http.MethodGet, "/a/b/c", mockHandler)
	assert.PanicsWithValue(t, "web: 路由冲突，重复注册[/a/b/c]", func() {
		r.addRoute(http.MethodGet, "/a/b/c", mockHandler)
	})
}

// string 返回一个错误信息，帮助我们排查问题
// bool 是代表是否真的相等
func (r *router) equal(y *router) (string, bool) {
	for k, v := range r.trees {
		dst, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("找不到对应的 http method. %s", k), false
		}

		// v, dst 要相等
		msg, equal := v.equal(dst)
		if !equal {
			return msg, false
		}
	}

	return "", true
}

func (n *node) equal(y *node) (string, bool) {
	if n.path != y.path {
		return fmt.Sprintf("节点路径不匹配"), false
	}

	// 比较 handler
	nHandler := reflect.ValueOf(n.handler)
	yHandler := reflect.ValueOf(y.handler)
	if nHandler != yHandler {
		return fmt.Sprintf("handler 不相等"), false
	}

	if len(n.children) != len(y.children) {
		return fmt.Sprintf("子节点数量不相等"), false
	}

	for path, c := range n.children {
		dst, ok := y.children[path]
		if !ok {
			return fmt.Sprintf("子节点 %s 不存在", path), false
		}
		msg, ok := c.equal(dst)
		if !ok {
			return msg, false
		}
	}

	return "", true
}

func Test_router_findRoute(t *testing.T) {
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
	}

	r := newRouter()
	mockHandler := func(ctx *Context) {}
	for _, route := range testRoutes {
		r.addRoute(route.method, route.path, mockHandler)
	}

	testCases := []struct {
		// 子用例名
		name string
		// 入参
		method string
		path   string
		// 返回值
		wantFound bool
		wantNode  *node
	}{
		{
			// 方法都不存在
			name:      "method not found",
			method:    http.MethodHead,
			wantFound: false,
			wantNode:  nil,
		},
		{
			// 路径不存在
			name:      "path not found",
			method:    http.MethodGet,
			path:      "/abc",
			wantFound: false,
			wantNode:  nil,
		},
		{
			name:      "root",
			method:    http.MethodGet,
			path:      "/",
			wantFound: true,
			wantNode: &node{
				path:    "/",
				handler: mockHandler,
			},
		},
		{
			name:      "user",
			method:    http.MethodGet,
			path:      "/user",
			wantFound: true,
			wantNode: &node{
				path:    "user",
				handler: mockHandler,
			},
		},
		{
			name:      "no handler",
			method:    http.MethodPost,
			path:      "/order",
			wantFound: true,
			wantNode: &node{
				path: "order",
			},
		},
		{
			// 完全命中
			name:      "two layer",
			method:    http.MethodPost,
			path:      "/order/create",
			wantFound: true,
			wantNode: &node{
				handler: mockHandler,
				path:    "create",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			n, found := r.findRoute(tc.method, tc.path)
			assert.Equal(t, tc.wantFound, found)
			if !found {
				return
			}

			assert.Equal(t, tc.wantNode.path, n.path)
			nHandler := reflect.ValueOf(tc.wantNode.handler)
			yHandler := reflect.ValueOf(n.handler)
			assert.Equal(t, yHandler, nHandler)
		})
	}
}
