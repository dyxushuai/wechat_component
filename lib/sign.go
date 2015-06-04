package lib

import (
	"crypto/sha1"
	"encoding/hex"
	"math/rand"
	"net/http"
	"sort"
	"time"
)

var (
	chars       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	defaultRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// 生成nonce
func createNonceStr(length int) string {
	var str string
	for i := 0; i < length; i++ {
		tmpI := defaultRand.Intn(len(chars) - 1)
		str += chars[tmpI : tmpI+1]
	}
	return str
}

// 微信公众号 明文模式/URL认证 签名
func Sign(token, timestamp, nonce string) (signature string) {
	strs := sort.StringSlice{token, timestamp, nonce}
	strs.Sort()

	n := len(token) + len(timestamp) + len(nonce)
	buf := make([]byte, 0, n)

	buf = append(buf, strs[0]...)
	buf = append(buf, strs[1]...)
	buf = append(buf, strs[2]...)

	hashsum := sha1.Sum(buf)
	return hex.EncodeToString(hashsum[:])
}

// 微信公众号/企业号 密文模式消息签名
func MsgSign(token, timestamp, nonce, encryptedMsg string) (signature string) {
	strs := sort.StringSlice{token, timestamp, nonce, encryptedMsg}
	strs.Sort()

	n := len(token) + len(timestamp) + len(nonce) + len(encryptedMsg)
	buf := make([]byte, 0, n)

	buf = append(buf, strs[0]...)
	buf = append(buf, strs[1]...)
	buf = append(buf, strs[2]...)
	buf = append(buf, strs[3]...)

	hashsum := sha1.Sum(buf)
	return hex.EncodeToString(hashsum[:])
}

func CheckSignature(t string, w http.ResponseWriter, r *http.Request) bool {
	r.ParseForm()
	signature := r.FormValue("signature")
	timestamp := r.FormValue("timestamp")
	nonce := r.FormValue("nonce")

	return Sign(t, timestamp, nonce) == signature
}
