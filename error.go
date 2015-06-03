// 微信错误返回
package wechat

type ApiError struct {
	ErrCode string `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (a *ApiError) isError() bool {
	return a.ErrCode != ""
}
