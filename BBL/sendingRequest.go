package bbl

import (
	"bytes"
	"encoding/json"
	"net/http"

	"objectswaterfall.com/core/models"
	"objectswaterfall.com/core/services"
)

type sendingService struct {
}

func NewSendingService() services.Send {
	return sendingService{}
}

func (s sendingService) SendRequest(host string, obj interface{}, headers map[string]string) (models.ResponseResult, error) {
	var err error
	var buffer *bytes.Buffer
	if str, ok := obj.(string); ok {
		buffer = bytes.NewBufferString(str)
	} else {
		data, err := json.Marshal(obj)
		buffer = bytes.NewBuffer(data)
		if err != nil {
			return models.ResponseResult{}, err
		}
	}
	req, err := http.NewRequest("POST", host, buffer)
	if err != nil {
		return models.ResponseResult{}, err
	}
	defer req.Body.Close()
	if len(headers) != 0 {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.ResponseResult{}, err
	}

	return models.NewResponseResult(resp.StatusCode, resp.Status), nil
}
