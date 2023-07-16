package testcase

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/go-playground/locales/en"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

type User struct {
	Name string `validate:"min=6,max=10"`
	Age  int    `validate:"min=1,max=100"`
}

// 校验的时候三步
// 1. 调用 validator.New() 初始化一个校验器；
// 2. 将【待校验的结构体】传入我们的校验器的 Struct 方法中；
// 3. 校验返回的 error 是否为 nil 即可。
func TestValidator(t *testing.T) {
	validate := validator.New()

	u1 := User{Name: "lidajun", Age: 18}
	err := validate.Struct(u1)
	fmt.Println(err)

	u2 := User{Name: "dj", Age: 101}
	err = validate.Struct(u2)
	fmt.Println(err)
}

// User2 contains user infomation
type User2 struct {
	FirstName      string     `validate:"required"`
	LastName       string     `validate:"required"`
	Age            uint8      `validate:"gte=0,lte=130"`
	Email          string     `validate:"required,email"`
	FavouriteColor string     `validate:"iscolor"`                // alias for 'hexcolor|rgb|rgba|hsl|hsla'
	Addresses      []*Address `validate:"required,dive,required"` // a person can have a home and cottage...
}

// Address houses a users address information
type Address struct {
	Street string `validate:"required"`
	City   string `validate:"required"`
	Planet string `validate:"required"`
	Phone  string `validate:"required"`
}

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func TestValidatorStruct(t *testing.T) {
	validate = validator.New()

}

func validateStruct() {
	address := &Address{
		Street: "Eavesdown Docks",
		Planet: "Persphone",
		Phone:  "none",
	}

	user := &User2{
		FirstName:      "Badger",
		LastName:       "Smith",
		Age:            135,
		Email:          "Badger.Smith@gmail.com",
		FavouriteColor: "#000-",
		Addresses:      []*Address{address},
	}

	// returns nil or ValidatorErrors ( []FieldError )
	err := validate.Struct(user)
	if err != nil {
		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this
	}
}

type SignUpParam struct {
	Age        uint8  `json:"age" binding:"gte=1,lte=130"`
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

func TestBind(t *testing.T) {
	r := gin.Default()

	r.POST("/signup", func(ctx *gin.Context) {
		u := &SignUpParam{}
		if err := ctx.ShouldBind(u); err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
			return
		}

		// 保存入库等业务逻辑代码...

		ctx.JSON(http.StatusOK, "success")
	})

	_ = r.Run(":8999")
}

// 定义一个全局翻译器
var trans ut.Translator

// InitTrans 初始化翻译器
func InitTrans(locale string) error {
	// 修改 gin 框架中的 Validator 引擎属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		zhT := zh.New() // 中文翻译器
		enT := en.New() // 英文翻译器

		// 第一个参数是备用（fallback)的语言环境
		// 后面的参数是应该支持的语言环境（支持多个）
		// uni := ut.New(zhT, zhT) 也是可以的
		uni := ut.New(enT, zhT, enT)

		// locale 通常取决于 http 请求头的 'Accept-Language'
		var ok bool
		// 也可以使用 uni.FindTranslator(...) 传入多个 locale 进行查找
		trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}

		// 注册翻译器
		var err error
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		}

		return err
	}
	return nil
}

func TestTranslator(t *testing.T) {
	if err := InitTrans("zh"); err != nil {
		fmt.Println("init trans failed, err:", err)
		return
	}

	u := &SignUpParam{}
	r := gin.Default()

	r.POST("/signup", func(ctx *gin.Context) {
		if err := ctx.ShouldBind(&u); err != nil {
			// 获取 validator.ValidationErrors 类型的 errors
			errs, ok := err.(validator.ValidationErrors)
			if !ok {
				// 非 validator.ValidationErrors 类型错误直接返回
				ctx.JSON(http.StatusOK, gin.H{
					"msg": err.Error(),
				})
				return
			}

			// validator.ValidationErrors 类型错误则进行翻译
			ctx.JSON(http.StatusOK, gin.H{
				"msg": errs.Translate(trans),
			})
			return
		}
		// 保存入库等具体业务逻辑代码...

		ctx.JSON(http.StatusOK, "success")
	})

	_ = r.Run(":8999")
}
