package upload_svc

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testUploadCase = []struct {
	req            UploadRequest
	expectedStatus int
}{
	{
		req: UploadRequest{
			UserID: "bwbwchen",
			Data:   getJPG(),
			ImgID:  "",
		},
		expectedStatus: http.StatusOK,
	},
	{
		req: UploadRequest{
			UserID: "bwbwchen",
			Data:   getPNG(),
			ImgID:  "",
		},
		expectedStatus: http.StatusOK,
	},
	{
		req: UploadRequest{
			UserID: "bwbwchen",
			Data:   getErrorImg(),
			ImgID:  "",
		},
		expectedStatus: http.StatusBadRequest,
	},
}

func TestUpload(t *testing.T) {
	auth := NewService()
	for _, testcase := range testUploadCase {
		resp := auth.Upload(context.TODO(), testcase.req)
		assert.Equal(t, testcase.expectedStatus, resp.StatusCode, resp)
	}
}

var testInfoCase = []struct {
	req            UploadRequest
	expectedStatus int
	msg            string
}{
	{
		req: UploadRequest{
			UserID: "bwbwchen",
			Data:   nil,
			ImgID:  "123456",
		},
		expectedStatus: http.StatusBadRequest,
		msg:            "not exist image id",
	},
	{
		req: UploadRequest{
			UserID: "bwbwchen",
			Data:   nil,
			ImgID:  "test",
		},
		expectedStatus: http.StatusOK,
		msg:            "exist image id",
	},
}

func TestInfo(t *testing.T) {
	auth := NewService()
	for _, testcase := range testInfoCase {
		resp := auth.Info(context.TODO(), testcase.req)
		assert.Equal(t, testcase.expectedStatus, resp.StatusCode, testcase.msg)
	}
}

func getJPG() []byte {
	jpgImg, _ := ioutil.ReadFile("../lena.jpg")
	return jpgImg
}

func getPNG() []byte {
	pngImg, _ := ioutil.ReadFile("../lena.png")
	return pngImg
}

func getErrorImg() []byte {
	errorImg, _ := ioutil.ReadFile("../main.go")
	return errorImg
}

func init() {
	addPerData()
}
