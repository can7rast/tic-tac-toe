package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"log"
	"os"
	"time"
)

type Provider struct {
	secretKey   []byte
	accessTime  time.Duration
	refreshTime time.Duration
}

func NewProvider() *Provider {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		log.Fatal("JWT_SECRET environment variable not set")
		return nil
	}

	return &Provider{
		secretKey:   []byte(secretKey),
		accessTime:  time.Duration(time.Hour),
		refreshTime: time.Duration(time.Hour * 24),
	}
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Login  string    `json:"login"`
	jwt.RegisteredClaims
}

func (p *Provider) GenerateAccessToken(userID uuid.UUID, login string) (string, error) {
	claims := Claims{
		UserID: userID,
		Login:  login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(p.accessTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "candidas",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(p.secretKey)
}

func (p *Provider) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(p.refreshTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "candidas",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(p.secretKey)
}

func (p *Provider) ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return p.secretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

func (p *Provider) GetUserId(tokenStr string) (uuid.UUID, error) {
	claims, err := p.ValidateToken(tokenStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid token: %v", err)
	}
	return claims.UserID, nil
}
