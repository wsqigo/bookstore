package jwt

import (
	"errors"
	"time"

	"github.com/spf13/viper"

	"github.com/golang-jwt/jwt/v4"
)

const jwtExpKey = "auth.jwt_expire"

var MySecret = []byte("夏天夏天悄悄过去")

// MyClaims 自定义声明结构体并内嵌 jwt.RegisteredClaims
// jwt 包自带的 jwt.RegisteredClaims 只包含了官方字段
// 假设我们这里需要额外记录一个 username 字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体
type MyClaims struct {
	UserID   int64  `json:"user_id"`
	UserName string `json:"user_name"`
	jwt.RegisteredClaims
}

// GenToken 生成 JWT
func GenToken(userID int64, username string) (string, error) {
	// 创建一个我们自己的声明
	claims := MyClaims{
		UserID:   userID,
		UserName: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(viper.GetDuration(jwtExpKey))), // 过期时间
			Issuer:    "bluebell",
		},
	}

	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用指定的 secret 签名并获得完整的编码后的字符串 token
	return token.SignedString(MySecret)
}

// ParseToken 解析 JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析 token
	mc := &MyClaims{}
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}

	// 校验 token
	if token.Valid {
		return mc, nil
	}

	return nil, errors.New("invalid token")
}
