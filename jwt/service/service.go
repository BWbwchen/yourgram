package jwt_svc

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

type AuthorizationService interface {
	CreateJWT(ctx context.Context, request AuthorizationRequest) AuthorizationResponse
	VerifyJWT(ctx context.Context, request AuthorizationRequest) AuthorizationResponse
}

type AuthorizationWorker struct{}

func NewService() AuthorizationService {
	return &AuthorizationWorker{}
}

type MyCustomClaims struct {
	UserData string `json:"userdata"`
	jwt.StandardClaims
}

var secretKey = []byte(os.Getenv("SECRETKEY"))

func (aw AuthorizationWorker) CreateJWT(ctx context.Context,
	request AuthorizationRequest) AuthorizationResponse {

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	claims := MyCustomClaims{
		request.UserData,
		jwt.StandardClaims{
			ExpiresAt: 0, // no expire time for dev
			Issuer:    "BWbwchen",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return AuthorizationResponse{
			StatusCode: http.StatusInternalServerError,
			Return:     "",
		}
	} else {
		return AuthorizationResponse{
			StatusCode: http.StatusOK,
			Return:     tokenString,
		}
	}
}

func (aw AuthorizationWorker) VerifyJWT(ctx context.Context,
	request AuthorizationRequest) AuthorizationResponse {

	userClaims := MyCustomClaims{}
	gotToken, err := jwt.ParseWithClaims(request.JWTToken, &userClaims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		fmt.Println(err)
	}

	if gotToken.Valid {
		return AuthorizationResponse{
			StatusCode: http.StatusOK,
			Return:     gotToken.Claims.(*MyCustomClaims).UserData,
		}
	} else {
		return AuthorizationResponse{
			StatusCode: http.StatusForbidden,
			Return:     "",
		}
	}
}
