package account_endpoint

import (
	"context"

	account_svc "yourgram/account/service"

	"github.com/go-kit/kit/endpoint"
)

func MakeCreateAccountEndPoint(s account_svc.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(account_svc.AuthRequest)
		res := s.CreateAccount(ctx, req)
		return res, nil
	}
}

func MakeUserLoginEndPoint(s account_svc.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(account_svc.AuthRequest)
		res := s.UserLogin(ctx, req)
		return res, nil
	}
}
