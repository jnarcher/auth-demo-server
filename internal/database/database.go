package database

import "auth-demo/internal/model"

type DB interface {
	CreateAccount(*model.Account) error
	DeleteAccount(int) error
	UpdateAccount(*model.Account) error
    GetAccounts() ([]*model.Account, error)
	GetAccountById(int) (*model.Account, error)
}

