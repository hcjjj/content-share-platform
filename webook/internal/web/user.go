// Package web -----------------------------
// @file      : user.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-01-29 20:45
// -------------------------------------------
package web

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/service"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-contrib/sessions"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
)

// UserHandler 定义用户有关的路由
type UserHandler struct {
	svc         service.UserService
	codeSvc     service.CodeService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	// 组合
	jwtHandler
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService) *UserHandler {
	const (
		emailRegexPattern    = "\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*"
		passwordRegexPattern = "^(?=.*\\d)(?=.*[a-zA-Z])(?=.*[^\\da-zA-Z\\s]).{1,16}$"
	)
	// 预编译所需的正则表达式
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
		codeSvc:     codeSvc,
		jwtHandler:  NewJwtHandler(),
	}
}

//func (u *UserHandler) RegisterRoutesV1(ug *gin.RouterGroup) {
//	ug.POST("/signup", u.SignUp)
//	ug.POST("/login", u.Login)
//	ug.POST("/edit", u.Edit)
//	ug.GET("/profile", u.Profile)
//}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	//server.POST("/users/signup", u.SignUp)
	//server.POST("/users/login", u.Login)
	//server.POST("/users/edit", u.Edit)
	//server.GET("/users/profile", u.Profile)
	ug := server.Group("/users")
	ug.POST("/signup", u.SignUp)
	//ug.POST("/login", u.Login)
	ug.POST("/login", u.LoginJWT)
	//ug.POST("/logout", u.Logout)
	ug.POST("/edit", u.Edit)
	//ug.GET("/profile", u.Profile)
	ug.GET("/profile", u.ProfileJWT)
	// PUT "/login/sms/code" 发送验证码
	// POST "/login/sms/code" 校验验证码
	// POST "/sms/login/code"
	ug.POST("/login_sms/code/send", u.SendLoginSMSCode)
	ug.POST("/login_sms", u.LoginSMS)
	ug.POST("/refresh_token", u.RefreshToken)
}

const biz = "login"

func (u *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 可以加上各种校验

	ok, err := u.codeSvc.Verify(ctx, biz, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "验证码有误",
		})
		return
	}
	// 这个手机号是否会是新用户？
	user, err := u.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
	// 登录成功
	if err = u.setJWTToken(ctx, user.Id); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
	if err = u.setRefreshToken(ctx, user.Id); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 4,
		Msg:  "登录成功",
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
	// 需要添加使用正则表达式校验 手机号格式
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "输入有误",
		})
	}
	err := u.codeSvc.Send(ctx, biz, req.Phone)
	// 这边可以细化错误
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

func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}
	var req SignUpReq
	// 接收请求
	// Bind 方法会根据 Content-type 解析数据到 req
	// 解析出错回直接写回一个 400 错误
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 数据校验
	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "邮箱格式不对")
		return
	}

	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次输入的密码不一致")
		return
	}

	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		// 记录日志
		// ...
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "至少包含字母、数字、特殊字符，1-16位")
		return
	}

	// 调用一下 svc 的方法
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrUserDuplicate {
		ctx.String(http.StatusOK, "邮箱或者手机号码冲突")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	ctx.String(http.StatusOK, "注册成功")
	//fmt.Printf("%v\n", req)
}

// RefreshToken 可以同时刷新长短 token，用 redis 来记录是否有效，即 refresh_token 是一次性的
// 参考登录校验部分，比较 User-Agent 来增强安全性
func (u *UserHandler) RefreshToken(ctx *gin.Context) {
	// 只有这个接口，拿出来的才是 refresh_token，其它地方都是 access token
	// 需要和前端配合， 前端需要拿 refresh_token 来刷新
	refreshToken := ExtractToken(ctx)
	var rc RefreshClaims
	token, err := jwt.ParseWithClaims(refreshToken, &rc, func(token *jwt.Token) (interface{}, error) {
		return u.rtKey, nil
	})
	if err != nil || !token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// 搞个新的 access_token
	err = u.setJWTToken(ctx, rc.Uid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "刷新成功",
	})
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
		ctx.String(http.StatusOK, "邮箱或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	// 到这里就登录成功了
	// 用 JWT 设置登录态
	// 生成长短 token
	if err := u.setJWTToken(ctx, user.Id); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if err := u.setRefreshToken(ctx, user.Id); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	//fmt.Println(user)
	//fmt.Println(tokenStr)
	ctx.String(http.StatusOK, "登录成功")
	return
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
		ctx.String(http.StatusOK, "邮箱或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	// 到这里就登录成功了
	// 设置 session
	sess := sessions.Default(ctx)
	// 设置放在 session 里面的值
	sess.Set("userId", user.Id)
	sess.Options(sessions.Options{
		// 10 min 过期
		MaxAge: 10 * 60,
	})
	sess.Save()
	ctx.String(http.StatusOK, "登录成功")
	return
}

func (u *UserHandler) Edit(ctx *gin.Context) {
	type Req struct {
		// 注意，其它字段，尤其是密码、邮箱和手机，
		// 修改都要通过别的手段
		// 邮箱和手机都要验证
		// 密码更加不用多说了
		Nickname string `json:"nickname"`
		// 2023-01-01
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}

	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 你可以尝试在这里校验。
	// 比如说你可以要求 Nickname 必须不为空
	// 校验规则取决于产品经理
	if req.Nickname == "" {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "昵称不能为空"})
		return
	}

	if len(req.AboutMe) > 1024 {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "关于我过长"})
		return
	}
	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		// 也就是说，我们其实并没有直接校验具体的格式
		// 而是如果你能转化过来，那就说明没问题
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "日期格式不对"})
		return
	}

	// 这边要处理数据库和缓存了
	// .....
	fmt.Printf("[%s %s %s]\n", req.Nickname, birthday, req.AboutMe)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, Result{Msg: "OK"})

}

func (u *UserHandler) Logout(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	sess.Options(sessions.Options{
		MaxAge: -1,
	})
	sess.Save()
	ctx.String(http.StatusOK, "退出登录成功")

}
func (u *UserHandler) Profile(ctx *gin.Context) {
	type Profile struct {
		Email string
	}
	sess := sessions.Default(ctx)
	id := sess.Get("userId").(int64)
	user, err := u.svc.Profile(ctx, id)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, Profile{
		Email: user.Email,
	})
	//ctx.String(http.StatusOK, "这是你的 profile")
}

func (c *UserHandler) ProfileJWT(ctx *gin.Context) {
	type Profile struct {
		Email    string
		Phone    string
		Nickname string
		Birthday string
		AboutMe  string
	}
	uc, _ := ctx.Get("claims")
	claims, ok := uc.(*UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	u, err := c.svc.Profile(ctx, claims.Uid)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, Profile{
		Email:    u.Email,
		Phone:    u.Phone,
		Nickname: u.Nickname,
		Birthday: u.Birthday.Format(time.DateOnly),
		AboutMe:  u.AboutMe,
	})

	//c, _ := ctx.Get("claims")
	//if !ok {
	//	ctx.String(http.StatusOK, "系统错误")
	//	return
	//}
	// 类型断言
	//claims, ok := c.(*UserClaims)
	//if !ok {
	//	ctx.String(http.StatusOK, "系统错误")
	//	return
	//}
	//fmt.Println(claims.Uid)
	//ue, _ := u.svc.Profile(ctx, claims.Uid)
	//ctx.String(http.StatusOK, fmt.Sprintf("用户Id：%d\n邮箱：%s\n手机号：%s\n创建时间：%s", ue.Id, ue.Email, ue.Phone, ue.Ctime))
}
