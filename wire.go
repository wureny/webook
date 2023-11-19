//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/wureny/webook/webook/Internal/repository"
	"github.com/wureny/webook/webook/Internal/repository/cache"
	"github.com/wureny/webook/webook/Internal/repository/dao"
	"github.com/wureny/webook/webook/Internal/service"
	"github.com/wureny/webook/webook/Internal/service/sms"
	"github.com/wureny/webook/webook/Internal/web"
	"github.com/wureny/webook/webook/Internal/web/jwt"
	"github.com/wureny/webook/webook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(ioc.InitDB, ioc.InitRedis, dao.NewUserDAO, cache.NewUserCache, cache.NewCodeCache, repository.NewCodeRepository, repository.NewUserRepository,
		service.NewUserService, service.NewCodeService,
		//	memory.NewService,
		web.NewUserHandler,
		ioc.InitMiddlewares,
		ioc.InitWebServer,
		sms.Newseservice,
		jwt.NewRedisJWTHandler,
	)
	return gin.Default()
}
