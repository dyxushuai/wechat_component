// regular tasks api
// these api will be run on time
package wechat

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

type RegularApi interface {
	GetAccessToken(ticket string) (string, time.Duration)
	GetPreAuthCode(accessToken string) (string, time.Duration)
}

type regularApi struct {
	wt *WechatThird
}

// type accessToken struct {
// 	createTime time.Time     // when refresh access token
// 	expire     time.Duration // expire duration
// }

// // 定时循环判断token 有没有超过到期时间的一半
// func (wt *WechatThird) createAccessOnTime() {
// 	token := accessToken{time.Now(), 1}

// 	for _ = range time.Tick(5 * time.Second) {
// 		if time.Since(token.createTime) >= token.expire/2 {
// 			access, expire := wt.GetAccessToken()
// 			if access != "" {
// 				token.createTime = time.Now()
// 				token.expire = expire
// 				wt.accessLock.Lock()
// 				wt.accessToken = access
// 				wt.accessLock.Unlock()
// 			}
// 		}
// 	}
// }

// 获取第三方平台令牌
func (ra *regularApi) GetAccessToken(ticket string) (string, time.Duration) {
	postForm := url.Values{}
	postForm.Set("component_appid", ra.wt.appId)
	postForm.Set("component_appsecret", ra.wt.appSecret)
	postForm.Set("component_verify_ticket", ticket)
	res, err := http.PostForm(apiComponentToken, postForm)

	result := &struct {
		CAT       string `json:"component_access_token"`
		ExpiresIn int64  `json:"expires_in"`
	}{}

	ae, err := unmarshalResponseToJson(res, result)
	if err != nil {
		log.Println("Parse access token failed: ", err)

	} else if ae != nil {
		if ae.isError() {
			log.Println("getAccessToken api failed: ", ae.ErrMsg)
		}
	} else {
		return result.CAT, time.Duration(result.ExpiresIn * 1000 * 1000 * 1000)
	}
	return "", 0
}

// 获取预授权码，用于公众号oauth
func (ra *regularApi) GetPreAuthCode(accessToken string) (string, time.Duration) {
	postForm := url.Values{}
	postForm.Set("component_appid", ra.wt.appId)
	res, err := http.PostForm(fmt.Sprintf(apiCreatePreAuthCode, accessToken), postForm)

	result := &struct {
		PAC       string `json:"pre_auth_code"`
		ExpiresIn int64  `json:"expires_in"`
	}{}

	ae, err := unmarshalResponseToJson(res, result)
	if err != nil {
		log.Println("Parse pre auth token failed: ", err)

	} else if ae != nil {
		if ae.isError() {
			log.Println("getPreAuthCode api failed: ", ae.ErrMsg)
		}
	} else {
		return result.PAC, time.Duration(result.ExpiresIn * 1000 * 1000 * 1000)
	}
	return "", 0
}
