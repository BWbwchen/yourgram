package upload_endpoint

import (
	"context"

	upload_svc "yourgram/upload/service"

	"github.com/go-kit/kit/endpoint"
)

func MakeUploadEndPoint(s upload_svc.UploadService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(upload_svc.UploadRequest)
		res := s.Upload(ctx, req)
		return res, nil
	}
}

func MakeInfoEndPoint(s upload_svc.UploadService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(upload_svc.UploadRequest)
		res := s.Info(ctx, req)
		return res, nil
	}
}
