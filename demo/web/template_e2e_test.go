//go:build e2e

package web

import (
	"fmt"
	"html/template"
	"testing"
)

func TestLoginPage(t *testing.T) {
	tpl, err := template.ParseGlob("testdata/tpls/*.gohtml")
	if err != nil {
		t.Fatal(err)
	}
	engine := &GoTemplateEngine{
		T: tpl,
	}

	s := NewHTTPServer(ServerWithTemplateEngine(engine))
	s.Get("/login", func(ctx *Context) {
		err = ctx.Render("login.gohtml", nil)
		if err != nil {
			fmt.Println(err)
		}
	})
	s.Start(":8081")
}
