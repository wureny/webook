package web

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/wureny/webook/webook/Internal/domain"
	"github.com/wureny/webook/webook/Internal/service"
	"net/http"
	"time"
)

type UserHandler struct {
	svc         *service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern = "^(?=.*[A-Za-z])(?=.*\\d)(?=.*[$@$!%*#?&])[A-Za-z\\d$@$!%*#?&]{8,}$"
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
	}
}
func (u *UserHandler) SignUp(ctx *gin.Context) {
	//内部结构
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}
	var req SignUpReq
	// Bind 方法会根据 Content-Type 来解析你的数据到 req 里面
	// 解析错了，就会直接写回一个 400 的错误
	if err := ctx.Bind(&req); err != nil {
		return
	}
	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "system failed")
		fmt.Println("system failed")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "wrong format email")
		return
	}
	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次输入的密码不一致")
		return
	}
	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		// 记录日志
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码必须大于8位，包含数字、特殊字符")
		return
	}
	fmt.Println("web.test")
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrUserDuplicateEmail {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}

	ctx.String(http.StatusOK, "注册成功")
}
func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	// 步骤2
	// 在这里登录成功了
	// 设置 session
	sess := sessions.Default(ctx)
	// 我可以随便设置值了
	// 你要放在 session 里面的值
	sess.Set("userId", user.Id)
	sess.Options(sessions.Options{
		MaxAge:   60,
		Secure:   false,
		HttpOnly: true,
	})
	sess.Save()
	ctx.String(http.StatusOK, "登录成功")
	return
}
func (u *UserHandler) LogOut(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	sess.Options(sessions.Options{
		MaxAge: -1,
	})
	sess.Save()
	ctx.String(http.StatusOK, "log out!")
}
func (u *UserHandler) Edit(ctx *gin.Context) {
	type edit struct {
		Birthday string `json:"birthday"`
		UserName string `json:"userName"`
		Bio      string `json:"bio"`
	}
	session := sessions.Default(ctx)
	userid := session.Get("userId")
	var s edit
	err := ctx.Bind(&s)
	if err != nil {
		ctx.String(http.StatusOK, "miss something")
		return
	}
	_, er := time.Parse("2006-01-02", s.Birthday)
	if er != nil {
		ctx.String(http.StatusOK, "wrong birthday")
		return
	}
	maxlen := 500
	if len(s.Bio) > maxlen {
		ctx.String(http.StatusOK, "Bio is too long!")
		return
	}
	userID, _ := userid.(uint64)
	e := u.svc.Edit(ctx, domain.User{
		Id:       userID,
		Birthday: s.Birthday,
		UserName: s.UserName,
		Bio:      s.Bio,
	})
	if e != nil {
		ctx.String(http.StatusOK, "edit未成功")
		return
	}
	ctx.String(http.StatusOK, "edit成功")
}

func (u *UserHandler) Profile(ctx *gin.Context) {
	type Userinfo struct {
		Email    string
		Birthday string
		Bio      string
		UserName string
	}
	session := sessions.Default(ctx)
	userid := session.Get("userId")
	userId := userid.(uint64)
	user, err := u.svc.GetUser(ctx, userId)
	if err != nil {
		ctx.String(http.StatusOK, "failed to get the info")
		return
	}
	tmp := Userinfo{
		Email:    user.Email,
		Birthday: user.Birthday,
		Bio:      user.Bio,
		UserName: user.UserName,
	}
	ctx.JSON(http.StatusOK, tmp)
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid uint64
}

func (u *UserHandler) LoginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	// 步骤2
	// 在这里用 JWT 设置登录态
	// 生成一个 JWT token
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		},
		Uid: user.Id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	ctx.Header("x-jwt-token", tokenStr)
	fmt.Println(user)
	ctx.String(http.StatusOK, "登录成功")
	return
}

func (u *UserHandler) ProfileJWT(ctx *gin.Context) {
	c, _ := ctx.Get("claims")
	// 你可以断定，必然有 claims
	//if !ok {
	//	// 你可以考虑监控住这里
	//	ctx.String(http.StatusOK, "系统错误")
	//	return
	//}
	// ok 代表是不是 *UserClaims
	claims, ok := c.(*UserClaims)
	if !ok {
		// 你可以考虑监控住这里
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	println(claims.Uid)
	// 这边就是你补充 profile 的其它代码
	user, err := u.svc.GetUser(ctx, claims.Uid)
	if err != nil {
		ctx.String(http.StatusOK, "failed to get the info")
		return
	}
	type Userinfo struct {
		Email    string
		Birthday string
		Bio      string
		UserName string
	}
	tmp := Userinfo{
		Email:    user.Email,
		Birthday: user.Birthday,
		Bio:      user.Bio,
		UserName: user.UserName,
	}
	ctx.JSON(http.StatusOK, tmp)
}

func (u *UserHandler) EditJWT(ctx *gin.Context) {
	c, _ := ctx.Get("claims")
	// 你可以断定，必然有 claims
	//if !ok {
	//	// 你可以考虑监控住这里
	//	ctx.String(http.StatusOK, "系统错误")
	//	return
	//}
	// ok 代表是不是 *UserClaims
	claims, ok := c.(*UserClaims)
	if !ok {
		// 你可以考虑监控住这里
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	println(claims.Uid)
	type edit struct {
		Birthday string `json:"birthday"`
		UserName string `json:"userName"`
		Bio      string `json:"bio"`
	}
	var s edit
	err := ctx.Bind(&s)
	if err != nil {
		ctx.String(http.StatusOK, "miss something")
		return
	}
	_, er := time.Parse("2006-01-02", s.Birthday)
	if er != nil {
		ctx.String(http.StatusOK, "wrong birthday")
		return
	}
	maxlen := 500
	if len(s.Bio) > maxlen {
		ctx.String(http.StatusOK, "Bio is too long!")
		return
	}
	userID := claims.Uid
	e := u.svc.Edit(ctx, domain.User{
		Id:       userID,
		Birthday: s.Birthday,
		UserName: s.UserName,
		Bio:      s.Bio,
	})
	if e != nil {
		ctx.String(http.StatusOK, "edit未成功")
		return
	}
	ctx.String(http.StatusOK, "edit成功")
}
