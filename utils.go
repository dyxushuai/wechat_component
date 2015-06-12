package wechat

import (
	"encoding/json"
	"io/ioutil"
	"net/url"

	"github.com/franela/goreq"
)

func unmarshalResponseToJson(res *goreq.Response, v interface{}) error {
	b, err := ioutil.ReadAll(res.Body)
	apiErr := &ApiError{}

	err = json.Unmarshal(b, apiErr)

	if err != nil {
		return err
	}
	if apiErr.isError() {
		return apiErr
	}
	return json.Unmarshal(b, v)
}

func UrlEncoded(str string) (string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
