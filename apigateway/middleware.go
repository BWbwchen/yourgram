package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

type JWTClaim struct {
	UserName string `json:"userdata"`
	jwt.StandardClaims
}

func JWTAuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			// Visitor
			c.Set("user", "__visitor__")
			c.Next()
			return
		}

		token := strings.Split(auth, "Bearer ")[1]
		jwtSecret := []byte(os.Getenv("SECRETKEY"))

		// parse and validate token for six things:
		// validationErrorMalformed => token is malformed
		// validationErrorUnverifiable => token could not be verified because of signing problems
		// validationErrorSignatureInvalid => signature validation failed
		// validationErrorExpired => exp validation failed
		// validationErrorNotValidYet => nbf validation failed
		// validationErrorIssuedAt => iat validation failed
		tokenClaims, err := jwt.ParseWithClaims(token, &JWTClaim{},
			func(token *jwt.Token) (i interface{}, err error) {
				return jwtSecret, nil
			})

		if err != nil {
			var message string
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					message = "token is malformed"
				} else if ve.Errors&jwt.ValidationErrorUnverifiable != 0 {
					message = "token could not be verified because of signing problems"
				} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
					message = "signature validation failed"
				} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
					message = "token is expired"
				} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
					message = "token is not yet valid before sometime"
				} else {
					message = "can not handle this token"
				}
			}
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": message,
			})
			c.Abort()
			return
		}

		if claims, ok := tokenClaims.Claims.(*JWTClaim); ok {
			c.Set("user", claims.UserName)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Something went wrong",
			})
			c.Abort()
			return
		}
	}
}
