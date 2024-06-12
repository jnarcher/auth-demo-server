package middleware

import (
	hlp "auth-demo/internal/helpers"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("super-duper-secret-key")

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := w.Header().Get("Authorization")

		if tokenString == "" {
            hlp.SendError(w, "Missing auth token", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(tokenString, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
            hlp.SendError(w, "Invalid auth token", http.StatusUnauthorized)
			return
		}

		tokenString = tokenParts[1]

		claims, err := verifyToken(tokenString)
		if err != nil {
            hlp.SendError(w, "Invalid auth token", http.StatusUnauthorized)
			return
		}

		user, ok := claims["user"].(string)
		if !ok {
            hlp.SendError(w, "Unable to parse jwt", http.StatusUnauthorized)
			return
		}

		w.Header().Set("user", user)
		next.ServeHTTP(w, r)
	})
}

func CreateToken(user string, pwd string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": user,
			"password": pwd,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		},
	)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Printf("Unable to create jwt token: %v", err)
		return "", err
	}

	return tokenString, nil
}

func verifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if token.Method.Alg() != "HS256" {
			return nil, fmt.Errorf("Invalid signing method")
		}

		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("Invalid token")
}
