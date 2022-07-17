package auth

import (
	"context"
	"errors"
	"time"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	jwtSecret = []byte{}
)

const (
	bearer = "bearer"
	issuer = "grin-authur"
	expire = 8 * time.Hour
)

var (
	ErrUnauthorized error = errors.New("unauthorized")
	ErrExpired      error = errors.New("token is expired")
)

type Claims struct {
	jwt.Claims
	Username string
	Ip       string
}

// Generate 生成token
func Generate(username string) (string, error) {
	deadline := time.Now().Add(expire)
	claims := Claims{
		Claims: jwt.StandardClaims{
			ExpiresAt: deadline.Unix(),
			Issuer:    issuer,
		},
		Username: username,
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

// ParseToken 解析token
func ParseToken(ctx context.Context) (*Claims, error) {
	token, err := grpc_auth.AuthFromMD(ctx, bearer)
	if err != nil {
		return nil, err
	}

	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, ErrUnauthorized
	}
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, ErrUnauthorized
}
