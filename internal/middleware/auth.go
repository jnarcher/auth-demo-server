package middleware

import (
	"auth-demo/internal/auth"
	hlp "auth-demo/internal/helpers"
	"errors"
	"net/http"
	// "strings"
)

var secretKey = []byte("super-duper-secret-key")

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        c, err := r.Cookie("token")
		if err != nil {
            switch {
            case errors.Is(err, http.ErrNoCookie):
                hlp.SendError(w, "Auth cookie not found", http.StatusBadRequest)
            default:
                hlp.SendError(w, "server error", http.StatusInternalServerError)
            }
			return
		}

        tokenString := c.Value
		// tokenParts := strings.Split(tokenString, " ")
		// if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		// 	hlp.SendError(w, "Invalid auth token", http.StatusUnauthorized)
		// 	return
		// }

		// tokenString = tokenParts[1]

		claims, err := auth.VerifyToken(tokenString)
		if err != nil {
			hlp.SendError(w, "Invalid auth token", http.StatusUnauthorized)
			return
		}

		user, ok := claims["user"].(string)
		if !ok {
			hlp.SendError(w, "Unable to parse jwt", http.StatusUnauthorized)
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			hlp.SendError(w, "Unable to parse jwt", http.StatusUnauthorized)
			return
		}

		pwdHash, ok := claims["pwdHash"].(string)
		if !ok {
			hlp.SendError(w, "Unable to parse jwt", http.StatusUnauthorized)
			return
		}

        r.Header.Add("user", user)
        r.Header.Add("role", role)
        r.Header.Add("pwdHash", pwdHash)
		next.ServeHTTP(w, r)
	})
}
