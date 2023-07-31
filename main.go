package main

import (
	"github.com/gin-gonic/gin"
	"github.com/wureny/webook/webook/Internal/web"
)

func main() {
	e := gin.Default()
	web.RegisterRoutes(e)
	e.Run(":8080")
}
