package upload_endpoint

import (
	"context"

	upload_svc "yourgram/upload/service"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	Upload endpoint.Endpoint
	Info   endpoint.Endpoint
}
type Req struct {
	UserID string `json:"UserID"`
	Data   []byte `json:"Data"`
	ImgID  string `json:"ImgID"`
}

type Resp struct {
	StatusCode int     `json:"StatusCode"`
	Info       ImgInfo `json:"Info"`
}

type ImgInfo struct {
	User   string `json:"User"`
	ImgID  string `json:"ImgID"`
	ImgURL string `json:"ImgURL"`
}

func MakeEndpoints(s upload_svc.Service) Endpoints {
	return Endpoints{
		Upload: makeUploadEndPoint(s),
		Info:   makeInfoEndPoint(s),
	}
}

func makeUploadEndPoint(s upload_svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Req)
		res := s.Upload(ctx, upload_svc.UploadRequest{
			UserID: req.UserID,
			ImgID:  req.ImgID,
			Data:   req.Data,
		})
		return Resp{
			StatusCode: res.StatusCode,
			Info: ImgInfo{
				User:   res.Info.User,
				ImgID:  res.Info.ImgID,
				ImgURL: res.Info.ImgURL,
			},
		}, nil
	}
}

func makeInfoEndPoint(s upload_svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Req)
		res := s.Info(ctx, upload_svc.UploadRequest{
			UserID: req.UserID,
			ImgID:  req.ImgID,
			Data:   req.Data,
		})
		return Resp{
			StatusCode: res.StatusCode,
			Info: ImgInfo{
				User:   res.Info.User,
				ImgID:  res.Info.ImgID,
				ImgURL: res.Info.ImgURL,
			},
		}, nil
	}
}
