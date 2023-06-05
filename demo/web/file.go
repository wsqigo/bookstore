package web

import (
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type FileUploader struct {
	// FileField 对应于文件在表单中的字段名字
	FileField string
	// DstPathFunc 用于计算目标路径
	// 为什么要用户传？
	// 要考虑文件名冲突的问题
	// 所以很多时候，目标文件名字，都是随机的
	DstPathFunc func(fh *multipart.FileHeader) string
}

func (u FileUploader) Handle() HandleFunc {
	// 这里可以额外做一些检测，下面那种模式也可以做检测
	//if u.FileField == "" {
	//	// 这种方案默认值我其实不是很喜欢
	//	// 因为我们需要教会用户说，这个 file 是指什么意思
	//	u.FileField = "file"
	//}

	return func(ctx *Context) {
		// 上传文件的逻辑在这里
		// 第一步：读到文件内容
		file, fileHeader, err := ctx.Req.FormFile(u.FileField)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败" + err.Error())
			return
		}
		defer file.Close()
		// 第二步：计算出目标路径
		// 这种做法就是，将目标路径计算逻辑，交给用户
		// O_WRONLY 写入数据
		// O_TRUNC 如果文件本身存在，清空数据
		// O_CREATE 创建一个新的

		dstPath := u.DstPathFunc(fileHeader)
		// 可以尝试把 dst 上不存在的目录都全部建立起来
		//os.MkdirAll()
		dst, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败" + err.Error())
			return
		}
		defer dst.Close()
		// 第三步：保存文件
		// buf 会影响你的性能
		// 你要考虑复用
		_, err = io.CopyBuffer(dst, file, nil)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败" + err.Error())
			return
		}
		// 第四步：返回
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("上传成功")
	}
}

type FileUploaderOption func(loader *FileUploader)

func NewFileUploader(opts ...FileUploaderOption) *FileUploader {
	res := &FileUploader{
		FileField: "file",
		DstPathFunc: func(fh *multipart.FileHeader) string {
			return filepath.Join("testdata", "uploader", uuid.New().String())
		},
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

// HandleFunc 这种设计方案也是可以的，但是不如上一种灵活。
// 它可以直接用来注册路由
// 上一种可以在返回 HandleFunc 之前可以继续检测一下传下的字段
// 这种形态和 Option 模式配合就很好
func (u FileUploader) HandleFunc(ctx *Context) {
	// 文件传
}
