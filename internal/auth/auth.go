package auth

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte("super-duper-secret-key")

func CreateToken(user string, pwdHash string, role string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user": user,
            "pwdHash": pwdHash,
			"role": role,
			"exp":  time.Now().Add(time.Hour * 24).Unix(),
		},
	)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Printf("Unable to create jwt token: %v", err)
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
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

func HashPassword(pwd string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), 14)
    return string(bytes), err
}

func CheckPasswordHash(pwd, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
    return err == nil
}
