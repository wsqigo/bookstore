package template_demo

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"html/template"
	"testing"
)

const serviceTpl = `
{{- $service := .GenName -}}
type {{ $service }} struct {
	EndPoint string
	Path string
	Client http.Client
}
`

func TestHelloWorld(t *testing.T) {
	type User struct {
		Name string
	}
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`Hello, {{.Name}}`)
	if err != nil {
		t.Fatal(err)
	}
	bs := &bytes.Buffer{}
	err = tpl.Execute(bs, &User{Name: "Tom"})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Hello, Tom", bs.String())
}

func TestMapData(t *testing.T) {
	tpl := template.New("map-data")
	tpl, err := tpl.Parse(`Hello, {{.Name}}`)
	if err != nil {
		t.Fatal(err)
	}
	bs := &bytes.Buffer{}
	err = tpl.Execute(bs, map[string]string{"Name": "Jerry"})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Hello, Jerry", bs.String())
}

func TestSliceData(t *testing.T) {
	tpl := template.New("slice-data")
	tpl, err := tpl.Parse(`Hello, {{index . 0}}`)
	if err != nil {
		t.Fatal(err)
	}

	bs := &bytes.Buffer{}
	err = tpl.Execute(bs, []string{"John"})
	if err != nil {
		t.Fatal()
	}
	assert.Equal(t, "Hello, John", bs.String())
}

func TestBasicData(t *testing.T) {
	tpl := template.New("slice-data")
	tpl, err := tpl.Parse(`Hello, {{.}}`)
	if err != nil {
		t.Fatal(err)
	}

	bs := &bytes.Buffer{}
	err = tpl.Execute(bs, 123)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Hello, 123", bs.String())
}

type FuncCall struct {
	Slice []string
}

func (f FuncCall) Hello(firstName string, lastName string) string {
	return fmt.Sprintf("Hello, %s·%s", firstName, lastName)
}

func TestFuncCall(t *testing.T) {
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`
切片长度: {{len .Slice}}
say Hello: {{.Hello "Tom" "Jerry"}}
打印数字: {{printf "%.2f" 1.234}} 
`)
	assert.Nil(t, err)

	bs := &bytes.Buffer{}
	err = tpl.Execute(bs,
		&FuncCall{Slice: []string{"Tom", "Jerry"}})
	assert.Nil(t, err)
	assert.Equal(t, `
切片长度: 2
say Hello: Hello, Tom Jerry
打印数字: 1.23
`, bs.String())
}

func TestForLoop(t *testing.T) {
	// 用一点小技巧来实现 for i 循环
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`
{{- range $idx, $elem := . -}}
下标: {{$idx -}}
{{- end}}
`)
	assert.Nil(t, err)
	bs := &bytes.Buffer{}
	// 假设我们要从 0 迭代到 100，即 [0 100]
	// 这里的切片可以是任意类型，[]bool, []byte 都可以
	// 因为我们本身并不关心里面元素，只是借用一下下标而已
	data := make([]bool, 100)
	err = tpl.Execute(bs, data)
	assert.Nil(t, err)
	t.Log(bs.String())
}

func TestIfElseBlock(t *testing.T) {
	type User struct {
		Age int
	}
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`
{{- if and (gt .Age 0) (le .Age 6) -}} 
儿童 (0, 6]
{{- else if and (gt .Age 6) (le .Age 18) -}}
少年 (6, 18]
{{- else -}}
成人 > 18
{{- end -}}
`)
	assert.Nil(t, err)
	bs := &bytes.Buffer{}
	err = tpl.Execute(bs, User{Age: 4})
	assert.Nil(t, err)
	assert.Equal(t, `儿童 (0, 6]`, bs.String())
}
