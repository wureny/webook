package web

import "github.com/gin-gonic/gin"

func (user *UserHandler) RegisterRoutes(e *gin.Engine) *gin.Engine {
	//	user := NewUserHandler()

	e.POST("/users/signup", user.SignUp)
	e.POST("users/login", user.Login)
	e.POST("users/edit", user.Edit)
	e.GET("users/profile", user.Profile)
	return e
}
