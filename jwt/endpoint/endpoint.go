package jwt_endpoint

import (
	"context"

	jwt_svc "yourgram/jwt/service"

	"github.com/go-kit/kit/endpoint"
)

func MakeCreateJWTEndPoint(s jwt_svc.AuthorizationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(jwt_svc.AuthorizationRequest)
		res := s.CreateJWT(ctx, req)
		return res, nil
	}
}

func MakeVerifyJWTEndPoint(s jwt_svc.AuthorizationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(jwt_svc.AuthorizationRequest)
		res := s.VerifyJWT(ctx, req)
		return res, nil
	}
}
