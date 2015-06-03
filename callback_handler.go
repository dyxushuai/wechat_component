// 公众号消息与事件处理器

package wechat

import (
	"net/http"
)

type CBHandler interface {
	http.Handler
}

type CBHandle struct {
}
