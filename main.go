package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/wureny/webook/webook/Internal/web"
	"strings"
	"time"
)

func main() {
	e := gin.Default()
	web.RegisterRoutes(e)
	e.Use(cors.New(cors.Config{
		//	AllowAllOrigins:        false,
		//	AllowOrigins:           nil,
		//	AllowMethods:           nil,
		AllowHeaders:     []string{"content-type", "Authorization"},
		AllowCredentials: true,
		//	ExposeHeaders:          nil,
		MaxAge: 7 * time.Hour,
		//	AllowWildcard:          false,
		//	AllowBrowserExtensions: false,
		//	AllowWebSockets:        false,
		//	AllowFiles:             false,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				//开发环境，因为只有他有localhost
				return true
			}
			return strings.Contains(origin, "yourcompany.com")
		},
	}))
	e.Run(":8080")
}
