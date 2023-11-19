package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/wureny/webook/webook/Internal/repository"
	"github.com/wureny/webook/webook/Internal/repository/cache"
	"github.com/wureny/webook/webook/Internal/repository/dao"
	"github.com/wureny/webook/webook/Internal/service"
	"github.com/wureny/webook/webook/Internal/service/sms/memory"
	"github.com/wureny/webook/webook/Internal/web"
	"github.com/wureny/webook/webook/Internal/web/jwt"
	"github.com/wureny/webook/webook/Internal/web/middleware"
	"github.com/wureny/webook/webook/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

func main() {
	/*db := initDB()
	server := initWebServer()
	rdb := initRedis()
	user := initUser(db, rdb)
	user.RegisterRoute(server)*/
	//		server.Run(":8081")

	//	server := gin.Default()
	server := initWebServer()
	server.GET("/hello", func(ctx *gin.Context) {
		fmt.Println("tttest")
		ctx.String(http.StatusOK, "你好，你来了")
	})
	fmt.Println("halo")
	server.GET("/try", func(ctx *gin.Context) {
		fmt.Println("nice try")
	})
	server.Run(":8083")
}
func initDB() *gorm.DB {
	//	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webooktest?charset=utf8mb4&parseTime=True&loc=Local"))
	//	db, err := gorm.Open(mysql.Open("root:root@tcp(webook-live-mysql:11309)/webook"))
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	//	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:30002)/webook"))
	fmt.Println(config.Config.DB.DSN)
	fmt.Println("fuck")
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
func initUser(db *gorm.DB, rdb redis.Cmdable) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	uc := cache.NewUserCache(rdb)
	repo := repository.NewUserRepository(ud, uc)
	svc := service.NewUserService(repo)
	codeCache := cache.NewCodeCache(rdb)
	codeRepo := repository.NewCodeRepository(codeCache)
	smsSvc := memory.NewService()
	codeSvc := service.NewCodeService(codeRepo, &smsSvc)
	jwtHandler := jwt.NewRedisJWTHandler(rdb)
	u := web.NewUserHandler(svc, codeSvc, jwtHandler)
	return u
}
func initWebServer() *gin.Engine {
	e := gin.Default()
	e.Use(func(ctx *gin.Context) {
		println("这是第一个 middleware")
	})

	e.Use(func(ctx *gin.Context) {
		fmt.Println("first middleware")
	})
	//redisClient := redis.NewClient(&redis.Options{
	//		Addr: "localhost:6379",
	//	Addr: "webook-live-redis:11479",
	//	Addr: config.Config.Redis.Addr,
	//	})
	fmt.Println(config.Config.Redis.A)
	//e.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())
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
	jwtHandler := jwt.NewRedisJWTHandler(initRedis())
	e.Use(middleware.NewLoginJWTMiddlewareBuilder(jwtHandler).
		IgnorePaths("/users/signup").
		IgnorePaths("/hello").
		IgnorePaths("/login_sms/code/send").
		IgnorePaths("/login_sms").
		IgnorePaths("/users/loginJWT").Build())
	return e
}

func mid(ctx *gin.Context) {
	fmt.Println("first middleware")
}

func initRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	return redisClient
}
