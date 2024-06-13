package auth

import (
	"auth-demo/internal/model"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func CreateToken(acc *model.Account) (string, error) {
	secretKey, err := getSecretKey()
	if err != nil {
		return "", err
	}

	claims := model.AuthClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
            IssuedAt: jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer: "auth-demo-server",
		},
		AccountId: acc.Id,
		User:      acc.User,
		PwdHash:   acc.PwdHash,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Printf("Unable to create jwt token: %v", err)
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (*model.AuthClaims, error) {
	secretKey, err := getSecretKey()
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(
        tokenString, 
        &model.AuthClaims{}, 
        func(token *jwt.Token) (interface{}, error) {
            if token.Method.Alg() != "HS256" {
                return nil, fmt.Errorf("Invalid signing method")
            }
            return secretKey, nil
        },
    )
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*model.AuthClaims); ok && token.Valid {
		if err := claims.Validate(); err != nil {
			return nil, err
		}
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

func SetAuthCookie(w http.ResponseWriter, acc *model.Account) error {
	tkn, err := CreateToken(acc)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    tkn,
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, cookie)
	return nil
}

func getSecretKey() ([]byte, error) {
	secretKey, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		return []byte{}, fmt.Errorf("JWT secret key not set in env")
	}
	return []byte(secretKey), nil
}
