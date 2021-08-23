package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UploadService struct {
	GateWayStruct
}

func NewUploadService() Service {
	as := UploadService{
		GateWayStruct: NewService(),
	}
	as.serviceName = "upload_service"
	as.tags = []string{"primary"}
	return as
}

type UploadRequest struct {
	UserID string `json:"userid"`
	Data   []byte `json:"data"`
	ImgID  string `json:"imgid"`
}

type UploadResponse struct {
	StatusCode int     `json:"status"`
	Info       ImgInfo `json:"info"`
}

type ImgInfo struct {
	User   string `json:"user"`
	ImgID  string `json:"imgid"`
	ImgURL string `json:"imgurl"`
}

var us Service = NewUploadService()

func UploadGateway(r *gin.Engine) *gin.Engine {
	service := r.Group("/v1/img")
	{
		service.POST("/upload", us.(UploadService).proxy)
		service.GET("/info", us.(UploadService).proxy)
	}
	return r
}

func (s UploadService) proxy(c *gin.Context) {
	s.originalRequest = c.Request
	s.data = getBodyData(c.Request)

	endpoint := s.GetEndpoint(c, s.encodeRequest, s.decodeResponse)

	var request UploadRequest
	json.NewDecoder(c.Request.Body).Decode(&request)

	response, _ := endpoint(c, request)

	resp := response.(UploadResponse)
	w := c.Writer
	w.WriteHeader(resp.StatusCode)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(resp)
}

func getBodyData(req *http.Request) []byte {
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return data
}

func (s UploadService) encodeRequest(ctx context.Context, request *http.Request, r interface{}) error {
	request.Body = ioutil.NopCloser(bytes.NewReader(s.data))
	for header, values := range s.originalRequest.Header {
		for _, value := range values {
			request.Header.Add(header, value)
		}
	}
	return nil
}

func (s UploadService) decodeResponse(ctx context.Context, resp *http.Response) (interface{}, error) {
	var response UploadResponse
	json.NewDecoder(resp.Body).Decode(&response)
	return response, nil
}
