// 授权事件处理器

package wechat

import (
	"encoding/xml"
	"errors"
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

type AuthHandler interface {
	M(http.ResponseWriter, *http.Request) (*AuthResult, error)
}

func (ar *AuthResult) IsTicket() bool {
	return ar.InfoType == "component_verify_ticket"
}

type authHandle struct {
	wt     *WechatThird
	cipher lib.IOChipher
}

func newAuthHandle(wt *WechatThird) (*authHandle, error) {
	c, err := lib.NewCipher(wt.token, wt.cryptoKey, wt.appId)
	if err != nil {
		return nil, err
	}
	return &authHandle{
		wt:     wt,
		cipher: c,
	}, nil
}

func (ah *authHandle) M(w http.ResponseWriter, r *http.Request) (*AuthResult, error) {
	defer w.Write([]byte("success"))

	if !lib.CheckSignature(ah.wt.token, w, r) {
		return nil, errors.New("sign error form weixin paltform")
	}
	data, err := ah.cipher.Decrypt(r.Body)
	if err != nil {
		return nil, err
	}

	ar := &AuthResult{}

	err = xml.Unmarshal(data, ar)
	if err != nil {
		return nil, err
	}

	return ar, nil
}
