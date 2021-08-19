package authentication

import (
	"context"

	authsvc "yourgram/authentication/service"

	"github.com/go-kit/kit/endpoint"
)

func MakeCreateAccountEndPoint(s authsvc.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(authsvc.AuthRequest)
		res := s.CreateAccount(ctx, req)
		return res, nil
	}
}

func MakeUserLoginEndPoint(s authsvc.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(authsvc.AuthRequest)
		res := s.UserLogin(ctx, req)
		return res, nil
	}
}
