package services

import "objectswaterfall.com/core/models"

type Send interface {
	SendRequest(host string, obj interface{}, headers map[string]string) (models.ResponseResult, error)
}
