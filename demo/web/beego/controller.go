package beego

import "github.com/beego/beego/v2/server/web"

type UserController struct {
	web.Controller
}

func (c *UserController) GetUser() {
	c.Ctx.WriteString("你好，我是wsqigo")
}

type User struct {
	Name string
}

// CreateUser 用户操作请求和响应是通过 Ctx 来达成的。
func (c *UserController) CreateUser() {
	u := &User{}

	err := c.Ctx.BindJSON(u)
	if err != nil {
		c.Ctx.WriteString(err.Error())
		return
	}

	_ = c.Ctx.JSONResp(u)
}
