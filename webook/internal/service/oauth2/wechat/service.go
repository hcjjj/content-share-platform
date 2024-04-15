package wechat

import (
	"basic-go/webook/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var redirectURI = url.PathEscape("https://meoying.com/oauth2/wechat/callback")

type Service interface {
	AuthURL(ctx context.Context, state string) (string, error)
	VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error)
}

type service struct {
	appId     string
	appSecret string
	client    *http.Client
	//cmd       redis.Cmdable
}

// 不偷懒的写法
func NewServiceV1(appId string, appSecret string, client *http.Client) Service {
	return &service{
		appId:     appId,
		appSecret: appSecret,
		client:    client,
	}
}

func NewService(appId string, appSecret string) Service {
	return &service{
		appId:     appId,
		appSecret: appSecret,
		// 依赖注入，但是没完全注入
		client: http.DefaultClient,
	}
}

func (s *service) VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error) {
	const targetPattern = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	target := fmt.Sprintf(targetPattern, s.appId, s.appSecret, code)
	//resp, err := http.Get(target)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	//req, err := http.NewRequest(http.MethodGet, target, nil)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	// 会产生复制，性能极差，比如说你的 URL 很长
	//req = req.WithContext(ctx)

	resp, err := s.client.Do(req)
	if err != nil {
		return domain.WechatInfo{}, err
	}

	// 只读一遍
	decoder := json.NewDecoder(resp.Body)
	var res Result
	err = decoder.Decode(&res)

	// 整个响应都读出来，不推荐，因为 Unmarshal 再读一遍，合计两遍
	//body, err := io.ReadAll(resp.Body)
	//err = json.Unmarshal(body, &res)

	if err != nil {
		return domain.WechatInfo{}, err
	}

	if res.ErrCode != 0 {
		return domain.WechatInfo{},
			fmt.Errorf("微信返回错误响应，错误码：%d，错误信息：%s", res.ErrCode, res.ErrMsg)
	}

	// 攻击者的 state
	//str := s.cmd.Get(ctx, "my-state"+state).String()
	//if str != state {
	//	// 不相等
	//}

	return domain.WechatInfo{
		OpenID:  res.OpenID,
		UnionID: res.UnionID,
	}, nil
}

func (s *service) AuthURL(ctx context.Context, state string) (string, error) {
	const urlPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect"
	// 如果在这里存 state，假如说我存 redis
	//s.cmd.Set(ctx, "my-state"+state, state, time.Minute)
	return fmt.Sprintf(urlPattern, s.appId, redirectURI, state), nil
}

type Result struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`

	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`

	OpenID  string `json:"openid"`
	Scope   string `json:"scope"`
	UnionID string `json:"unionid"`
}
