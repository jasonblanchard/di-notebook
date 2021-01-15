package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func bearerHeaderToID(header string) (string, error) {
	tokenString := strings.Replace(header, "Bearer ", "", 1)

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return "", nil
	})

	if token == nil {
		return "", errors.New("Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok != true {
		return "", errors.New("Token does not contain any claims")
	}

	return fmt.Sprintf("%s", claims["uesrUuid"]), nil
}
