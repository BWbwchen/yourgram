package account_trans

import (
	"context"
	account_endp "yourgram/account/endpoint"
	"yourgram/account/pb"

	"github.com/go-kit/kit/log"

	gt "github.com/go-kit/kit/transport/grpc"
)

type gRPCServer struct {
	createAccount gt.Handler
	userLogin     gt.Handler
}

func NewGRPCServer(endpoints account_endp.Endpoints, logger log.Logger) pb.AuthServiceServer {
	return &gRPCServer{
		createAccount: gt.NewServer(
			endpoints.CreateAccount,
			decodeRequestRPC,
			encodeResponseRPC,
		),
		userLogin: gt.NewServer(
			endpoints.UserLogin,
			decodeRequestRPC,
			encodeResponseRPC,
		),
	}
}

func (s *gRPCServer) CreateAccount(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	_, resp, err := s.createAccount.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.AuthResponse), nil
}

func (s *gRPCServer) UserLogin(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	_, resp, err := s.userLogin.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.AuthResponse), nil
}

func decodeRequestRPC(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.AuthRequest)
	return account_endp.Req{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	}, nil
}

func encodeResponseRPC(ctx context.Context, response interface{}) (interface{}, error) {
	resp := response.(account_endp.Resp)
	return &pb.AuthResponse{
		StatusCode: int32(resp.StatusCode),
		JWTToken:   resp.JWTToken,
	}, nil
}
