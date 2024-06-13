package model

import (
	"time"
)

type SignupRequest struct {
	User      string `json:"user"`
	Pwd       string `json:"pwd"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type LoginRequest struct {
	User      string `json:"user"`
	Pwd       string `json:"pwd"`
}

type Account struct {
	Id        int    `json:"id"`
	User      string `json:"user"`
	PwdHash   string `json:"-"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

