package testcase

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt/v4"
)

// 用于签名的字段穿
var mySignKey = []byte("wsqigo")

// GenRegisteredClaims 使用默认声明创建 jwt
// 官方字段
// issuer 签发人
// expiration time 过期时间
// subject 主题
// audience 受众
// Not Before 生效时间
// Issued At 签发时间
// JWT ID 编号
func GenRegisteredClaims() (string, error) {
	// 创建 Claims
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // 过期时间
		Issuer:    "wsqigo",                                           // 签发人
	}

	// 生成 token 对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 生成签名字符串对付穿
	return token.SignedString(mySignKey)
}

// ParseRegisteredClaims 解析 jwt
func ParseRegisteredClaims(tokenString string) bool {
	// 解析 token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return mySignKey, nil
	})
	// 解析 token 失败
	if err != nil {
		return false
	}

	return token.Valid
}

// CustomClaims 自定义声明类型 并内嵌 jwt.RegisteredClaims
// jwt 包自带的 jwt.RegisteredClaims 只包含了官方字段
// 假设我们这里需要额外记录一个 username 字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体
type CustomClaims struct {
	// 可根据需要自行添加字段
	Username             string `json:"username"`
	jwt.RegisteredClaims        // 内嵌标准的声明
}

const TokenExpireDuration = time.Hour * 24

// CustomSecret 用于签名的字符串
var CustomSecret = []byte("夏天夏天悄悄过去")

// GenToken 生成 jwt
func GenToken(username string) (string, error) {
	// 创建一个我们自己的声明
	claims := CustomClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			Issuer:    "my-project", // 签发人
		},
	}

	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用指定的 secret 签名并获得完整的编码后的字符串 token
	return token.SignedString(CustomSecret)
}

// ParseToken 解析 JWT
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 解析 token
	// 如果是自定义 claim 结构体则需要使用 ParseWithClaims 方法
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 直接使用标准的 Claim 则可以直接使用 Parse 方法
		// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return CustomSecret, nil
	})
	if err != nil {
		return nil, err
	}

	// 对 token 对象进行断言
	// 校验 token
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func authHandler(ctx *gin.Context) {
	// 用户发送用户名和密码过来
	u := &user{}
	err := ctx.ShouldBind(u)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 2001,
			"msg":  "invalid param",
		})
		return
	}

	// 校验用户名和密码是否正确
	if u.Username == "q1mi" && u.Password == "q1mi123" {
		// 生成 Token
		tokenString, _ := GenToken(u.Username)
		ctx.JSON(http.StatusOK, gin.H{
			"code": 2000,
			"msg":  "success",
			"data": gin.H{"token": tokenString},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 2002,
		"msg":  "鉴权失败",
	})
}

func JWTAuthMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// 客户端携带 Token 有三种方式
		// 1. 放在请求头 2. 放在请求体 3. 放在URI
		// 这里假设 Token 放在 Header 的 Authorization 中，并使用 Bearer 开头
		// Authorization: Bearer xxxxxx.xxxxx
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusOK, gin.H{
				"code": 2003,
				"msg":  "请求头中 auth 为空",
			})
			ctx.Abort()
			return
		}

		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			ctx.JSON(http.StatusOK, gin.H{
				"code": 2004,
				"msg":  "无效的 Token",
			})
		}

		// parts[1] 是获取到的 tokenString，我们使用之前定义好的解析 JWT 的函数来解析它
		mc, err := ParseToken(parts[1])
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": 2005,
				"msg":  "无效的 Token",
			})
			return
		}

		// 将当前请求的 username 信息保存到请求的上下文 ctx 上
		ctx.Set("username", mc.Username)
		// 后续的处理函数可以通过 ctx.Get("username") 来获取当前请求的用户信息
		ctx.Next()
	}
}

type MyClaims struct {
	UserID int64

	jwt.RegisteredClaims
}

// GenTokenV2 生成 access token 和 refresh token
func GenTokenV2(userID int64) (string, string, error) {
	// 创建一个我们自己的声明
	claims := &MyClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			Issuer:    "bluebell",
		},
	}

	// 加密并获得完整的编码后的字符串 token
	aToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(mySignKey)
	if err != nil {
		return "", "", err
	}

	// refresh token 不需要存任何自定义数据
	rToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Hour * 24)),
		Issuer:    "bluebell",
	}).SignedString(mySignKey)

	return aToken, rToken, nil
}

// RefreshToken 刷新 AccessToken
func RefreshToken(aToken, rToken string) (string, string, error) {
	// refresh token 无效直接返回
	if _, err := jwt.Parse(rToken, func(token *jwt.Token) (interface{}, error) {
		return rToken, nil
	}); err != nil {
		return "", "", err
	}

	// 从就access token 中解析出claims数据
	claims := &MyClaims{}
	_, err := jwt.ParseWithClaims(aToken, claims, nil)
	v, _ := err.(*jwt.ValidationError)

	// 当 access token 是过期错误 并且 refresh token 没有过期时就创建一个新的 access token
	if v.Errors == jwt.ValidationErrorExpired {
		return GenTokenV2(claims.UserID)
	}

	return aToken, rToken, nil
}
