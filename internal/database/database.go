package database

import (
	"auth-demo/internal/auth"
	"auth-demo/internal/model"
	"time"
)

type DB interface {
	CreateAccount(*model.Account) error
	DeleteAccount(int) error
	UpdateAccount(*model.Account) error
    GetAccounts() ([]*model.Account, error)
	GetAccountById(int) (*model.Account, error)
}

func NewAccount(req model.SignupRequest) (*model.Account, error) {
    hash, err := auth.HashPassword(req.Pwd)
    if err != nil {
        return nil, err
    }

	return &model.Account{
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
