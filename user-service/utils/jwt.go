package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var jwtSecretKey []byte

func init() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// Initialize the secret key
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		log.Fatal("JWT_SECRET_KEY environment variable is not set")
	}
	jwtSecretKey = []byte(secret)
}

func GenerateJWT(userID uint) (string, error) {
	// Check if the secret key is set
	if len(jwtSecretKey) == 0 {
		return "", errors.New("JWT secret key is not set")
	}

	// Create the token claims
	claims := jwt.MapClaims{
		"authorized": true,
		"user_id":    userID,
		"exp":        time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract user ID from claims
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return 0, errors.New("user_id claim is invalid")
		}
		userID := uint(userIDFloat)
		return userID, nil
	} else {
		return 0, errors.New("invalid token claims")
	}
}
