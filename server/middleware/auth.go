package middleware

import (
	"context"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/metadata"
)

var jwtKey = []byte("my_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func Authenticate(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("missing metadata")
	}
	token := md["authorization"]
	if len(token) == 0 {
		return "", errors.New("missing token")
	}
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token[0], claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !tkn.Valid {
		return "", errors.New("invalid token")
	}
	return claims.Username, nil
}
