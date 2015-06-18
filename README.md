#### 微信第三方公众号平台

##### Usage

```go
import (
	"net/http"
	"git.ishopex.cn/xushuai/wechat"
)

var WechatComponentSDK wechat.WechatComponent

func init() {
	//appId 第三方公众平台 appid
	//appSecret 第三方公众平台 appSecret
	//cryptoKey 公众号消息加解密Key
	//token 公众号消息校验Token
	WechatComponentSDK = wechat.New(appId, appSecret, cryptoKey, token)
}

// 使用微信加解密算法封装网络读写接口
var cipher, err = WechatComponentSDK.GetCipher()

type request struct {
	rc io.ReadCloser
}

func (r *request) Read(p []byte) (n int, err error) {
	data, err := cipher.Decrypt(r.rc)
	n = copy(p, data)
	return
}

func (r *request) Close() error {
	return r.rc.Close()
}

type response struct {
	http.ResponseWriter
}

func (r *response) Write(p []byte) (n int, err error) {
	n = len(p)
	err = cipher.Encrypt(r.ResponseWriter, p)
	return
}

// 检测微信的签名
func CheckSign(w http.ResponseWriter, r *http.Request) {
	if cipher.CheckSignature(c.Writer, c.Request) {
		// check pass
	}
}


// 定时api接口
type RegularApi interface {
	GetAccessToken(ticket string) (string, float64)
	GetPreAuthCode(accessToken string) (string, float64)
}

func Job() {
	// 获取第三方平台 access token
	WechatComponentSDK.GetRegularApi().GetAccessToken(ticket)
	// 获取第三方平台 pre auth code
	WechatComponentSDK.GetRegularApi().GetPreAuthCode(accessToken)
}

func GetOauthUrl() string {
	// redirectUrl 公众号oauth授权后 callback url
	// preAuthCode 获取第三方平台 pre auth code
	return WechatComponentSDK.OAuthUrl(redirectUrl, preAuthCode)
}

// 通用api接口
type NormalApi interface {
 	// accessToken 第三方平台 access token 
 	// authCode 公众号oauth获取的授权码
	GetPublicInfo(accessToken, authCode string) (*PublicInfo, error)
	// appId 公众号 appId
	// refreshToken 公众号 刷新 access token 令牌
	GetAuthAccessToken(accessToken, appId, refreshToken string) (*PublicToken, error)
	GetAuthProfile(accessToken, appId string) (*PublicProfile, error)
	GetAuthOption(accessToken, appId, option string) (*PublicOption, error)
	SetAuthOption(accessToken, appId, optionName, optionValue string) error
}
func Call() {
	WechatComponentSDK.GetNormalApi().GetPublicInfo(accessToken, authCode)
	...
}

```
