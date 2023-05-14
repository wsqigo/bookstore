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
		// 通配符测试用例
		{
			method: http.MethodGet,
			path:   "/order/*",
		},
		{
			method: http.MethodGet,
			path:   "/*",
		},
		{
			method: http.MethodGet,
			path:   "/*/*",
		},
		{
			method: http.MethodGet,
			path:   "/*/abc",
		},
		{
			method: http.MethodGet,
			path:   "/*/abc/*",
		},
		// 参数路由
		{
			method: http.MethodGet,
			path:   "/param/:id",
		},
		{
			method: http.MethodGet,
			path:   "/param/:id/detail",
		},
		{
			method: http.MethodGet,
			path:   "/param/:id/*",
		},
		// 正则路由
		{
			method: http.MethodDelete,
			path:   "/reg/:id(.*)",
		},
		{
			method: http.MethodDelete,
			path:   "/:name(^.+$)/abc",
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
			http.MethodGet: {
				path: "/",
				children: map[string]*node{
					"user": {path: "user", children: map[string]*node{
						"home": {path: "home", handler: mockHandler, typ: nodeTypeStatic},
					}, handler: mockHandler, typ: nodeTypeStatic},
					"order": {path: "order", children: map[string]*node{
						"detail": {path: "detail", handler: mockHandler, typ: nodeTypeStatic},
					}, starChild: &node{path: "*", handler: mockHandler, typ: nodeTypeAny}},
					"param": {
						path: "param",
						paramChild: &node{
							path:      ":id",
							paramName: "id",
							children: map[string]*node{
								"detail": {path: "detail", handler: mockHandler},
							},
							starChild: &node{path: "*", handler: mockHandler},
							handler:   mockHandler,
							typ:       nodeTypeParam,
						},
					},
				},
				starChild: &node{
					path: "*",
					children: map[string]*node{
						"abc": {path: "abc", handler: mockHandler, starChild: &node{path: "*", handler: mockHandler, typ: nodeTypeStatic}},
					},
					starChild: &node{path: "*", handler: mockHandler},
					handler:   mockHandler,
					typ:       nodeTypeAny,
				},
				handler: mockHandler,
				typ:     nodeTypeStatic,
			},
			http.MethodPost: {path: "/", children: map[string]*node{
				"order": {path: "order", children: map[string]*node{
					"create": {path: "create", handler: mockHandler, typ: nodeTypeStatic},
				}},
				"login": {path: "login", handler: mockHandler, typ: nodeTypeStatic},
			}, typ: nodeTypeStatic},
			http.MethodDelete: {path: "/", children: map[string]*node{
				"reg": {path: "reg", typ: nodeTypeStatic, regChild: &node{
					path: ":id(.*)", paramName: "id", typ: nodeTypeReg, handler: mockHandler},
				},
			}, regChild: &node{
				path: ":name(^.+$)",
				children: map[string]*node{
					"abc": {path: "abc", typ: nodeTypeStatic, handler: mockHandler},
				},
				typ:       nodeTypeReg,
				paramName: "name",
			}},
		},
	}

	// 断言两者相等
	msg, ok := wantRouter.equal(&r)
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

	// 根节点重复注册
	r.addRoute(http.MethodGet, "/", mockHandler)
	assert.PanicsWithValue(t, "web: 路由冲突[/]", func() {
		r.addRoute(http.MethodGet, "/", mockHandler)
	})

	// 普通节点重复注册
	r.addRoute(http.MethodGet, "/a/b/c", mockHandler)
	assert.PanicsWithValue(t, "web: 路由冲突，重复注册[/a/b/c]", func() {
		r.addRoute(http.MethodGet, "/a/b/c", mockHandler)
	})

	// 多个 /
	assert.PanicsWithValue(t, "web: 非法路由。不允许使用 //a/b, /a//b 之类的路由,[/a/b//c]", func() {
		r.addRoute(http.MethodGet, "/a/b//c", mockHandler)
	})

	// 同时注册通配符路由和参数路由
	assert.PanicsWithValue(t, "web: 非法路由，已有通配符路由。不允许同时注册通配符路由和参数路由 [:id]", func() {
		r.addRoute(http.MethodGet, "/a/*", mockHandler)
		r.addRoute(http.MethodGet, "/a/:id", mockHandler)
	})
	assert.PanicsWithValue(t, "web: 非法路由，已有路径参数路由。不允许同时注册通配符路由和参数路由 [*]", func() {
		r.addRoute(http.MethodGet, "/a/b/:id", mockHandler)
		r.addRoute(http.MethodGet, "/a/b/*", mockHandler)
	})

	// 同时注册正则路由和参数路由
	r = newRouter()
	assert.PanicsWithValue(t, "web: 非法路由，已有路径参数路由。不允许同时注册正则路由和参数路由 [:id(.*)]", func() {
		r.addRoute(http.MethodGet, "/a/b/:id", mockHandler)
		r.addRoute(http.MethodGet, "/a/b/:id(.*)", mockHandler)
	})
	r = newRouter()
	assert.PanicsWithValue(t, "web: 非法路由，已有正则路由。不允许同时注册通配符路由和正则路由 [*]", func() {
		r.addRoute(http.MethodGet, "/a/b/:id(.*)", mockHandler)
		r.addRoute(http.MethodGet, "/a/b/*", mockHandler)
	})

	r = newRouter()
	assert.PanicsWithValue(t, "web: 非法路由，已有正则路由。不允许同时注册正则路由和参数路由 [:id]", func() {
		r.addRoute(http.MethodGet, "/a/b/:id(.*)", mockHandler)
		r.addRoute(http.MethodGet, "/a/b/:id", mockHandler)
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

	// 参数冲突
	assert.PanicsWithValue(t, "web: 路由冲突，参数路由冲突，已有 :id，新注册 :name", func() {
		r.addRoute(http.MethodGet, "/a/b/c/:id", mockHandler)
		r.addRoute(http.MethodGet, "/a/b/c/:name", mockHandler)
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
	if y == nil {
		return "目标节点为 nil", false
	}
	if n.path != y.path {
		return fmt.Sprintf("%s 节点路径不匹配 x %s, y %s", n.path, n.path, y.path), false
	}

	if n.starChild != nil {
		msg, ok := n.starChild.equal(y.starChild)
		if !ok {
			return msg, ok
		}
	}

	if n.paramChild != nil {
		msg, ok := n.paramChild.equal(y.paramChild)
		if !ok {
			return msg, ok
		}
	}

	if n.paramName != y.paramName {
		return fmt.Sprintf("%s 节点参数名字不相等 x %s, y %s", n.path, n.paramName, y.paramName), false
	}

	if n.regChild != nil {
		msg, ok := n.regChild.equal(y.regChild)
		if !ok {
			return msg, ok
		}
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
		{
			method: http.MethodPost,
			path:   "/order/*",
		},
		{
			method: http.MethodGet,
			path:   "/user/*/home",
		},
		// 参数路由
		{
			method: http.MethodGet,
			path:   "/param/:id",
		},
		{
			method: http.MethodGet,
			path:   "/param/:id/detail",
		},
		{
			method: http.MethodGet,
			path:   "/param/:id/*",
		},
		// 正则
		{
			method: http.MethodDelete,
			path:   "/reg/:id(.*)",
		},
		{
			method: http.MethodDelete,
			path:   "/:id([0-9]+)/home",
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
		info      *matchInfo
	}{
		{
			// 方法都不存在
			name:      "method not found",
			method:    http.MethodHead,
			wantFound: false,
			info:      nil,
		},
		{
			// 路径不存在
			name:      "path not found",
			method:    http.MethodGet,
			path:      "/abc",
			wantFound: false,
			info:      nil,
		},
		{
			name:      "root",
			method:    http.MethodGet,
			path:      "/",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					path:    "/",
					handler: mockHandler,
				},
			},
		},
		{
			name:      "user",
			method:    http.MethodGet,
			path:      "/user",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					path:    "user",
					handler: mockHandler,
				},
			},
		},
		{
			name:      "no handler",
			method:    http.MethodPost,
			path:      "/order",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					path: "order",
				},
			},
		},
		{
			// 完全命中
			name:      "two layer",
			method:    http.MethodPost,
			path:      "/order/create",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					handler: mockHandler,
					path:    "create",
				},
			},
		},
		// 通配符匹配
		{
			// 命中/order/*
			name:      "star match",
			method:    http.MethodPost,
			path:      "/order/delete",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					path:    "*",
					handler: mockHandler,
				},
			},
		},
		{
			// 命中通配符在中间的
			// /user/*/home
			name:      "star in middle",
			method:    http.MethodGet,
			path:      "/user/Tom/home",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					path:    "home",
					handler: mockHandler,
				},
			},
		},
		{
			// 比 /order/* 多了一段
			name:      "overflow",
			method:    http.MethodPost,
			path:      "/order/delete/123",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					path:    "*",
					handler: mockHandler,
				},
			},
		},
		// 参数匹配
		{
			// 命中 /param/:id
			name:      ":id",
			method:    http.MethodGet,
			path:      "/param/123",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					path:    ":id",
					handler: mockHandler,
				},
				pathParams: map[string]string{
					"id": "123",
				},
			},
		},
		{
			// 命中 /param/:id/detail
			name:      ":id/*",
			method:    http.MethodGet,
			path:      "/param/123/detail",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					path:    "detail",
					handler: mockHandler,
				},
				pathParams: map[string]string{
					"id": "123",
				},
			},
		},
		{
			// 命中 /reg/:id(.*)
			name:      ":id(.*)",
			method:    http.MethodDelete,
			path:      "/reg/123",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					path:    ":id(.*)",
					handler: mockHandler,
				},
				pathParams: map[string]string{
					"id": "123",
				},
			},
		},
		{
			// 命中 /:id([0-9]+)
			name:      ":id([0-9]+)",
			method:    http.MethodDelete,
			path:      "/123/home",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					path:    "home",
					handler: mockHandler,
				},
				pathParams: map[string]string{
					"id": "123",
				},
			},
		},
		{
			// 未命中
			name:   "not :id([0-9]+)",
			method: http.MethodDelete,
			path:   "/abc/home",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mi, found := r.findRoute(tc.method, tc.path)
			assert.Equal(t, tc.wantFound, found)
			if !found {
				return
			}

			assert.Equal(t, tc.info.pathParams, mi.pathParams)
			assert.Equal(t, tc.info.n.path, mi.n.path)
			wantVal := reflect.ValueOf(tc.info.n.handler)
			nVal := reflect.ValueOf(mi.n.handler)
			assert.Equal(t, nVal, wantVal)
		})
	}
}
