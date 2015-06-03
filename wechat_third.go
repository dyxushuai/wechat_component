package wechat

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/rakyll/ticktock"
	"github.com/rakyll/ticktock/t"
)

var (
	weixinComponentHost    = "https://api.weixin.qq.com/cgi-bin/component"
	apiComponentToken      = "/api_component_token"
	apiCreatePreAuthCode   = "/api_create_preauthcode"
	apiQueryAuth           = "/api_query_auth"
	apiAuthorizerToken     = "/api_authorizer_token"
	apiGetAuthorizerInfo   = "/api_get_authorizer_info"
	apiGetAuthorizerOption = "/api_get_authorizer_option"
	apiSetAuthorizerOption = "/api_set_authorizer_option"
)

// 微信第三方公众号平台 实现 ServeHTTP interface
type WechatThird struct {
	Config
	verifyTicket string
	access       *accessToken
	authHandle   AuthHandler
	cbHandle     CBHandler
	authPath     string
	cbPath       string
}

func (wt *WechatThird) Init() error {
	if err := wt.access.Run(); err != nil {
		return err
	}
	go wt.getTicket()
	go func() {
		// 没30分钟获取一次accessToken
		ticktock.Schedule("get-access-token", wt.access, &t.When{Each: "30m"})
		ticktock.Start()
	}()
	return nil
}

type accessToken struct {
	w         *WechatThird
	CAT       string        `json:"component_access_token"`
	ExpiresIn time.Duration `json:"expires_in"`
}

func (at *accessToken) Run() error {
	at.CAT, at.ExpiresIn = at.w.getAccessToken()
	return nil
}

type Config struct {
	AppId       string // 第三方应用id
	AppSecret   string // 第三方应用secret
	CryptoKey   string // 公众号消息加解密Key
	Token       string // 公众号消息校验Token
	AuthUrl     string // 授权事件接收URL
	PublicCBUrl string // 公众号消息与事件接收URL
}

func New(c Config) (*WechatThird, error) {
	ah, err := NewAuthHandle(c.Token, c.CryptoKey, c.AppId)
	if err != nil {
		return nil, err
	}

	authUrl, err := url.Parse(c.AuthUrl)
	if err != nil {
		return nil, err
	}

	publicUrl, err := url.Parse(c.PublicCBUrl)
	if err != nil {
		return nil, err
	}

	wt := &WechatThird{
		Config:       c,
		verifyTicket: "",
		access:       &accessToken{},
		authHandle:   ah,
		authPath:     authUrl.Path,
		cbPath:       publicUrl.Path,
	}
	err = wt.Init()
	if err != nil {
		return nil, err
	}
	return wt, nil
}

//定时获取ticket
func (wt *WechatThird) getTicket() {
	for {
		wt.verifyTicket = wt.authHandle.GetTicket()
	}
}

func (wt *WechatThird) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == wt.authPath {
		wt.authHandle.ServeHTTP(w, r)
		return
	}
	wt.cbHandle.ServeHTTP(w, r)
}

// 获取第三方平台令牌
func (wt *WechatThird) getAccessToken() (string, time.Duration) {
	postForm := url.Values{}
	postForm.Set("component_appid", wt.Config.AppId)
	postForm.Set("component_appsecret", wt.Config.AppSecret)
	postForm.Set("component_verify_ticket", wt.verifyTicket)
	res, err := http.PostForm(weixinComponentHost+apiComponentToken, postForm)

	result := &struct {
		CAT       string `json:"component_access_token"`
		ExpiresIn int64  `json:"expires_in"`
	}{}

	ae, err := unmarshalResponseToJson(res, result)
	if err != nil {
		log.Println("Parse access token failed: ", err)

	} else if ae != nil {
		if ae.isError() {
			log.Println("Parse access token failed: ", ae.ErrMsg)
		}
	} else {
		return result.CAT, time.Duration(result.ExpiresIn * 1000 * 1000 * 1000)
	}
	return "", 0
}

// 获取预授权码，用于公众号oauth
func getPreAuthCode() string {
	return ""
}
