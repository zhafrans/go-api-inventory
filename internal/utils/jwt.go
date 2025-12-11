package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	UserName string `json:"user_name"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID, email, userName, role, secret string, expireHours int) (string, error) {
	expirationTime := time.Now().Add(time.Duration(expireHours) * time.Hour)
	
	claims := &Claims{
		UserID: userID,
		Email:  email,
		UserName: userName,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "inventory-api",
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateJWT(tokenString, secret string) (*Claims, error) {
	claims := &Claims{}
	
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	
	return claims, nil
}