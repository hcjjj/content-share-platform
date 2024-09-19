package web

import (
	"net/http"
	"time"

	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/errs"
	"basic-go/webook/internal/service"
	ijwt "basic-go/webook/internal/web/jwt"
	"basic-go/webook/pkg/ginx"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	bizLogin             = "login"
)

type UserHandler struct {
	ijwt.Handler
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
	svc            service.UserService
	codeSvc        service.CodeService
}

func NewUserHandler(svc service.UserService,
	hdl ijwt.Handler,
	codeSvc service.CodeService) *UserHandler {
	return &UserHandler{
		// 预编译正则表达式来提高校验速度
		emailRexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:            svc,
		codeSvc:        codeSvc,
		Handler:        hdl,
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	// REST 风格
	//server.POST("/user", h.SignUp)
	//server.PUT("/user", h.SignUp)
	//server.GET("/users/:username", h.Profile)
	ug := server.Group("/users")
	// POST /users/signup
	ug.POST("/signup", ginx.WrapBody(h.SignUp))
	// POST /users/login
	//ug.POST("/login", h.Login)
	ug.POST("/login", ginx.WrapBody(h.LoginJWT))
	ug.POST("/logout", h.LogoutJWT)
	// POST /users/edit
	ug.POST("/edit", ginx.WrapBodyAndClaims(h.Edit))
	// GET /users/profile
	ug.GET("/profile", ginx.WrapClaims(h.Profile))
	ug.GET("/refresh_token", h.RefreshToken)

	// 手机验证码登录相关功能
	ug.POST("/login_sms/code/send", ginx.WrapBody(h.SendSMSLoginCode))
	ug.POST("/login_sms", ginx.WrapBody(h.LoginSMS))
}

func (h *UserHandler) LoginSMS(ctx *gin.Context,
	req LoginSMSReq) (ginx.Result, error) {
	ok, err := h.codeSvc.Verify(ctx, bizLogin, req.Phone, req.Code)
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统异常",
			//Msg: err.Error(),
		}, err
	}
	if !ok {
		return ginx.Result{
			Code: 4,
			Msg:  "验证码不对，请重新输入",
		}, nil
	}
	u, err := h.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	err = h.SetLoginToken(ctx, u.Id)
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	return ginx.Result{
		Msg: "登录成功",
	}, nil
}

func (h *UserHandler) SendSMSLoginCode(ctx *gin.Context,
	req SendSMSCodeReq) (ginx.Result, error) {
	// 这边可以校验 Req
	if req.Phone == "" {
		return ginx.Result{
			Code: 4,
			Msg:  "请输入手机号码",
		}, nil
	}
	err := h.codeSvc.Send(ctx, bizLogin, req.Phone)
	switch err {
	case nil:
		return ginx.Result{
			Msg: "发送成功",
		}, nil
	case service.ErrCodeSendTooMany:
		// 事实上，防不住有人不知道怎么触发了
		// 少数这种错误，是可以接受的
		// 但是频繁出现，就代表有人在搞的系统
		return ginx.Result{
			Code: 4,
			Msg:  "短信发送太频繁，请稍后再试",
		}, nil
	default:
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
}

func (h *UserHandler) SignUp(ctx *gin.Context, req SignUpReq) (ginx.Result, error) {
	isEmail, err := h.emailRexExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return ginx.Result{
			Code: errs.UserInternalServerError,
			Msg:  "系统错误",
		}, err
	}
	if !isEmail {
		return ginx.Result{
			Code: errs.UserInvalidInput,
			Msg:  "非法邮箱格式",
		}, nil
	}

	if req.Password != req.ConfirmPassword {
		return ginx.Result{
			Code: errs.UserInvalidInput,
			Msg:  "两次输入的密码不相等",
		}, nil
	}

	isPassword, err := h.passwordRexExp.MatchString(req.Password)
	if err != nil {
		return ginx.Result{
			Code: errs.UserInternalServerError,
			Msg:  "系统错误",
		}, err
	}
	if !isPassword {
		return ginx.Result{
			Code: errs.UserInvalidInput,
			Msg:  "密码必须包含字母、数字、特殊字符",
		}, nil
	}

	err = h.svc.Signup(ctx.Request.Context(), domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	switch err {
	case nil:
		return ginx.Result{
			Msg: "OK",
		}, nil
	case service.ErrDuplicateEmail:
		return ginx.Result{
			Code: errs.UserDuplicateEmail,
			Msg:  "邮箱冲突",
		}, nil
	default:
		return ginx.Result{
			Code: errs.UserInternalServerError,
			Msg:  "系统错误",
		}, err
	}
}

func (h *UserHandler) LoginJWT(ctx *gin.Context, req LoginJWTReq) (ginx.Result, error) {
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		err = h.SetLoginToken(ctx, u.Id)
		if err != nil {
			return ginx.Result{
				Code: 5,
				Msg:  "系统错误",
			}, err
		}
		return ginx.Result{
			Msg: "OK",
		}, nil
	case service.ErrInvalidUserOrPassword:
		return ginx.Result{Msg: "用户名或者密码错误"}, nil
	default:
		return ginx.Result{Msg: "系统错误"}, err
	}
}

