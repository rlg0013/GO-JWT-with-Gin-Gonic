package helper

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type SignedDetails struct {
	Email, FirstName, LastName, UID, UserType string
	jwt.RegisteredClaims
}

func secret() ([]byte, error) {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		return nil, fmt.Errorf("JWT_SECRET is not set")
	}
	return []byte(s), nil
}
func GenerateAllTokens(email, firstName, lastName, uid, userType string) (string, string, error) {
	s, err := secret()
	if err != nil {
		return "", "", err
	}
	claims := SignedDetails{Email: email, FirstName: firstName, LastName: lastName, UID: uid, UserType: userType, RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), IssuedAt: jwt.NewNumericDate(time.Now())}}
	refreshClaims := claims
	refreshClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour))
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s)
	if err != nil {
		return "", "", err
	}
	refresh, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(s)
	return token, refresh, err
}
func ValidateToken(tokenString string) (*SignedDetails, error) {
	s, err := secret()
	if err != nil {
		return nil, err
	}
	claims := &SignedDetails{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s, nil
	})
	if err != nil || !token.Valid {
		if err == nil {
			err = fmt.Errorf("invalid token")
		}
		return nil, err
	}
	return claims, nil
}
