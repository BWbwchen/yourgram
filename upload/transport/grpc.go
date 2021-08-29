package upload_trans

import (
	"context"
	upload_endp "yourgram/upload/endpoint"
	"yourgram/upload/pb"

	"github.com/go-kit/kit/log"

	gt "github.com/go-kit/kit/transport/grpc"
)

type gRPCServer struct {
	upload gt.Handler
	info   gt.Handler
}

func NewGRPCServer(endpoints upload_endp.Endpoints, logger log.Logger) pb.UploadServiceServer {
	return &gRPCServer{
		upload: gt.NewServer(
			endpoints.Upload,
			decodeRequestRPC,
			encodeResponseRPC,
		),
		info: gt.NewServer(
			endpoints.Info,
			decodeRequestRPC,
			encodeResponseRPC,
		),
	}
}

func (s *gRPCServer) Upload(ctx context.Context, req *pb.UploadRequest) (*pb.UploadResponse, error) {
	_, resp, err := s.upload.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.UploadResponse), nil
}

func (s *gRPCServer) Info(ctx context.Context, req *pb.UploadRequest) (*pb.UploadResponse, error) {
	_, resp, err := s.info.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.UploadResponse), nil
}

func decodeRequestRPC(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.UploadRequest)
	return upload_endp.Req{
		UserID: req.UserID,
		Data:   req.Data,
		ImgID:  req.ImgID,
	}, nil
}

func encodeResponseRPC(ctx context.Context, response interface{}) (interface{}, error) {
	resp := response.(upload_endp.Resp)
	return &pb.UploadResponse{
		StatusCode: int32(resp.StatusCode),
		Info: &pb.ImgInfo{
			User:   resp.Info.User,
			ImgID:  resp.Info.ImgID,
			ImgURL: resp.Info.ImgURL,
		},
	}, nil
}
