package pkg

import (
	"encoding/base64"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
)

var cost = 12

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(hash), err
}

func CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func ParseBasicAuth(authHeader string) (login, password string, err error) {
	if authHeader == "" || !strings.HasPrefix(authHeader, "Basic ") {
		return "", "", errors.New("invalid auth header")
	}

	decoded, err := base64.StdEncoding.DecodeString(authHeader[6:])
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	creds := strings.SplitN(string(decoded), ":", 2)
	if len(creds) != 2 {
		return "", "", err
	}

	return creds[0], creds[1], nil
}
