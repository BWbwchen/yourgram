package health

import (
	"context"
	"yourgram/view_image/pb"

	"github.com/go-kit/log"
)

type service struct {
	log log.Logger
}

func NewService(log log.Logger) pb.HealthServer {
	return &service{
		log: log,
	}
}

func (s *service) Check(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	s.log.Log("status", "check")
	return &pb.HealthCheckResponse{
		Status: pb.HealthCheckResponse_SERVING,
	}, nil
}
