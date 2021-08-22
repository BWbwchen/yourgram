package account_svc

import (
	"context"
	"net/http"
)

type AuthService interface {
	CreateAccount(ctx context.Context, request AuthRequest) AuthResponse
	UserLogin(ctx context.Context, request AuthRequest) AuthResponse
}

type AuthenticateWorker struct{}

func NewService() AuthService {
	return &AuthenticateWorker{}
}

func (aw AuthenticateWorker) CreateAccount(ctx context.Context,
	request AuthRequest) AuthResponse {
	if db.CreateUser(UserInfo(request)) {
		return AuthResponse{
			StatusCode: http.StatusOK,
			JWTToken:   "",
		}
	} else {
		return AuthResponse{
			StatusCode: http.StatusBadRequest,
			JWTToken:   "",
		}
	}
}

func (aw AuthenticateWorker) UserLogin(ctx context.Context,
	request AuthRequest) AuthResponse {
	JWTToken := db.UserLogin(UserInfo(request))
	if JWTToken != "" {
		return AuthResponse{
			StatusCode: http.StatusOK,
			JWTToken:   JWTToken,
		}
	} else {
		return AuthResponse{
			StatusCode: http.StatusUnauthorized,
			JWTToken:   JWTToken,
		}
	}
}
