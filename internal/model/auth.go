package model

import (
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type AuthClaims struct {
	AccountId int64    `json:"accountId"`
	User      string `json:"user"`
	PwdHash   string `json:"pwdHash"`

	jwt.RegisteredClaims
}

func (ac AuthClaims) Validate() error {
	if strings.TrimSpace(ac.User) == "" {
		return errors.New("no user found")
	}
	if strings.TrimSpace(ac.PwdHash) == "" {
		return errors.New("no pwd hash found")
	}
	return nil
}
