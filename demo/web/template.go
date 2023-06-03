package web

import "context"

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
