package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

func Cors(next http.Handler) http.Handler {
	opts := cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowCredentials: true,
		Debug:            true,
	}
	c := cors.New(opts)
	return c.Handler(next)
}
