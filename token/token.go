package token

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

func GenerateToken(apiSecret, apiKey string) string {
	var jwtKey = []byte(apiSecret)
	expirationTime := time.Now().Add(20 * time.Minute)
	claims := &jwt.StandardClaims{
		Audience:  " ",
		Issuer:    apiKey,
		ExpiresAt: expirationTime.Unix(),
		IssuedAt:  time.Now().Unix(),
	}
	tokenJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, _ := tokenJwt.SignedString(jwtKey)
	return token
}
