package url_filter

import (
	"bytes"
	"encoding/json"
)

type CheckUrlRequest struct {
	Url string `json:"url"`
}

type CheckUrlResponse struct {
	Result bool `json:"exists"`
}

func CreateCheckUrlReq(url string, buf *bytes.Buffer) error {
	rq := CheckUrlRequest{
		Url: url,
	}
	err := json.NewEncoder(buf).Encode(rq)
	return err
}
