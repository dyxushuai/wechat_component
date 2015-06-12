// 微信错误返回
package wechat

import (
	"fmt"
)

type ApiError struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (a *ApiError) isError() bool {
	return a.ErrMsg != ""
}

func (a *ApiError) Error() string {
	return fmt.Sprintf("Api Error: Code: %d Message: %s", a.ErrCode, a.ErrMsg)
}
