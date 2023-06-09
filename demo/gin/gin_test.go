package gin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"path/filepath"
	"testing"
)

func TestGinFirst(t *testing.T) {
	// 创建一个默认的路由引擎
	r := gin.Default()
	// Get: 请求方式； /hello: 请求路径
	r.GET("/hello", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "Hello, world",
		})
	})
	r.Run(":8080")
}

func TestTemplate(t *testing.T) {
	r := gin.Default()
	r.LoadHTMLGlob("templates/**/*.gohtml")
	//r.LoadHTMLFiles("templates/posts/index.gohtml", "templates/users/index.gohtml")
	r.GET("/posts/index", func(context *gin.Context) {
		context.HTML(http.StatusOK, "posts/index.gohtml", gin.H{
			"title": "posts/index",
		})
	})

	r.GET("users/index", func(context *gin.Context) {
		context.HTML(http.StatusOK, "users/index.gohtml", gin.H{
			"title": "users/index",
		})
	})

	r.Run(":8080")
}

func TestQueryString(t *testing.T) {
	r := gin.Default()
	r.GET("/user/search", func(context *gin.Context) {
		username := context.DefaultQuery("username", "小王子")
		// username := context.Query("username")
		address := context.Query("address")
		// 输出 json 结果给调用方
		context.JSON(http.StatusOK, gin.H{
			"message":  "ok",
			"username": username,
			"address":  address,
		})
	})
	r.Run(":8080")
}

func TestFormData(t *testing.T) {
	r := gin.Default()
	r.POST("/user/search", func(context *gin.Context) {
		// DefaultPostForm取不到值时会返回指定的默认值
		// context.DefaultPostForm("username", "小王子")
		username := context.PostForm("username")
		address := context.PostForm("address")
		// 输出json结果给调用方
		context.JSON(http.StatusOK, gin.H{
			"message":  "ok",
			"username": username,
			"address":  address,
		})
	})
	r.Run(":8080")
}

func TestParamPath(t *testing.T) {
	r := gin.Default()
	r.GET("/user/search/:username/:address", func(context *gin.Context) {
		username := context.Param("username")
		address := context.Param("address")
		// 输出json结果给调用方
		context.JSON(http.StatusOK, gin.H{
			"message":  "ok",
			"username": username,
			"address":  address,
		})
	})
	r.Run(":8080")
}

// Binding for JSON
type Login struct {
	User     string `form:"user" json:"user" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// ShouldBind 会按照下面的顺序解析请求中的数据完成绑定:
// 1. 如果是 Get 请求，只使用 Form 绑定引擎(query)
// 2.
func TestShouldBind(t *testing.T) {
	router := gin.Default()

	// 绑定JSON的实例 ({"user": "q1mi", "password": "123456"})
	router.POST("/loginJSON", func(context *gin.Context) {
		login := &Login{}

		// ShouldBind() 会根据请求的 Content-Type 自行选择绑定器
		err := context.ShouldBind(login)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Printf("login info:%#v", login)
		context.JSON(http.StatusOK, gin.H{
			"user":     login.User,
			"password": login.Password,
		})
	})

	// 绑定form表单示例 (user=q1mi&password=123456)
	router.POST("/loginForm", func(context *gin.Context) {
		login := &Login{}

		// ShouldBind() 会根据请求的 Content-Type 自行选择绑定器
		err := context.ShouldBind(login)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Printf("login info:%#v", login)
		context.JSON(http.StatusOK, gin.H{
			"user":     login.User,
			"password": login.Password,
		})
	})

	// 绑定QueryString示例 (loginQuery?user=q1mi&password=123456)
	router.GET("/loginForm", func(context *gin.Context) {
		login := &Login{}

		// ShouldBind() 会根据请求的 Content-Type 自行选择绑定器
		err := context.ShouldBind(login)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Printf("login info:%#v", login)
		context.JSON(http.StatusOK, gin.H{
			"user":     login.User,
			"password": login.Password,
		})
	})
	router.Run(":8080")
}

func TestUploadFile(t *testing.T) {
	router := gin.Default()

	router.LoadHTMLGlob("templates/**/*.gohtml")
	router.GET("/upload", func(context *gin.Context) {
		context.HTML(http.StatusOK, "uploads/index.gohtml", nil)
	})

	// 处理multipart forms提交文件时默认的内存限制是32 MiB
	// 可以通过下面的方式修改
	// router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.POST("/upload", func(context *gin.Context) {
		// 单个文件
		file, err := context.FormFile("f1")
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		log.Println(file.Filename)
		dst := filepath.Join("templates", "uploads", file.Filename)
		// 上传文件到指定路径
		context.SaveUploadedFile(file, dst)
		context.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("'%s' uploaded!", file.Filename),
		})
	})

	router.Run()
}
