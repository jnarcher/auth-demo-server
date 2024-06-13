package model

import (
	"auth-demo/internal/auth"
	"time"
)

type CreateAccountRequest struct {
	User      string `json:"user"`
	Pwd       string `json:"pwd"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type Account struct {
	Id        int    `json:"id"`
	User      string `json:"user"`
	PwdHash   string `json:"pwdHash"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

func NewAccount(req CreateAccountRequest) (*Account, error) {
    hash, err := auth.HashPassword(req.Pwd)
    if err != nil {
        return nil, err
    }

	return &Account{
		User:      req.User,
		PwdHash:   hash,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
	}, nil
}
