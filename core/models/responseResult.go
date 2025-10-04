package models

import (
	"fmt"
)

type ResponseResult struct {
	code    int
	message string
}

func NewResponseResult(code int, message string) ResponseResult {
	return ResponseResult{
		code:    code,
		message: message,
	}
}

func (r *ResponseResult) String() string {
	return fmt.Sprintf("Code: %d | Message %s", r.code, r.message)
}
