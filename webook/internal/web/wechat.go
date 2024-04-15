package web

import (
	"basic-go/webook/internal/service"
	"basic-go/webook/internal/service/oauth2/wechat"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
)

type OAuth2WechatHandler struct {
	svc     wechat.Service
	userSvc service.UserService
	jwtHandler
	stateKey []byte
	cfg      WechatHandlerConfig
}

type WechatHandlerConfig struct {
	Secure bool
	//StateKey
}

func NewOAuth2WechatHandler(svc wechat.Service,
	userSvc service.UserService,
	cfg WechatHandlerConfig) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:      svc,
		userSvc:  userSvc,
		stateKey: []byte("jaks3jgvkjoiGezwd4QbE9ujPZp0fL8p"),
		cfg:      cfg,
	}
}

func (h *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", h.AuthURL)
	g.Any("/callback", h.Callback)
}

func (h *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {
	state := uuid.New()
	url, err := h.svc.AuthURL(ctx, state)
	// 要把我的 state 存好
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "构造扫码登录URL失败",
		})
		return
	}
	if err = h.setStateCookie(ctx, state); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: url,
	})
}

func (h *OAuth2WechatHandler) setStateCookie(ctx *gin.Context, state string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, StateClaims{
		State: state,
		RegisteredClaims: jwt.RegisteredClaims{
			// 过期时间，你预期中一个用户完成登录的时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		},
	})
	tokenStr, err := token.SignedString(h.stateKey)
	if err != nil {
		return err
	}
	ctx.SetCookie("jwt-state", tokenStr,
		600, "/oauth2/wechat/callback",
		"", h.cfg.Secure, true)
	return nil
}

func (h *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	code := ctx.Query("code")
	err := h.verifyState(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "登录失败",
		})
		return
	}
	info, err := h.svc.VerifyCode(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	// 这里怎么办？
	// 从 userService 里面拿 uid
	u, err := h.userSvc.FindOrCreateByWechat(ctx, info)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	err = h.setJWTToken(ctx, u.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	err = h.setRefreshToken(ctx, u.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
	// 验证微信的 code
}

func (h *OAuth2WechatHandler) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")
	// 校验一下我的 state
	ck, err := ctx.Cookie("jwt-state")
	if err != nil {
		return fmt.Errorf("拿不到 state 的 cookie, %w", err)
	}

	var sc StateClaims
	token, err := jwt.ParseWithClaims(ck, &sc, func(token *jwt.Token) (interface{}, error) {
		return h.stateKey, nil
	})
	if err != nil || !token.Valid {
		return fmt.Errorf("token 已经过期了, %w", err)
	}

	if sc.State != state {
		return errors.New("state 不相等")
	}
	return nil
}

type StateClaims struct {
	State string
	jwt.RegisteredClaims
}

//type OAuth2Handler struct {
//	wechatService
//	dingdingService
//	feishuService
//
//	svcs map[string]OAuth2Service
//}

//func (h *OAuth2Handler) RegisterRoutes(server *gin.Engine) {
//	// 统一处理所有的 OAuth2 的
//	g := server.Group("/oauth2")
//	g.GET("/:platform/authurl", h.AuthURL)
//	g.Any("/:platform/callback", h.Callback)
//}

//func (h *OAuth2Handler) AuthURL(ctx *gin.Context) {
//	platform := ctx.Param("platform")
//	switch platform {
//	case "wechat":
//		h.wechatService.AuthURL
//	}
//
//	svc := h.svcs[platform]
//	svc.
//}

//func (h *OAuth2Handler) Callback(ctx *gin.Context) {
//
//}
