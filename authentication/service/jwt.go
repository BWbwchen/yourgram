package authentication_service

import (
	"os"

	"github.com/golang-jwt/jwt"
)

type MyCustomClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

var secretKey = []byte(os.Getenv("SECRETKEY"))

func generateJWTToken(user UserInfo) string {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	claims := MyCustomClaims{
		user.Email,
		jwt.StandardClaims{
			ExpiresAt: 2592000, // 1 month
			Issuer:    "BWbwchen",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, _ := token.SignedString(secretKey)
	Logger.Log("JWTToken", tokenString)

	return tokenString
}
