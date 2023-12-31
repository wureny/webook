package web

import "github.com/gin-gonic/gin"

func (user *UserHandler) RegisterRoute(e *gin.Engine) *gin.Engine {
	//	user := NewUserHandler()

	e.POST("/users/signup", user.SignUp)
	e.POST("/users/login", user.Login)
	e.POST("/users/edit", user.Edit)
	e.GET("/users/profile", user.Profile)
	e.POST("/users/loginJWT", user.LoginJWT)
	e.POST("/users/editJWT", user.EditJWT)
	e.GET("/users/profileJWT", user.ProfileJWT)
	e.POST("/users/login_sms/code/send", user.SendLoginSMSCode)
	e.POST("/users/login_sms", user.LoginSMS)

	return e
}
