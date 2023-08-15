package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/wureny/webook/webook/Internal/web"
	"log"
	"net/http"
	"strings"
	"time"
)

type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 {
			// 没登录，有人瞎搞
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		claims := &web.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), nil
		})
		if err != nil {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid || claims.Uid == 0 {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claims.UserAgent != ctx.Request.UserAgent() {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		now := time.Now()
		// 每十秒钟刷新一次
		if claims.ExpiresAt.Sub(now) < time.Second*50 {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, err = token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
			if err != nil {
				// 记录日志
				log.Println("jwt 续约失败", err)
			}
			ctx.Header("x-jwt-token", tokenStr)
		}
		ctx.Set("claims", claims)
		ctx.Set("userId", claims.Uid)
	}
}
