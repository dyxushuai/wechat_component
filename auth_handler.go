// 授权事件处理器

package wechat

import (
	"encoding/xml"
	"log"
	"net/http"

	"git.ishopex.cn/xushuai/wechat/lib"
)

type AuthResult struct {
	AppId                 string
	CreateTime            string
	InfoType              string
	ComponentVerifyTicket string
	AuthorizerAppid       string
}

func (ar *AuthResult) isTicket() bool {
	return ar.InfoType == "component_verify_ticket"
}

type authHandle struct {
	token      string
	cipher     lib.IOChipher
	ticketChan chan string
}

func newAuthHandle(token, key, appID string) (*authHandle, error) {
	c, err := lib.NewCipher(token, key, appID)
	if err != nil {
		return nil, err
	}
	return &authHandle{
		token:      token,
		cipher:     c,
		ticketChan: make(chan string),
	}, nil
}

func (ah *authHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 微信要求返回 success
	defer w.Write([]byte("success"))

	if !lib.CheckSignature(ah.token, w, r) {
		log.Println("sign error form weixin paltform")
		return
	}
	data, err := ah.cipher.Decrypt(r.Body)
	if err != nil {
		log.Println(err)
		return
	}

	ar := &AuthResult{}

	err = xml.Unmarshal(data, ar)
	if err != nil {
		log.Println(err)
		return
	}
	if ar.isTicket() {
		ah.ticketChan <- ar.ComponentVerifyTicket
	}
}

func (ah *authHandle) getTicket() <-chan string {
	return ah.ticketChan
}
