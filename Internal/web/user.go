package web

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/wureny/webook/webook/Internal/domain"
	"github.com/wureny/webook/webook/Internal/service"
	ijwt "github.com/wureny/webook/webook/Internal/web/jwt"
	"net/http"
	"time"
)

const biz = "login"

var _ handler = (*UserHandler)(nil)

type UserHandler struct {
	svc         service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	codeSvc     *service.CodeService
	ijwt.Handler
	cmd redis.Cmdable
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {

}

func NewUserHandler(svc service.UserService, codeSvc *service.CodeService, jwtHdl ijwt.Handler) *UserHandler {
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
		codeSvc:     codeSvc,
		Handler:     jwtHdl,
	}
}

func (u *UserHandler) LogoutJWT(ctx *gin.Context) {
	err := u.ClearToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "退出登录失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "退出登录OK",
	})
}

// RefreshToken 可以同时刷新长短 token，用 redis 来记录是否有效，即 refresh_token 是一次性的
// 参考登录校验部分，比较 User-Agent 来增强安全性
func (u *UserHandler) RefreshToken(ctx *gin.Context) {
	// 只有这个接口，拿出来的才是 refresh_token，其它地方都是 access token
	refreshToken := u.ExtractToken(ctx)
	var rc ijwt.RefreshClaims
	token, err := jwt.ParseWithClaims(refreshToken, &rc, func(token *jwt.Token) (interface{}, error) {
		return ijwt.RtKey, nil
	})
	if err != nil || !token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	err = u.CheckSession(ctx, rc.Ssid)
	if err != nil {
		// 要么 redis 有问题，要么已经退出登录
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// 生成新的 token
	err = u.SetJWTToken(ctx, rc.Uid, rc.Ssid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "刷新成功",
	})
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
	userID, _ := userid.(int64)
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
	userId := userid.(int64)
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
	Uid       int64
	UserAgent string
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
	if err := u.setJWTToken(ctx, user.Id); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	fmt.Println(user)
	ctx.String(http.StatusOK, "登录成功")
	return
}

func (u *UserHandler) setJWTToken(ctx *gin.Context, uid int64) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		},
		Uid:       uid,
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

func (u *UserHandler) ProfileJWT(ctx *gin.Context) {
	fmt.Println("this is profilejwt")
	c, _ := ctx.Get("claims")
	// 你可以断定，必然有 claims
	//if !ok {
	//	// 你可以考虑监控住这里
	//	ctx.String(http.StatusOK, "系统错误")
	//	return
	//}
	// ok 代表是不是 *UserClaims
	fmt.Println("getting profilejwt")
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

func (u *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	biz := "login"
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// 这边，可以加上各种校验
	ok, err := u.codeSvc.Verify(ctx, biz, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码有误",
		})
		return
	}

	// 我这个手机号，会不会是一个新用户呢？
	// 这样子
	user, err := u.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	// 这边要怎么办呢？
	// 从哪来？
	if err = u.setJWTToken(ctx, user.Id); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Msg: "验证码校验通过",
	})
}

func (u *UserHandler) SendLoginSMSCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 是不是一个合法的手机号码
	// 考虑正则表达式
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "输入有误",
		})
		return
	}
	biz := "login"
	err := u.codeSvc.Send(ctx, biz, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case service.ErrCodeSendTooMany:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送太频繁，请稍后再试",
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
}
