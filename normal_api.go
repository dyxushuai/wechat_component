package wechat

import (
	"fmt"
	"net/http"
	"net/url"
)

type NormalApi interface {
	GetPublicInfo(accessToken, authCode string) (*PublicInfo, *ApiError, error)
	GetAuthAccessToken(accessToken, appId, refreshToken string) (*PublicToken, *ApiError, error)
	GetAuthProfile(accessToken, appId string) (*PublicProfile, *ApiError, error)
	GetAuthOption(accessToken, appId, option string) (*PublicOption, *ApiError, error)
	SetAuthOption(accessToken, appId, optionName, optionValue string) (*ApiError, error)
}

type normalApi struct {
	wt *WechatThird
}

type PublicInfo struct {
	AuthorizationInfo struct {
		AppId        string `json:"authorizer_appid"`
		AccessToken  string `json:"authorizer_access_token"`
		ExpiresIn    int64  `json:"expires_in"`
		RefreshToken string `json:"authorizer_refresh_token"`
		FuncInfo     []struct {
			Funcscope struct {
				Id int64 `json:"id"`
			} `json:"funcscope_category"`
		} `json:"func_info"`
	} `json:"authorization_info"`
}

func (na *normalApi) GetPublicInfo(accessToken, authCode string) (*PublicInfo, *ApiError, error) {
	postForm := url.Values{}
	postForm.Set("component_appid", na.wt.appId)
	postForm.Set("authorization_code", authCode)
	res, err := http.PostForm(fmt.Sprintf(apiQueryAuth, accessToken), postForm)
	if err != nil {
		return nil, nil, err

	}
	result := &PublicInfo{}

	ae, err := unmarshalResponseToJson(res, result)
	if err != nil {
		return nil, nil, err

	}
	if ae != nil {
		return nil, ae, nil
	}
	return result, nil, nil
}

// authorizer access token and refresh token
type PublicToken struct {
	AccessToken  string `json:"authorizer_access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"authorizer_refresh_token"`
}

// accessToken component access token
// appId authorizer appId
// refreshToken authorizer refresh token
func (na *normalApi) GetAuthAccessToken(accessToken, appId, refreshToken string) (*PublicToken, *ApiError, error) {
	postForm := url.Values{}
	postForm.Set("component_appid", na.wt.appId)
	postForm.Set("authorizer_appid", appId)
	postForm.Set("authorizer_refresh_token", refreshToken)
	res, err := http.PostForm(fmt.Sprintf(apiAuthorizerToken, accessToken), postForm)
	if err != nil {
		return nil, nil, err

	}
	result := &PublicToken{}

	ae, err := unmarshalResponseToJson(res, result)
	if err != nil {
		return nil, nil, err

	}
	if ae != nil {
		return nil, ae, nil
	}
	return result, nil, nil
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

func (na *normalApi) GetAuthProfile(accessToken, appId string) (*PublicProfile, *ApiError, error) {
	postForm := url.Values{}
	postForm.Set("component_appid", na.wt.appId)
	postForm.Set("authorizer_appid", appId)
	res, err := http.PostForm(fmt.Sprintf(apiGetAuthorizerInfo, accessToken), postForm)
	if err != nil {
		return nil, nil, err

	}
	result := &PublicProfile{}

	ae, err := unmarshalResponseToJson(res, result)
	if err != nil {
		return nil, nil, err

	}
	if ae != nil {
		return nil, ae, nil
	}
	return result, nil, nil
}

type PublicOption struct {
	AppId       string `json:"authorizer_appid"`
	OptionName  string `json:"option_name"`
	OptionValue string `json:"option_value"`
}

func (na *normalApi) GetAuthOption(accessToken, appId, option string) (*PublicOption, *ApiError, error) {
	postForm := url.Values{}
	postForm.Set("component_appid", na.wt.appId)
	postForm.Set("authorizer_appid", appId)
	postForm.Set("option_name", option)
	res, err := http.PostForm(fmt.Sprintf(apiGetAuthorizerOption, accessToken), postForm)
	if err != nil {
		return nil, nil, err

	}
	result := &PublicOption{}

	ae, err := unmarshalResponseToJson(res, result)
	if err != nil {
		return nil, nil, err

	}
	if ae != nil {
		return nil, ae, nil
	}
	return result, nil, nil
}

func (na *normalApi) SetAuthOption(accessToken, appId, optionName, optionValue string) (*ApiError, error) {
	postForm := url.Values{}
	postForm.Set("component_appid", na.wt.appId)
	postForm.Set("authorizer_appid", appId)
	postForm.Set("option_name", optionName)
	postForm.Set("option_value", optionValue)
	res, err := http.PostForm(fmt.Sprintf(apiSetAuthorizerOption, accessToken), postForm)
	if err != nil {
		return nil, err

	}
	result := &ApiError{}

	return unmarshalResponseToJson(res, result)

}
