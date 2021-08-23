package upload_svc

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client
var bucketName string = "yourgram"

func initCloudStorage() {
	endpoint := os.Getenv("minio_url")
	accessKeyID := os.Getenv("minio_user")
	secretAccessKey := os.Getenv("minio_password")
	useSSL := false

	mc, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Println(err)
	}
	minioClient = mc

	// Make a new bucket called mymusic.
	bucketName := "yourgram"

	ctx := context.TODO()
	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Println(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

}

func storeImg(ctx context.Context, data []byte) (info ImgInfo, err error) {
	info = ImgInfo{
		ImgID:  uuid.NewString(),
		ImgURL: "",
	}
	err = nil

	contentType := http.DetectContentType(data)

	// Upload the zip file with PutObject
	_, err = minioClient.PutObject(ctx, bucketName, info.ImgID, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Println(err)
		return
	}

	// Generates a presigned url which expires in a day.
	presignedURL, err := minioClient.PresignedGetObject(context.Background(), bucketName, info.ImgID, time.Second*10, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	info.ImgURL = presignedURL.String()
	log.Println("Successfully generated presigned URL", presignedURL.String())
	return
}
