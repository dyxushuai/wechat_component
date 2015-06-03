package wechat

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"time"
)

var (
	chars       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	defaultRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// xml cdata
type CharData struct {
	Text []byte `xml:",innerxml"`
}

func NewCharData(s string) CharData {
	return CharData{[]byte("<![CDATA[" + s + "]]>")}
}

// 生成nonce
func createNonceStr(length int) string {
	var str string
	for i := 0; i < length; i++ {
		tmpI := defaultRand.Intn(len(chars) - 1)
		str += chars[tmpI : tmpI+1]
	}
	return str
}

func aesKeyDecode(encodedAESKey string) (key []byte, err error) {
	if len(encodedAESKey) != 43 {
		err = errors.New("the length of encodedAESKey must be equal to 43")
		return
	}
	key, err = base64.StdEncoding.DecodeString(encodedAESKey + "=")
	if err != nil {
		return
	}
	if len(key) != 32 {
		err = errors.New("encodingAESKey invalid")
		return
	}
	return
}

func unmarshalResponseToJson(res *http.Response, v interface{}) (*ApiError, error) {
	defer res.Body.Close()

	d := json.NewDecoder(res.Body)

	apiErr := &ApiError{}

	err := d.Decode(apiErr)
	if err != nil {
		return nil, err
	}
	if apiErr.isError() {
		return apiErr, nil
	}

	return nil, d.Decode(v)
}
