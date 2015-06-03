// 授权事件处理器

package wechat

import (
	"encoding/xml"
	"log"
	"net/http"
)

type AuthHandler interface {
	http.Handler
	GetTicket() string // 定时获取ticket
}

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

type AuthHandle struct {
	token      string
	cipher     IOChipher
	ticketChan chan string
}

func NewAuthHandle(token, encodingAESKey, appID string) (*AuthHandle, error) {
	c, err := NewCipher(token, encodingAESKey, appID)
	if err != nil {
		return nil, err
	}
	return &AuthHandle{
		token:      token,
		cipher:     c,
		ticketChan: make(chan string),
	}, nil
}

func (ah *AuthHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 微信要求返回 success
	defer w.Write([]byte("success"))

	if !checkSignature(ah.token, w, r) {
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

func (ah *AuthHandle) GetTicket() string {
	return <-ah.ticketChan
}
