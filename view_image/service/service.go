package service

import (
	"context"
	"yourgram/view_image/pb"

	"github.com/go-kit/kit/log"
)

type Worker struct {
	log log.Logger
}

func NewService(log log.Logger) pb.ViewImageServiceServer {
	initService()
	return &Worker{
		log: log,
	}
}

func initService() {
	initDB()
}

func (w Worker) GetImage(ctx context.Context, request *pb.ViewImageRequest) (*pb.ViewImageResponse, error) {
	w.log.Log("status", "get image")
	ImgURLs, err := db.Query(request.UserID)

	return &pb.ViewImageResponse{
		ImgURLs: ImgURLs,
	}, err
}
