package account_endpoint

import (
	"context"

	account_svc "yourgram/account/service"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateAccount endpoint.Endpoint
	UserLogin     endpoint.Endpoint
}

type Req struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Resp struct {
	StatusCode int    `json:"status"`
	JWTToken   string `json:"jwt"`
}

func MakeEndpoints(s account_svc.Service) Endpoints {
	return Endpoints{
		CreateAccount: makeCreateAccountEndpoint(s),
		UserLogin:     makeUserLoginEndpoint(s),
	}
}

func makeCreateAccountEndpoint(s account_svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Req)
		res := s.CreateAccount(ctx, account_svc.Input{
			Email:    req.Email,
			Name:     req.Name,
			Password: req.Password,
		})
		return Resp{
			StatusCode: res.StatusCode,
			JWTToken:   res.JWTToken,
		}, nil
	}
}

func makeUserLoginEndpoint(s account_svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Req)
		res := s.UserLogin(ctx, account_svc.Input{
			Email:    req.Email,
			Name:     req.Name,
			Password: req.Password,
		})
		return Resp{
			StatusCode: res.StatusCode,
			JWTToken:   res.JWTToken,
		}, nil
	}
}
