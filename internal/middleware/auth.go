package middleware

import (
	"auth-demo/internal/auth"
	"auth-demo/internal/helpers"
	"auth-demo/internal/model"
	"errors"
	"fmt"
	"log"

	// "fmt"
	"net/http"
)

func WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				_ = helpers.WriteJson(
					w,
					http.StatusBadRequest,
					model.ApiError{Error: "permision denied"},
				)
			default:
				_ = helpers.WriteJson(
					w,
					http.StatusInternalServerError,
					model.ApiError{Error: "server error"},
				)
			}
			return
		}

		tokenString := c.Value
        claims, err := auth.VerifyToken(tokenString)
		if err != nil {
            log.Printf("Error - unable to verify jwt token: %+v", err)
			_ = helpers.WriteJson(
				w,
				http.StatusUnauthorized,
				model.ApiError{Error: "permission denied"},
			)
			return
		}

        r.Header.Add("account_id", fmt.Sprintf("%d", claims.AccountId))
        r.Header.Add("account_user", claims.User)
        r.Header.Add("account_pwd_hash", claims.PwdHash)
		next.ServeHTTP(w, r)
	})
}
