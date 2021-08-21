package token

import (
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

type Claims struct {
	jwt.StandardClaims
}

func TokenGenerate(apiSecret string, apiKey string) string {
	var jwtKey = []byte(apiSecret)
	_, present := os.LookupEnv("ZOOM_TOKEN")
	claims := &Claims{}
	if present {
		tknStr := os.Getenv("ZOOM_TOKEN")

		tkn, _ := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if tkn.Valid {
			return tknStr
		}
		return tknStr
	}

	apikey := apiKey
	expirationTime := time.Now().Add(20 * time.Minute)
	claims = &Claims{
		StandardClaims: jwt.StandardClaims{
			Audience:  " ",
			Issuer:    apikey,
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)
	os.Setenv("ZOOM_TOKEN", tokenString)
	return tokenString
}
