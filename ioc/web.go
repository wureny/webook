package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/wureny/webook/webook/Internal/web"
	"github.com/wureny/webook/webook/Internal/web/jwt"
	"github.com/wureny/webook/webook/Internal/web/middleware"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoute(server)
	return server
}
func InitMiddlewares(redisclient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middleware.NewLoginJWTMiddlewareBuilder(jwt.NewRedisJWTHandler(redisclient)).
			IgnorePaths("/users/signup").
			IgnorePaths("/hello").
			IgnorePaths("/login_sms/code/send").
			IgnorePaths("/login_sms").
			IgnorePaths("/users/loginJWT").Build(),
		//	ratelimit.NewBuilder(redisclient, time.Second, 100).Build(),
	}
}

func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		//AllowOrigins: []string{"*"},
		//AllowMethods: []string{"POST", "GET"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		// 你不加这个，前端是拿不到的
		ExposeHeaders: []string{"x-jwt-token"},
		// 是否允许你带 cookie 之类的东西
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				// 你的开发环境
				return true
			}
			return strings.Contains(origin, "yourcompany.com")
		},
		MaxAge: 12 * time.Hour,
	})
}
