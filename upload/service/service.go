package upload_svc

import (
	"context"
	"net/http"
)

type UploadService interface {
	Upload(ctx context.Context, request UploadRequest) UploadResponse
	Info(ctx context.Context, request UploadRequest) UploadResponse
}

type UploadWorker struct{}

func NewService() UploadService {
	return &UploadWorker{}
}

func init() {
	initDB()
	initCloudStorage()
}

func (uw UploadWorker) Info(ctx context.Context, request UploadRequest) UploadResponse {
	info, err := getImgInfo(request.ImgID)
	if err != nil {
		return UploadResponse{
			StatusCode: http.StatusBadRequest,
			Info:       ImgInfo{},
		}
	}
	return UploadResponse{
		StatusCode: http.StatusOK,
		Info:       info,
	}
}

func (uw UploadWorker) Upload(ctx context.Context, request UploadRequest) UploadResponse {
	fileType := http.DetectContentType(request.Data)
	if fileType != "image/jpeg" && fileType != "image/png" {
		return UploadResponse{
			StatusCode: http.StatusBadRequest,
			Info:       ImgInfo{},
		}
	}

	info, err := storeImg(ctx, request.Data)
	if err != nil {
		panic(err)
	}
	info.User = request.UserID

	err = storeImgInfo(info)
	if err != nil {
		panic(err)
	}

	return UploadResponse{
		StatusCode: http.StatusOK,
		Info:       info,
	}
}
