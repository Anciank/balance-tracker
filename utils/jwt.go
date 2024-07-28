package utils

import (
	"balance-tracker/models"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secret_key")

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

// NumericDate represents a numeric date value as defined in RFC 7519
type NumericDate int64

// NewNumericDate returns a new NumericDate value from a time.Time object
func NewNumericDate(t time.Time) int64 {
	return int64(t.Unix())
}

// Time returns the time.Time object corresponding to the NumericDate value
func (d NumericDate) Time() time.Time {
	return time.Unix(int64(d), 0)
}

func GenerateToken(user models.User) (string, error) {
	expirationTime := NewNumericDate(time.Now().Add(365 * 24 * time.Hour))
	claims := &Claims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
			Issuer:    "balance-tracker",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
