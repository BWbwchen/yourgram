package main

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func MakeCreateJWTEndPoint(s AuthorizationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AuthorizationRequest)
		res := s.CreateJWT(ctx, req)
		return res, nil
	}
}

func MakeVerifyJWTEndPoint(s AuthorizationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AuthorizationRequest)
		res := s.VerifyJWT(ctx, req)
		return res, nil
	}
}
