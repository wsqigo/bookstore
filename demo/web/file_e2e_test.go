//go:build e2e

package web

import (
	"html/template"
	"mime/multipart"
	"path/filepath"
	"testing"
)

func TestUpload(t *testing.T) {
	tpl, err := template.ParseGlob("testdata/tpls/*.gohtml")
	if err != nil {
		t.Fatal(err)
	}
	engine := &GoTemplateEngine{
		T: tpl,
	}

	s := NewHTTPServer(ServerWithTemplateEngine(engine))
	s.Get("/upload", func(ctx *Context) {
		err = ctx.Render("upload.gohtml", nil)
		if err != nil {
			t.Fatal(err)
		}
	})

	fu := FileUploader{
		// 这里的 myfile 就是 <input type="file" name="myfile" />
		// 那个 name 的取值
		FileField: "myfile",
		DstPathFunc: func(fh *multipart.FileHeader) string {
			return filepath.Join("testdata", "upload", fh.Filename)
		},
	}
	s.Post("/upload", fu.Handle())
	s.Start(":8081")
}

func TestDownload(t *testing.T) {
	s := NewHTTPServer()

	fd := FileDownloader{
		Dir: filepath.Join("testdata", "download"),
	}
	s.Get("/download", fd.Handle())
	s.Start(":8081")
}

func TestStaticResourceHandler_Handle(t *testing.T) {
	s := NewHTTPServer()
	h, err := NewStaticResourceHandler(filepath.Join("testdata", "static"))
	if err != nil {
		t.Fatal(err)
	}

	// /static/js/:file

	// localhost:8081/static/xxx.jpg
	s.Get("/static/:file", h.Handle)
	s.Start(":8081")
}