//func (h *UserHandler) Logout(ctx *gin.Context) {
//	sess := sessions.Default(ctx)
//	sess.Options(sessions.Options{
//		MaxAge: -1,
//	})
//	sess.Save()
//}

func (h *UserHandler) Login(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	// Bind 方法会根据 HTTP 请求的 Content-Type 来决定怎么处理
	if err := ctx.Bind(&req); err != nil {
		return
	}
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		//Gin 的session 插件，会从 Cookie 里面找到 sess_id，再根据 sess_id 找到对应的 Session
		sess := sessions.Default(ctx)
		sess.Set("userId", u.Id)
		sess.Options(sessions.Options{
			// 十分钟
			MaxAge: 30,
		})
		err = sess.Save()
		if err != nil {
			ctx.String(http.StatusOK, "系统错误")
			return
		}
		ctx.String(http.StatusOK, "登录成功")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "用户名或者密码不对")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

func (h *UserHandler) Edit(ctx *gin.Context, req UserEditReq,
	uc ijwt.UserClaims) (ginx.Result, error) {
	// 嵌入一段刷新过期时间的代码
	//sess := sessions.Default(ctx)
	//sess.Get("uid")
	// 用户输入不对
	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		return ginx.Result{
			Code: 4,
			Msg:  "生日格式不对",
		}, err
	}
	err = h.svc.UpdateNonSensitiveInfo(ctx, domain.User{
		Id:       uc.Uid,
		Nickname: req.Nickname,
		Birthday: birthday,
		AboutMe:  req.AboutMe,
	})
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	return ginx.Result{
		Msg: "OK",
	}, nil
}

func (h *UserHandler) Profile(ctx *gin.Context,
	uc ijwt.UserClaims) (ginx.Result, error) {
	//us := ctx.MustGet("user").(UserClaims)
	//ctx.String(http.StatusOK, "这是 profile")
	u, err := h.svc.FindById(ctx, uc.Uid)
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	type User struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		AboutMe  string `json:"aboutMe"`
		Birthday string `json:"birthday"`
	}
	return ginx.Result{
		Data: User{
			Nickname: u.Nickname,
			Email:    u.Email,
			AboutMe:  u.AboutMe,
			Birthday: u.Birthday.Format(time.DateOnly),
		},
	}, nil
}

func (h *UserHandler) RefreshToken(ctx *gin.Context) {
	// 约定，前端在 Authorization 里面带上这个 refresh_token
	tokenStr := h.ExtractToken(ctx)
	var rc ijwt.RefreshClaims
	token, err := jwt.ParseWithClaims(tokenStr, &rc, func(token *jwt.Token) (interface{}, error) {
		return ijwt.RCJWTKey, nil
	})
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if token == nil || !token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = h.CheckSession(ctx, rc.Ssid)
	if err != nil {
		// token 无效或者 redis 有问题
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = h.SetJWTToken(ctx, rc.Uid, rc.Ssid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Msg: "OK",
	})
}

func (h *UserHandler) LogoutJWT(ctx *gin.Context) {
	err := h.ClearToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{Code: 5, Msg: "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{Msg: "退出登录成功"})
}
