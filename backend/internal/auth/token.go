package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrUnauthorized = errors.New("unauthorized")

type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// Tokens issues and validates JWTs.
type Tokens struct {
	secret string
	ttl    time.Duration
}

func NewTokens(secret string) *Tokens {
	return &Tokens{secret: secret, ttl: 720 * time.Hour} // 30 days
}

func (t *Tokens) Issue(userID, role string) (string, error) {
	now := time.Now()
	claims := Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(t.ttl)),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(t.secret))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return token, nil
}

func (t *Tokens) Parse(raw string) (User, error) {
	var claims Claims
	_, err := jwt.ParseWithClaims(raw, &claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return []byte(t.secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return User{}, ErrUnauthorized
	}
	if claims.Subject == "" || claims.Role == "" {
		return User{}, ErrUnauthorized
	}
	return User{ID: claims.Subject, Role: claims.Role}, nil
}

func (t *Tokens) pepper() string {
	return t.secret
}
