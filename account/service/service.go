package account_svc

import (
	"context"
	"net/http"

	"github.com/go-kit/log"
)

type Service interface {
	CreateAccount(ctx context.Context, request Input) Output
	UserLogin(ctx context.Context, request Input) Output
}

type Worker struct {
	logger log.Logger
}

func NewService(logger log.Logger) Service {
	initService()
	return &Worker{
		logger: logger,
	}
}

func initService() {
	initDB()
}

func (aw Worker) CreateAccount(ctx context.Context,
	request Input) Output {
	if db.CreateUser(UserInfo(request)) {
		return Output{
			StatusCode: http.StatusOK,
			JWTToken:   "",
		}
	} else {
		return Output{
			StatusCode: http.StatusBadRequest,
			JWTToken:   "",
		}
	}
}

func (aw Worker) UserLogin(ctx context.Context,
	request Input) Output {
	JWTToken := db.UserLogin(UserInfo(request))
	if JWTToken != "" {
		return Output{
			StatusCode: http.StatusOK,
			JWTToken:   JWTToken,
		}
	} else {
		return Output{
			StatusCode: http.StatusUnauthorized,
			JWTToken:   JWTToken,
		}
	}
}
