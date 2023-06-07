package web

import (
	"github.com/google/uuid"
	lru "github.com/hashicorp/golang-lru"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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

type FileDownloader struct {
	Dir string
}

func (d FileDownloader) Handle() HandleFunc {
	return func(ctx *Context) {
		// 用的是 xxx?file=xxx
		req, err := ctx.QueryValue("file")
		if err != nil {
			ctx.RespStatusCode = http.StatusBadRequest
			ctx.RespData = []byte("找不到目标文件")
			return
		}
		path := filepath.Join(d.Dir, filepath.Clean(req))
		// 做一个校验，防止相对路径引起攻击者下载了你的系统文件
		// path, err = filepath.Abs(path)
		// if strings.Contains(path, d.Dir) {
		//
		// }
		fn := filepath.Base(path)
		header := ctx.Resp.Header()
		header.Set("Content-Disposition", "attachment;filename="+fn)
		header.Set("Content-Disposition", "File Transfer")
		header.Set("Content-Type", "application/octet-stream")
		header.Set("Content-Transfer-Encoding", "binary")
		// 下面两个是控制缓存的
		header.Set("Expires", "0")
		header.Set("Cache-Control", "must-revalidate")
		header.Set("Pragma", "public")

		http.ServeFile(ctx.Resp, ctx.Req, path)
	}
}

type StaticResourceHandlerOption func(h *StaticResourceHandler)

type StaticResourceHandler struct {
	dir               string
	extContextTypeMap map[string]string

	// 缓存静态资源的限制
	cache       *lru.Cache
	maxFileSize int
}

type fileCacheItem struct {
	fileName    string
	fileSize    int
	contentType string
	data        []byte
}

func NewStaticResourceHandler(dir string, opts ...StaticResourceHandlerOption) (*StaticResourceHandler, error) {
	res := &StaticResourceHandler{
		dir: dir,
		extContextTypeMap: map[string]string{
			// 这里根据自己的需要不断添加
			"jpeg": "image/jpeg",
			"jpe":  "image/jpeg",
			"jpg":  "image/jpeg",
			"png":  "image/png",
			"pdf":  "image/pdf",
		},
	}
	for _, opt := range opts {
		opt(res)
	}

	return res, nil
}

// WithFileCache 静态文件将会被缓存
// maxFileSizeThreshold 超过这个大小的文件，就被认为是大文件，我们将不会缓存
// 所以我们最多缓存 maxFileSizeThreshold * maxCacheFileCnt
func WithFileCache(maxFileSizeThreshold int, maxCacheFileCnt int) StaticResourceHandlerOption {
	return func(h *StaticResourceHandler) {
		c, err := lru.New(maxCacheFileCnt)
		if err != nil {
			log.Printf("创建缓存失败，将不会缓存静态资源")
		}
		h.maxFileSize = maxFileSizeThreshold
		h.cache = c
	}
}

func WithMoreExtension(extMap map[string]string) StaticResourceHandlerOption {
	return func(h *StaticResourceHandler) {
		for ext, contentType := range extMap {
			h.extContextTypeMap[ext] = contentType
		}
	}
}

func (h *StaticResourceHandler) Handle(ctx *Context) {
	// 1. 拿到目标文件名
	req, err := ctx.PathValue("file")
	if err != nil {
		ctx.RespStatusCode = http.StatusBadRequest
		ctx.RespData = []byte("请求路径不对")
		return
	}
	item, ok := h.readFileFromData(req)
	if ok {
		log.Printf("从缓存中读取数据...")
		h.writeItemAsResponse(item, ctx.Resp)
	}

	// 2. 定位到目标文件，并且读出来
	path := filepath.Join(h.dir, req)
	data, err := os.ReadFile(path)
	if err != nil {
		ctx.RespStatusCode = http.StatusInternalServerError
		ctx.RespData = []byte("服务器错误")
		return
	}
	// 3. 返回给前端
	ext := getFileExt(path)
	t, ok := h.extContextTypeMap[ext]
	if !ok {
		ctx.RespStatusCode = http.StatusBadRequest
		return
	}
	item = &fileCacheItem{
		fileSize:    len(data),
		data:        data,
		contentType: t,
		fileName:    req,
	}
	h.cacheFile(item)
	h.writeItemAsResponse(item, ctx.Resp)
}

func getFileExt(name string) string {
	ext := filepath.Ext(name)
	return ext[1:]
}

func (h *StaticResourceHandler) readFileFromData(fileName string) (*fileCacheItem, bool) {
	if h.cache != nil {
		if item, ok := h.cache.Get(fileName); ok {
			return item.(*fileCacheItem), true
		}
	}

	return nil, false
}

func (h *StaticResourceHandler) cacheFile(item *fileCacheItem) {
	if h.cache != nil && item.fileSize < h.maxFileSize {
		h.cache.Add(item.fileName, item)
	}
}

func (h *StaticResourceHandler) writeItemAsResponse(item *fileCacheItem, writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", item.contentType)
	writer.Header().Set("Content-Type", strconv.Itoa(item.fileSize))
	_, _ = writer.Write(item.data)
}
