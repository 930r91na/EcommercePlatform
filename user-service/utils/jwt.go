package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"time"
)

var jwtKey = os.Getenv("JWT_KEY")

func GenerateJWT(userID uint) (string, error) {
	// Token works for a day
	expirationTime := time.Now().Add(72 * time.Hour)

	// Create the JWT claims, which includes the user ID and expiry time
	claims := &jwt.RegisteredClaims{
		// Convert uint to string for the Subject claim
		Subject:   fmt.Sprintf("%d", userID),
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "your-application-name",
	}

	// Create the token using the claims and sign it with the secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
