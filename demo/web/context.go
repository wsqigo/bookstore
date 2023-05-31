package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

/*
为什么都已经有了 http 包，还要开发 Web 框架?
高级路由功能，封装HTTP上下文以提供简单API、封装 Server 以提供生命周期控制、设计插件机制以提供无侵入式解决方案，提供如上传下载等默认功能
*/

type Context struct {
	Req *http.Request
	// Resp 原生的 ResponseWriter。当你直接使用 Resp 的时候，
	// 那么相当于你绕开了 RespData 和 RespStatusCode。
	// 响应数据直接被发送到前端，其它中间件将无法修改响应
	// 其实我们也可以考虑将这个做成私有的
	Resp http.ResponseWriter

	// 缓存的响应部分
	// 这部分数据会在最后刷新到前端
	RespData       []byte
	RespStatusCode int

	PathParams map[string]string

	//Ctx context.Context

	// 缓存的数据
	cacheQueryValues url.Values
	// 命中的路由
	MatchedRoute string
}

func (c *Context) SetCookie(cookie *http.Cookie) {
	// 不推荐
	// cookie.SameSite = c.cookieSameSite
	http.SetCookie(c.Resp, cookie)
}

func (c *Context) RespJSON(statusCode int, val any) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}

	//c.Resp.Header().Set("Content-Type", "application/json")
	c.RespStatusCode = statusCode
	c.RespData = data
	return err
}

func (c *Context) RespJSONOK(val any) error {
	return c.RespJSON(http.StatusOK, val)
}

func (c *Context) BindJSON(val any) error {
	if c.Req.Body == nil {
		return errors.New("web: body 为 nil")
	}

	decoder := json.NewDecoder(c.Req.Body)
	return decoder.Decode(val)
}

//func (c *Context) BindJSONOpt(val any, useNumber bool, disableUnknown bool) error {
//	if c.Req.Body == nil {
//		return errors.New("web: body 为 nil")
//	}
//
//	decoder := json.NewDecoder(c.Req.Body)
//	// useNumber => 数字就是用 Number 来表示
//	// 否则默认是 float64
//	// decoder.UseNumber()
//	if useNumber {
//		decoder.UseNumber()
//	}
//	// 如果要是有一个未知的字段，就会报错
//	// 比如说你 User 只有 Name 和 Email 两个字段
//	// JSON 里面额外多了一个 Age 字段，那么就会报错
//	if disableUnknown {
//		decoder.DisallowUnknownFields()
//	}
//	return decoder.Decode(val)
//}

func (c *Context) FormValue(key string) (string, error) {
	if err := c.Req.ParseForm(); err != nil {
		return "", err
	}

	return c.Req.FormValue(key), nil
}

func (c *Context) FormValueV2(key string) StringValue {
	if err := c.Req.ParseForm(); err != nil {
		return StringValue{err: err}
	}

	return StringValue{val: c.Req.FormValue(key)}
}

//func (c *Context) FormValueAsInt64(key string) (int64, error) {
//	val, err := c.FormValue(key)
//	if err != nil {
//		return 0, err
//	}
//
//	return strconv.ParseInt(val, 10, 64)
//}

// QueryValue Query 和表单比起来，它没有缓存
func (c *Context) QueryValue(key string) (string, error) {
	// 用户区别不出来是真的有值，但是值恰好是空字符串
	// 还是没有值
	//return c.Req.URL.Query().Get(key), nil
	if c.cacheQueryValues == nil {
		c.cacheQueryValues = c.Req.URL.Query()
	}

	vals, ok := c.cacheQueryValues[key]
	if !ok {
		return "", errors.New("web: 找不到这个 key")
	}
	return vals[0], nil
}

func (c *Context) QueryValueV2(key string) StringValue {
	// 用户区别不出来是真的有值，但是值恰好是空字符串
	// 还是没有值
	//return c.Req.URL.Query().Get(key), nil
	if c.cacheQueryValues == nil {
		c.cacheQueryValues = c.Req.URL.Query()
	}

	vals, ok := c.cacheQueryValues[key]
	if !ok {
		return StringValue{err: errors.New("web: 找不到这个 key")}
	}
	return StringValue{val: vals[0]}
}

func (c *Context) PathValue(key string) (string, error) {
	val, ok := c.PathParams[key]
	if !ok {
		return "", errors.New("web: 找不到这个 key")
	}
	return val, nil
}

func (c *Context) PathValueV2(key string) StringValue {
	val, ok := c.PathParams[key]
	if !ok {
		return StringValue{err: errors.New("web: 找不到这个 key")}
	}
	return StringValue{val: val}
}

type StringValue struct {
	val string
	err error
}

func (s StringValue) String() (string, error) {
	return s.val, s.err
}

func (s StringValue) ToInt64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}

	return strconv.ParseInt(s.val, 10, 64)
}
