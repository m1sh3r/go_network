package main

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func main() {
	claims := Claims{
		Username: "user123",
		Role:     "user",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("my-secret-key"))
	if err != nil {
		fmt.Printf("Не удалось подписать токен: %v\n", err)
		return
	}

	fmt.Printf("JWT: %s\n", tokenString)
}
