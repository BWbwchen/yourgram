package main

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func MakeCreateAccountEndPoint(s AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AuthRequest)
		res := s.CreateAccount(ctx, req)
		return res, nil
	}
}

func MakeUserLoginEndPoint(s AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AuthRequest)
		res := s.UserLogin(ctx, req)
		return res, nil
	}
}
