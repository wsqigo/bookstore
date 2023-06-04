package web

import (
	"bytes"
	"context"
	"html/template"
)

type TemplateEngine interface {
	// Render 渲染页面
	// tplName 模板的名字，按名索引
	// data 渲染页面所需要的数据
	Render(ctx context.Context, tplName string, data any) ([]byte, error)

	// 渲染页面，数据写入到 Writer 里面
	// Render(ctx, "aa", map[]{}, responseWriter)
	// Render(ctx context.Context, tplName string, data any, writer io.Writer) error
	// 不需要，让具体实现自己去管自己的模板
	// 加入或者删除一个模板，是具体模板引擎的事情，和 Web 框架什么关系都没有
	// AddTemplate(tplName string, tpl []byte) error
}

type GoTemplateEngine struct {
	T *template.Template
	// 也可以考虑设计为 map[string]*template.Template
	// 但是其实没太大必要，因为 template.Template 本身就提供了按名索引的功能
}

func (e *GoTemplateEngine) Render(ctx context.Context,
	tplName string, data any) ([]byte, error) {
	bs := &bytes.Buffer{}
	err := e.T.ExecuteTemplate(bs, tplName, data)
	return bs.Bytes(), err
}

// 以下这三个方法，是管理模板本身的方法。
// Web 框架根本不在意你从哪里把模板搞到，它只关心 Render 方法要实现，所以说管
// 理模板的方法并不算是 TemplateEngine 接口的一部分。
// 你封装了，你就得教会用户怎么用，例如文件怎么定位......

func (e *GoTemplateEngine) LoadFromGlob(pattern string) error {
	var err error
	e.T, err = template.ParseGlob(pattern)
	return err
}

func (e *GoTemplateEngine) LoadFromFiles(filenames ...string) error {
	var err error
	e.T, err = template.ParseFiles(filenames...)
	return err
}
