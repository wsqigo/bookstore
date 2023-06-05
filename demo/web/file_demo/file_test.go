package file_demo

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFile(t *testing.T) {
	f, err := os.Open("testdata/my_file.txt")
	assert.Nil(t, err)
	data := make([]byte, 64)
	n, err := f.Read(data)
	fmt.Println(n)
	assert.Nil(t, err)

	n, err = f.WriteString("hello")
	fmt.Println(n)
	// bad file descriptor 不可写
	assert.NotNil(t, err)
	f.Close()

	f, err = os.OpenFile("testdata/my_file.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	assert.Nil(t, err)
	n, err = f.WriteString("hello")
	fmt.Println(n)
	assert.Nil(t, err)
	f.Close()

	// 创建一个文件。如果文件已经存在，会被清空
	f, err = os.Create("testdata/my_file_copy.txt")
	assert.Nil(t, err)
	n, err = f.WriteString("hello, world")
	fmt.Println(n)
	assert.Nil(t, err)
	f.Close()
}
