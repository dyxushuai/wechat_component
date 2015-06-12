package wechat

import (
	"git.ishopex.cn/xushuai/wechat/lib"

	"fmt"
)

var (
	weixinComponentHost    = "https://api.weixin.qq.com/cgi-bin/component"
	apiComponentToken      = weixinComponentHost + "/api_component_token"
	apiCreatePreAuthCode   = weixinComponentHost + "/api_create_preauthcode?component_access_token=%s"
	apiQueryAuth           = weixinComponentHost + "/api_query_auth?component_access_token=%s"
	apiAuthorizerToken     = weixinComponentHost + "/api_authorizer_token?component_access_token=%s"
	apiGetAuthorizerInfo   = weixinComponentHost + "/api_get_authorizer_info?component_access_token=%s"
	apiGetAuthorizerOption = weixinComponentHost + "/api_get_authorizer_option?component_access_token=%s"
	apiSetAuthorizerOption = weixinComponentHost + "/api_set_authorizer_option?component_access_token=%s"

	oauthUrl = "https://mp.weixin.qq.com/cgi-bin/componentloginpage?component_appid=%s&pre_auth_code=%s&redirect_uri=%s"
)

// 微信第三方平台interface
type WechatComponent interface {
	GetRegularApi() RegularApi
	GetNormalApi() NormalApi
	OAuthUrl(redirectUrl, preAuthCode string) string
	GetCipher() (lib.IOCipher, error)
}

func New(appId, appSecret, cryptoKey, token string) WechatComponent {
	return &WechatThird{
		appId:     appId,
		appSecret: appSecret,
		cryptoKey: cryptoKey,
		token:     token,
	}
}

// 微信第三方公众号平台 实现 ServeHTTP interface
type WechatThird struct {
	appId     string // 第三方应用id
	appSecret string // 第三方应用secret
	cryptoKey string // 公众号消息加解密Key
	token     string // 公众号消息校验Token
}

func (wt *WechatThird) GetRegularApi() RegularApi {
	return &regularApi{
		wt: wt,
	}
}
func (wt *WechatThird) GetNormalApi() NormalApi {
	return &normalApi{
		wt: wt,
	}
}
func (wt *WechatThird) GetCipher() (lib.IOCipher, error) {
	return lib.NewCipher(wt.token, wt.cryptoKey, wt.appId)
}

func (wt *WechatThird) OAuthUrl(redirectUrl, preAuthCode string) string {
	u, _ := UrlEncoded(redirectUrl)

	return fmt.Sprintf(oauthUrl, wt.appId, preAuthCode, u)
}
