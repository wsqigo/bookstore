package iris

import (
	"github.com/kataras/iris/v12"
	"testing"
)

func TestIrisController(t *testing.T) {
	// 相当于Beego的HttpServer和Gin的Engine
	app := iris.New()

	app.Get("/", func(ctx iris.Context) {
		_, _ = ctx.HTML("Hello <strong>%s</strong>!", "World")
	})

	_ = app.Listen(":8083")
}
