package gin

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
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

// StatCost 是一个统计请求耗时的中间件
func StatCost() gin.HandlerFunc {
	return func(context *gin.Context) {
		start := time.Now()
		context.Set("name", "小王子") // 可以通过 context.Set 在请求中设置值，后续的处理函数能够取到该值
		// 调用该请求的剩余处理程序
		context.Next()
		// 不调用该请求的剩余处理程序
		// context.Abort()
		// 计算耗时
		cost := time.Since(start)
		log.Println(cost)
	}
}

// 记录响应体的中间件
type bodyLogWriter struct {
	gin.ResponseWriter // 嵌入 gin 框架 ResponseWriter

	body *bytes.Buffer // 我们记录用的 response
}

func (w bodyLogWriter) Write(data []byte) (int, error) {
	w.body.Write(data)                  //我们记录一分
	return w.ResponseWriter.Write(data) // 真正写入响应
}

// ginBodyLogMiddleware 一个记录返回给客户端响应体的中间件
// https://stackoverflow.com/questions/38501325/how-to-log-response-body-in-gin
func ginBodyLogMiddleware(context *gin.Context) {
	blw := &bodyLogWriter{
		body:           bytes.NewBuffer([]byte{}),
		ResponseWriter: context.Writer,
	}
	context.Writer = blw // 使用我们自定义的类型替换默认的

	context.Next() // 执行业务逻辑

	fmt.Println("Response body: " + blw.body.String()) // 事后按需记录返回的响应
}

func TestMiddleware(t *testing.T) {
	// 新建一个没有任何默认中间件的路由
	r := gin.New()
	// 注册一个全局中间件
	r.Use(GinLogger())

	r.GET("/test", func(context *gin.Context) {
		name := context.MustGet("name").(string) //从上下文取值
		log.Println(name)
		context.JSON(http.StatusOK, gin.H{
			"message": "Hello world!",
		})
	})

	// 给 /test2 路由单独注册中间件（可注册多个）
	//r.GET("/test2", StatCost(), func(context *gin.Context) {
	//	name := context.MustGet("name")
	//	log.Println(name)
	//	context.JSON(http.StatusOK, gin.H{
	//		"message": "Hello world!",
	//	})
	//})

	r.Run()
}

// gin 的中间件实现，函数责任链条
//type HandlerFunc func(*Context)

//type HandlersChain []HandlerFunc

//func (c *Context) Next() {
//	c.index++
//	for c.index < int8(len(c.handlers)) {
//		c.handlers[c.index](c)
//		c.index++
//	}
//}
