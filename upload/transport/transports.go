package upload_transport

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	upload_svc "yourgram/upload/service"
)

func DecodePOSTRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	file, _, err := req.FormFile("file")
	if err != nil {
		log.Println(err)
		return upload_svc.UploadRequest{}, err
	}
	data, _ := ioutil.ReadAll(file)

	request := upload_svc.UploadRequest{
		UserID: "test",
		ImgID:  "",
		Data:   data,
	}

	return request, nil

}

func DecodeGETRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var request upload_svc.UploadRequest
	if req.Body == nil {
		fmt.Println(req.Body)
	}
	err := json.NewDecoder(req.Body).Decode(&request)
	if err == io.EOF {
		fmt.Println("decode eof")
		return upload_svc.UploadRequest{}, nil
	}
	return request, nil
}

func DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	if req.Method == "POST" {
		return DecodePOSTRequest(ctx, req)
	} else if req.Method == "GET" {
		return DecodeGETRequest(ctx, req)
	}
	return upload_svc.UploadRequest{}, nil
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(upload_svc.UploadResponse)
	w.WriteHeader(resp.StatusCode)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(resp)
	return nil
}
