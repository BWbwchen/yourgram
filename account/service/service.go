package account_svc

import (
	"context"
	"net/http"

	"yourgram/account/pb"

	"github.com/go-kit/log"
)

type Worker struct {
	logger log.Logger
}

func NewService(logger log.Logger) pb.AuthServiceServer {
	initService()
	return &Worker{
		logger: logger,
	}
}

func initService() {
	initDB()
}

func (aw Worker) CreateAccount(ctx context.Context,
	request *pb.AuthRequest) (*pb.AuthResponse, error) {
	if db.CreateUser(UserInfo{
		request.Email,
		request.Name,
		request.Password,
	}) {
		return &pb.AuthResponse{
			StatusCode: http.StatusOK,
			JWTToken:   "",
			Email:      request.Email,
			Name:       request.Name,
		}, nil
	} else {
		return &pb.AuthResponse{
			StatusCode: http.StatusBadRequest,
			JWTToken:   "",
			Email:      request.Email,
			Name:       request.Name,
		}, nil
	}
}

func (aw Worker) UserLogin(ctx context.Context,
	request *pb.AuthRequest) (*pb.AuthResponse, error) {
	JWTToken, userInfo := db.UserLogin(UserInfo{
		request.Email,
		request.Name,
		request.Password,
	})
	if JWTToken != "" {
		return &pb.AuthResponse{
			StatusCode: http.StatusOK,
			JWTToken:   JWTToken,
			Email:      userInfo.Email,
			Name:       userInfo.Name,
		}, nil
	} else {
		return &pb.AuthResponse{
			StatusCode: http.StatusUnauthorized,
			JWTToken:   JWTToken,
			Email:      "",
			Name:       "",
		}, nil
	}
}
