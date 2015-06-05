// 微信错误返回
package wechat

type ApiError struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (a *ApiError) isError() bool {
	return a.ErrMsg != ""
}
