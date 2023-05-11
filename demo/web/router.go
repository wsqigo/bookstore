package web

import (
	"fmt"
	"strings"
)

/*
路由树归根结底就是设计一颗多叉树
*/

// 用来支持对路由树的操作
// 维持住了所有的路由树，它是整个路由注册和查找的总入口。代表路由树（森林）
type router struct {
	// Beego Gin HTTP method 对应一棵树
	// GET 有一棵树，POST 也有一棵树

	// trees 是按照HTTP方法来组织的
	// http method => 路由树根节点
	trees map[string]*node
}

func newRouter() router {
	return router{
		trees: map[string]*node{},
	}
}

// addRoute 注册路由
// method 是 HTTP 方法
// - 已经注册了的路由，无法被覆盖。例如 /user/home注册两次，会冲突
// - path 必须以 / 开头并且结尾不能有 /，中间也不允许有连续的 /
// - 不能在同一个位置注册不同的参数路由，例如 /user/:id 和 /user/:name 冲突
// 定义为私有的 addRoute
// 1. 用户只能通过 Get 或者 Post来注册，那么可以确保 method 参数永远是对的
// 2. addRoute 在接口里面是私有的，限制了用户将无法实现 Server。
func (r *router) addRoute(method string, path string, handleFunc HandleFunc) {
	// 对 path 进行校验
	if path == "" {
		panic("web: 路由是空字符串")
	}

	if path[0] != '/' {
		panic("web: 路径必须以 / 开头")
	}

	if path != "/" && path[len(path)-1] == '/' {
		panic("web: 路径不能以 / 结尾")
	}
	// 首先找到树来
	root, ok := r.trees[method]
	// 这是一个全新的 HTTP 方法，创建根节点
	if !ok {
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}

	// 根节点特殊处理一下
	if path == "/" {
		if root.handler != nil {
			panic("web: 路由冲突[/]")
		}
		root.handler = handleFunc
		return
	}

	// 切割这个 path
	segs := strings.Split(path[1:], "/")
	// 开始一段段处理
	for _, seg := range segs {
		if seg == "" {
			panic(fmt.Sprintf("web: 非法路由。不允许使用 //a/b, /a//b 之类的路由,[%s]", path))
		}
		// 递归下去，找准位置
		// 如果中途有节点不存在，你就要创建出来
		root = root.childOrCreate(seg)
	}
	if root.handler != nil {
		panic(fmt.Sprintf("web: 路由冲突，重复注册[%s]", path))
	}
	root.handler = handleFunc
}

// findRoute 查找对应的节点
// 注意，返回的 node 内部 handleFunc 不为nil才算是注册了路由
func (r *router) findRoute(method string, path string) (*matchInfo, bool) {
	// 基本上是不是也是沿着树深度查找下去？
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}

	if path == "/" {
		return &matchInfo{n: root}, true
	}

	// 这里把前置和后置的 / 都去掉，然后按照斜杠切割
	segs := strings.Split(strings.Trim(path, "/"), "/")
	mi := &matchInfo{}
	for _, seg := range segs {
		var matchParam bool
		root, matchParam, ok = root.childOf(seg)
		if !ok {
			return nil, false
		}
		if matchParam {
			mi.addValue(root.path[1:], seg)
		}
	}

	mi.n = root
	return mi, true
}

// node 代表路由树的节点
// 路由树的匹配顺序是：
// 1. 静态完全匹配
// 2. 路径参数匹配：形式: param_name
// 3. 通配符匹配：*
// 这是不回溯匹配
type node struct {
	path string

	// 静态匹配的节点
	// 子 path 到子节点的映射 子path => 子node
	children map[string]*node

	// handler 命中路由之后执行的逻辑
	handler HandleFunc

	// 路径参数
	paramChild *node

	// 通配符 * 表达的节点，任意匹配
	starChild *node
}

// childOrCreate 查找子节点
// 首先会判断 path 是不是通配符匹配
// 其次判断 path 是不是参数路径，即以 : 开头的路径
// 最后会从 children 里面查找，
// 如果没有找到，那么会创建一个新的节点，并且保存在 node 里面
func (n *node) childOrCreate(seg string) *node {
	childNode := &node{
		path: seg,
	}

	if seg == "*" {
		if n.paramChild != nil {
			panic(fmt.Sprintf("web: 非法路由，已有路径参数路由。不允许同时注册通配符路由和参数路由 [%s]", seg))
		}
		if n.starChild == nil {
			n.starChild = childNode
		}
		return n.starChild
	}

	// 以 : 开头，我们认为是参数路由
	if seg[0] == ':' {
		if n.starChild != nil {
			panic(fmt.Sprintf("web: 非法路由，已有通配符路由。不允许同时注册通配符路由和参数路由 [%s]", seg))
		}
		if n.paramChild == nil {
			n.paramChild = childNode
		} else {
			if n.paramChild.path != seg {
				panic(fmt.Sprintf("web: 路由冲突，参数路由冲突，已有 %s，新注册 %s", n.paramChild.path, seg))
			}
		}
		return n.paramChild
	}

	if n.children == nil {
		n.children = make(map[string]*node)
	}
	res, ok := n.children[seg]
	if !ok {
		// 要新建一个
		res = childNode
		n.children[seg] = res
	}

	return res
}

// childOf 返回子节点
// 第一个返回值 *node 是命中的节点
// 第二个返回值 bool 代表是否是命中参数路由
// 第三个返回值 bool 代表是否命中
func (n *node) childOf(path string) (*node, bool, bool) {
	if n.children == nil {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.starChild, false, n.starChild != nil
	}
	child, found := n.children[path]
	if !found {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.starChild, false, n.starChild != nil
	}
	return child, false, found
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
}

func (m *matchInfo) addValue(key string, value string) {
	if m.pathParams == nil {
		// 大多数情况，参数路径只会有一段
		m.pathParams = make(map[string]string, 1)
	}
	m.pathParams[key] = value
}
