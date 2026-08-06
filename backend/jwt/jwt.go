package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Login  string `json:"login"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func HashPassword(password string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("bcrypt hash: %w", err)
	}
	return string(h), nil
}

func ComparePassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func GenerateToken(secret string, expiryHours int, userID uint, login, role string) (string, error) {
	if secret == "" {
		return "", errors.New("jwt: empty secret")
	}
	if expiryHours <= 0 {
		return "", errors.New("jwt: expiryHours must be > 0")
	}

	now := time.Now()
	claims := Claims{
		UserID: userID,
		Login:  login,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expiryHours) * time.Hour)),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   login,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("jwt sign: %w", err)
	}
	return signed, nil
}

func ValidateToken(secret string, tokenString string) (*Claims, error) {
	if secret == "" {
		return nil, errors.New("jwt: empty secret")
	}
	if tokenString == "" {
		return nil, errors.New("jwt: empty token")
	}

	claims := &Claims{}
	parsed, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("jwt: unexpected signing method %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("jwt parse: %w", err)
	}
	if !parsed.Valid {
		return nil, errors.New("jwt: token is not valid")
	}
	return claims, nil
}
