package wechat

import (
	"fmt"

	"github.com/franela/goreq"
	"github.com/parnurzeal/gorequest"
)

type NormalApi interface {
	GetPublicInfo(accessToken, authCode string) (*PublicInfo, error)
	GetAuthAccessToken(accessToken, appId, refreshToken string) (*PublicToken, error)
	GetAuthProfile(accessToken, appId string) (*PublicProfile, error)
	GetAuthOption(accessToken, appId, option string) (*PublicOption, error)
	SetAuthOption(accessToken, appId, optionName, optionValue string) error
}

type normalApi struct {
	wt      *WechatThird
	request *gorequest.SuperAgent
}

type PublicInfo struct {
	AuthorizationInfo struct {
		AppId        string  `json:"authorizer_appid"`
		AccessToken  string  `json:"authorizer_access_token"`
		ExpiresIn    float64 `json:"expires_in"`
		RefreshToken string  `json:"authorizer_refresh_token"`
		FuncInfo     []struct {
			Funcscope struct {
				Id int64 `json:"id"`
			} `json:"funcscope_category"`
		} `json:"func_info"`
	} `json:"authorization_info"`
}

func (na *normalApi) GetPublicInfo(accessToken, authCode string) (*PublicInfo, error) {
	postData := struct {
		Component_appid    string `json:"component_appid"`
		Authorization_code string `json:"authorization_code"`
	}{
		Component_appid:    na.wt.appId,
		Authorization_code: authCode,
	}
	res, err := goreq.Request{
		Method: "POST",
		Uri:    fmt.Sprintf(apiQueryAuth, accessToken),
		Body:   postData,
	}.Do()
	if err != nil {
		return nil, err

	}
	result := &PublicInfo{}

	err = unmarshalResponseToJson(res, result)
	if err != nil {
		return nil, err

	}
	return result, nil
}

// authorizer access token and refresh token
type PublicToken struct {
	AccessToken  string  `json:"authorizer_access_token"`
	ExpiresIn    float64 `json:"expires_in"`
	RefreshToken string  `json:"authorizer_refresh_token"`
}

// accessToken component access token
// appId authorizer appId
// refreshToken authorizer refresh token
func (na *normalApi) GetAuthAccessToken(accessToken, appId, refreshToken string) (*PublicToken, error) {

	postData := struct {
		Component_appid          string `json:"component_appid"`
		Authorizer_appid         string `json:"authorizer_appid"`
		Authorizer_refresh_token string `json:"authorizer_refresh_token"`
	}{
		Component_appid:          na.wt.appId,
		Authorizer_appid:         appId,
		Authorizer_refresh_token: refreshToken,
	}

	res, err := goreq.Request{
		Method: "POST",
		Uri:    fmt.Sprintf(apiAuthorizerToken, accessToken),
		Body:   postData,
	}.Do()
	if err != nil {
		return nil, err

	}

	result := &PublicToken{}

	err = unmarshalResponseToJson(res, result)
	if err != nil {
		return nil, err

	}
	return result, nil
}

type PublicProfile struct {
	AuthorizerInfo struct {
		NickName        string `json:"nick_name"`
		HeadImg         string `json:"head_img"`
		ServiceTypeInfo struct {
			Id int64 `json:"id"`
		}
		VerifyTypeInfo struct {
			Id int64 `json:"id"`
		}
		UserName string `json:"user_name"`
		Alias    string `json:"alias"`
	} `json:"authorizer_info"`
	QR                string `json:"qrcode_url"`
	AuthorizationInfo struct {
		AppId    string `json:"appid"`
		FuncInfo []struct {
			Funcscope struct {
				Id int64 `json:"id"`
			} `json:"funcscope_category"`
		} `json:"func_info"`
	} `json:"authorization_info"`
}

func (na *normalApi) GetAuthProfile(accessToken, appId string) (*PublicProfile, error) {

	postForm := struct {
		Component_appid  string `json:"component_appid"`
		Authorizer_appid string `json:"authorizer_appid"`
	}{
		Component_appid:  na.wt.appId,
		Authorizer_appid: appId,
	}

	res, err := goreq.Request{
		Method: "POST",
		Uri:    fmt.Sprintf(apiGetAuthorizerInfo, accessToken),
		Body:   postForm,
	}.Do()
	if err != nil {
		return nil, err

	}
	result := &PublicProfile{}

	err = unmarshalResponseToJson(res, result)
	if err != nil {
		return nil, err

	}
	return result, nil
}

type PublicOption struct {
	AppId       string `json:"authorizer_appid"`
	OptionName  string `json:"option_name"`
	OptionValue string `json:"option_value"`
}

func (na *normalApi) GetAuthOption(accessToken, appId, option string) (*PublicOption, error) {
	postForm := struct {
		Component_appid  string `json:"component_appid"`
		Authorizer_appid string `json:"authorizer_appid"`
		Option_name      string `json:"option_name"`
	}{
		Component_appid:  na.wt.appId,
		Authorizer_appid: appId,
		Option_name:      option,
	}

	res, err := goreq.Request{
		Method: "POST",
		Uri:    fmt.Sprintf(apiGetAuthorizerOption, accessToken),
		Body:   postForm,
	}.Do()
	if err != nil {
		return nil, err

	}
	result := &PublicOption{}

	err = unmarshalResponseToJson(res, result)
	if err != nil {
		return nil, err

	}
	return result, nil
}

func (na *normalApi) SetAuthOption(accessToken, appId, optionName, optionValue string) error {
	postForm := struct {
		Component_appid  string `json:"component_appid"`
		Authorizer_appid string `json:"authorizer_appid"`
		Option_name      string `json:"option_name"`
		Option_value     string `json:"option_value"`
	}{
		Component_appid:  na.wt.appId,
		Authorizer_appid: appId,
		Option_name:      optionName,
		Option_value:     optionValue,
	}

	res, err := goreq.Request{
		Method: "POST",
		Uri:    fmt.Sprintf(apiSetAuthorizerOption, accessToken),
		Body:   postForm,
	}.Do()
	if err != nil {
		return err

	}
	result := &ApiError{}

	return unmarshalResponseToJson(res, result)

}
