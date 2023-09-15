package jwt_token

import (
	"gin-demo/core/settings"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"time"
)

// AuthCodeInfo 自定义JWT额外字段信息
type AuthCodeInfo struct {
	UserId   string `json:"userId"`
	UserName string `json:"username"`
}

// AuthCodeJwtClaims 自定义AuthCode JWT信息
type AuthCodeJwtClaims struct {
	jwt.RegisteredClaims
	AuthCodeInfo
}

var secret = settings.Config.AuthCodeJwt.Secret

// GenerateToken 生成JWT
func GenerateToken(authCodeInfo AuthCodeInfo, expire time.Duration) (string, error) {
	var (
		issuer = settings.Config.AuthCodeJwt.Issuer
		token  string
		err    error
	)
	claims := AuthCodeJwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),             //发行时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)), //过期时间
			NotBefore: jwt.NewNumericDate(time.Now()),             // 生效时间
			Issuer:    issuer,                                     //指定token发行人
		},
		AuthCodeInfo: authCodeInfo,
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tokenClaims.SignedString([]byte(secret))
	return token, err
}

func Secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil // 这是我的secret
	}
}

func ParseToken(tokenStr string) (*AuthCodeJwtClaims, error) {
	//用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*Token
	token, err := jwt.ParseWithClaims(tokenStr, &AuthCodeJwtClaims{}, Secret())
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("that's not even a token")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("token is expired")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("token not active yet")
			} else {
				return nil, errors.New("couldn't handle this token")
			}
		}
		return nil, err
	}
	if claims, ok := token.Claims.(*AuthCodeJwtClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("couldn't handle this token")
}
