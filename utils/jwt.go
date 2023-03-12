package utils

import (
	"errors"
	"fmt"
	"os"
	"time"
	
	"github.com/dgrijalva/jwt-go"
)

type customClaims struct {
	ID       string `json:"id"`
	UserType string `json:"type"`
	jwt.StandardClaims
}

func GenerateJWTToken(ID, UserType string) (string, error) {
	claims := customClaims{
		ID:       ID,
		UserType: UserType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
			Issuer:    "football-manager-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func GetClaimsFromToken(tokenFromHeader string) (string, string, error) {
	token, err := jwt.ParseWithClaims(
		tokenFromHeader,
		&customClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		},
	)

	if err != nil {
		fmt.Println("TOKEN EXPIRATION ERROR!!!!")
		fmt.Println(err.Error())
		return "", "", err
	}

	claims, ok := token.Claims.(*customClaims)
	if !ok {
		return "", "", errors.New("could not parse claims")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return "", "", errors.New("JWT Expired")
	}

	return claims.ID, claims.UserType, nil
}
