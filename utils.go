package wechat

import (
	"encoding/json"
	"net/http"
)

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
