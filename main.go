package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/wureny/webook/webook/Internal/repository"
	"github.com/wureny/webook/webook/Internal/repository/dao"
	"github.com/wureny/webook/webook/Internal/service"
	"github.com/wureny/webook/webook/Internal/web"
	"github.com/wureny/webook/webook/Internal/web/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

func main() {
	db := initDB()
	server := initWebServer()
	user := initUser(db)
	user.RegisterRoutes(server)
	server.Run(":8081")
	e := gin.Default()
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
func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webooktest?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		// 我只会在初始化过程中 panic
		// panic 相当于整个 goroutine 结束
		// 一旦初始化过程出错，应用就不要启动了
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}
func initWebServer() *gin.Engine {
	e := gin.Default()
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
	e.NoMethod(func(g *gin.Context) {
		g.String(405, "no such method for this route")
	})
	e.NoRoute(func(g *gin.Context) {
		g.String(http.StatusNotFound, "no such page!")
	})
	//store := cookie.NewStore([]byte("secret"))
	//	store, err := redis.NewStore(16,
	//		"tcp", "localhost:6379", "",
	//		[]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), []byte("0Pf2r0wZBpXVXlQNdpwCXN4ncnlnZSc3"))

	//	if err != nil {
	//		panic(err)
	//	}
	//	e.Use(sessions.Sessions("mysession", store))
	//	e.Use(middleware.NewLoginMiddlewareBuilder().
	//		IgnorePaths("/users/signup").
	//		IgnorePaths("/users/login").Build())
	e.Use(middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePaths("/users/signup").
		IgnorePaths("/users/loginJWT").Build())
	return e
}
