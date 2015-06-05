package wechat

import "github.com/franela/goreq"

func unmarshalResponseToJson(res *goreq.Response, v interface{}) (*ApiError, error) {
	apiErr := &ApiError{}
	err := res.Body.FromJsonTo(apiErr)
	if err != nil {
		return nil, err
	}
	if apiErr.isError() {
		return apiErr, nil
	}
	return nil, res.Body.FromJsonTo(v)
}
