package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
	"time"
)

var Jwt *JWT

const SigningKey = "test"

type JWT struct {
	SigningKey []byte
}

type CustomClaims struct {
	BaseClaims
	BufferTime int64
	jwt.RegisteredClaims
}

type BaseClaims struct {
	Uid      int64
	Username string
}

//var (
//	TokenExpired     = errors.New("Token is expired")
//	TokenNotValidYet = errors.New("Token not active yet")
//	TokenMalformed   = errors.New("That's not even a token")
//	TokenInvalid     = errors.New("Couldn't handle this token:")
//)

func InitJWT() {
	Jwt = &JWT{
		[]byte(SigningKey),
	}
}

func AuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := strings.TrimPrefix(ctx.GetHeader("Authorization"), "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) { return SigningKey, nil })
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": -1, "msg": fmt.Sprintf("access token parse error: %v", err)})
			return
		}
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			if !claims.VerifyExpiresAt(time.Now(), false) {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": -1, "msg": "access token expired"})
				return
			}
			ctx.Set("claims", claims)
		} else {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": -1, "msg": fmt.Sprintf("Claims parse error: %v", err)})
			return
		}
		ctx.Next()
	}
}

func (j *JWT) CreateClaims(baseClaims BaseClaims) CustomClaims {
	claims := CustomClaims{
		BaseClaims: baseClaims,
		BufferTime: 3600 * 24, // 缓冲时间1天 缓冲时间内会获得新的token刷新令牌 此时一个用户会存在两个有效令牌 但是前端只留一个 另一个会丢失
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1000 * time.Second)), // 签名生效时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 7 * time.Hour)),  // 过期时间 7天  配置文件
			Issuer:    "test",                                                  // 签名的发行者
		},
	}
	return claims
}

// CreateToken 创建一个token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

//// CreateTokenByOldToken 旧token 换新token 使用归并回源避免并发问题
//func (j *JWT) CreateTokenByOldToken(oldToken string, claims CustomClaims) (string, error) {
//	//todo SingleFlight
//	//v, err, _ := global.GVA_Concurrency_Control.Do("JWT:"+oldToken, func() (interface{}, error) {
//	return j.CreateToken(claims)
//	//})
//	//return v.(string), err
//}
//
//// ParseToken 解析 token
//func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
//	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
//		return j.SigningKey, nil
//	})
//	if err != nil {
//		if ve, ok := err.(*jwt.ValidationError); ok {
//			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
//				return nil, TokenMalformed
//			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
//				// Token is expired
//				return nil, TokenExpired
//			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
//				return nil, TokenNotValidYet
//			} else {
//				return nil, TokenInvalid
//			}
//		}
//	}
//	if token != nil {
//		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
//			return claims, nil
//		}
//		return nil, TokenInvalid
//
//	} else {
//		return nil, TokenInvalid
//	}
//}
