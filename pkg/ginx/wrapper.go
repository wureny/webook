package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

func WrapToken[C jwt.Claims](fn func(ctx *gin.Context, uc C) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c, ok := ctx.Get("users")
		if !ok {
			ctx.AbortWithStatus(http.StatusNonAuthoritativeInfo)
			return
		}
		val, ok := c.(C)
		if !ok {
			ctx.AbortWithStatus(http.StatusNonAuthoritativeInfo)
			return
		}
		res, err := fn(ctx, val)
		if err != nil {
			//TODO:日志处理错误
		}
		ctx.JSON(http.StatusOK, res)
		//继续执行一些逻辑
	}
}

func WrapBodyAndToken[Req any, C jwt.Claims](fn func(ctx *gin.Context, req Req, uc C) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.Bind(&req); err != nil {
			return
		}

		val, ok := ctx.Get("users")
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c, ok := val.(C)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		res, err := fn(ctx, req, c)
		if err != nil {
			//TODO: 日志处理错误
		}
		ctx.JSON(http.StatusOK, res)
	}
}

func WrapBodyV1[T any](fn func(ctx *gin.Context, req T) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req T
		if err := ctx.Bind(&req); err != nil {
			return
		}
		res, err := fn(ctx, req)
		if err != nil {
			//TODO：日志处理错误
		}
		ctx.JSON(http.StatusOK, res)
	}
}

func WrapBody[T any](fn func(ctx *gin.Context, req T) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req T
		if err := ctx.Bind(&req); err != nil {
			return
		}
		res, err := fn(ctx, req)
		if err != nil {
			//TODO: 日志处理错误
		}
		ctx.JSON(http.StatusOK, res)
	}
}
