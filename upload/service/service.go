package upload_svc

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/log"
)

type Service interface {
	Upload(ctx context.Context, request UploadRequest) UploadResponse
	Info(ctx context.Context, request UploadRequest) UploadResponse
}

type UploadWorker struct {
	log log.Logger
}

func NewService(log log.Logger) Service {
	initService()
	return &UploadWorker{
		log: log,
	}
}

func initService() {
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
	if !checkImageType(request.Data) {
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

func checkImageType(data []byte) bool {
	fileType := http.DetectContentType(data)
	if fileType != "image/jpeg" && fileType != "image/png" {
		return false
	}
	return true
}
