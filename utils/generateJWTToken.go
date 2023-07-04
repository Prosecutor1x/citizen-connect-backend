package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

func GenerateJWT(phone string) string {

	godotenv.Load(".env.local")

	// Define the secret key used to sign the token

	secretKey := os.Getenv("JWT_SECRET")

	// Create a new token object
	token := jwt.New(jwt.SigningMethodHS256)

	// Set the claims for the token
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = phone
	claims["exp"] = time.Now().Add(time.Hour * 4).Unix() // Token expiration time (1 hour from now)

	// Generate the JWT string
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		fmt.Println("Error generating token:", err)
		return "Error"
	}

	// Print the JWT string
	return tokenString
}
