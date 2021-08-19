package authentication_service

import (
	"context"
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
			StatusCode: 200,
			JWTToken:   "",
		}
	} else {
		return AuthResponse{
			StatusCode: 400,
			JWTToken:   "",
		}
	}
}

func (aw AuthenticateWorker) UserLogin(ctx context.Context,
	request AuthRequest) AuthResponse {
	JWTToken := db.UserLogin(UserInfo(request))
	if JWTToken != "" {
		return AuthResponse{
			StatusCode: 200,
			JWTToken:   JWTToken,
		}
	} else {
		return AuthResponse{
			StatusCode: 401,
			JWTToken:   JWTToken,
		}
	}
}
