package utils

import (
	"errors"
	"fmt"
	"football-manager-go/models/embedded"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"os"
	"time"
)

type customClaims struct {
	ID        string   `json:"id"`
	Username  string   `json:"username"`
	Firstname string   `json:"firstname"`
	Lastname  string   `json:"lastname"`
	Email     string   `json:"email"`
	Phone     string   `json:"phone"`
	Roles     []string `json:"roles"`
	jwt.StandardClaims
}

func GenerateJWT(id string, user embedded.User) (string, error) {
	claims := customClaims{
		ID:        id,
		Username:  user.Username,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Phone:     user.Phone,
		StandardClaims: jwt.StandardClaims{
			Issuer:    os.Getenv("JWT_VALID_ISSUER"),
			Id:        uuid.New().String(),
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		},
	}

	var roles []string

	for _, role := range user.Roles {
		if role.Name == embedded.Player {
			roles = append(roles, "player ")
		}
		if role.Name == embedded.Coach {
			roles = append(roles, "coach ")
		}
		if role.Name == embedded.Admin {
			roles = append(roles, "admin ")
		}
	}

	claims.Roles = roles

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func GetClaimsFromToken(tokenFromRequest string) (*customClaims, error) {
	token, err := jwt.ParseWithClaims(tokenFromRequest,
		&customClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

	if err != nil {
		fmt.Println("TOKEN EXPIRATION ERROR!!!!")
		fmt.Println(err.Error())
		return &customClaims{}, err
	}

	claims, ok := token.Claims.(*customClaims)
	if !ok {
		return &customClaims{}, errors.New("could not parse claims")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return &customClaims{}, errors.New("JWT Expired")
	}

	return claims, nil
}
