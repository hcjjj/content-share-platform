package web

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type jwtHandler struct {
	// assess_token key
	atKey []byte
	// refresh_token key
	rtKey []byte
}

func NewJwtHandler() jwtHandler {
	return jwtHandler{
		atKey: []byte("jaks3jgvkjoiGezwd4QbE9ujPZp0fL8p"),
		rtKey: []byte("jaks3jgvkjoiGezwd4QbE9ujPZp0fL8p"),
	}
}

func (h jwtHandler) setJWTToken(ctx *gin.Context, uid int64) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		Uid:       uid,
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	// 这里的 “ ” 得 和 中间件那边的一致
	tokenStr, err := token.SignedString(h.atKey)
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

func (h jwtHandler) setRefreshToken(ctx *gin.Context, uid int64) error {
	claims := RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		Uid: uid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(h.rtKey)
	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", tokenStr)
	return nil
}

func ExtractToken(ctx *gin.Context) string {
	// 我现在用 JWT 来校验
	tokenHeader := ctx.GetHeader("Authorization")
	//segs := strings.SplitN(tokenHeader, " ", 2)
	segs := strings.Split(tokenHeader, " ")
	if len(segs) != 2 {
		return ""
	}
	return segs[1]
}

type RefreshClaims struct {
	Uid  int64
	Ssid string
	jwt.RegisteredClaims
}

type UserClaims struct {
	jwt.RegisteredClaims
	// 声明你自己的要放进去 token 里面的数据
	// 敏感数据不要放这
	Uid int64
	// 自己随便加
	UserAgent string
}
