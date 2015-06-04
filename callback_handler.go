// 公众号消息与事件处理器

package wechat

import "net/http"

type cbHandle struct {
}

func (cb *cbHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
