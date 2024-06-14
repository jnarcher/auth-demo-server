package model

import "net/http"

type ApiError struct {
	Error string `json:"error"`
}

type ApiFunc func(w http.ResponseWriter, r *http.Request) error
