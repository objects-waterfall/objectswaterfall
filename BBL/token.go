package bbl

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"objectswaterfall.com/core/errors"
)

// Make functionality wider. Should give more then one auth option
type TokenService struct {
	authUrl      string
	authModel    string
	token        string
	expireTime   time.Time
	expiresAfter time.Duration
}

func NewTokenService(authUrl, authModel string, expiresAfter time.Duration) TokenService {
	return TokenService{
		authUrl:      authModel,
		authModel:    authModel,
		expireTime:   time.Now().Add(expiresAfter),
		expiresAfter: expiresAfter,
	}
}

func (t *TokenService) GetTokenFromUrl() (string, error) {
	resp, err := http.Post(t.authUrl, "application/json", bytes.NewBufferString(t.authModel))
	if err != nil {
		return "", errors.NewTockenRecievingError(err.Error())
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("code: %d. message: %s", resp.StatusCode, resp.Status)
	}
	token, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.NewTockenRecievingError(err.Error())
	}

	return string(token), nil
}

func (t *TokenService) Token() (string, error) {
	var err error
	if t.expireTime.Before(time.Now()) {
		t.token, err = t.GetTokenFromUrl()
		if err != nil {
			return "", err
		}
	}
	return t.token, nil
}
