package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"yourgram/pb"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
	"google.golang.org/grpc"
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
	service.Use(JWTAuthMiddleWare())
	{
		service.POST("/upload", us.(UploadService).proxy)
		service.GET("/info", us.(UploadService).proxy)
	}
	return r
}
func (s UploadService) GetEndpoint(c *gin.Context, enc httptransport.EncodeRequestFunc, dec httptransport.DecodeResponseFunc) endpoint.Endpoint {
	// get service instance
	instancer := consul.NewInstancer(s.client, s.logger, s.serviceName, s.tags, true)
	fmt.Println(instancer)

	relativePathSlice := strings.Split(c.Request.URL.Path, "/")
	relativePath := relativePathSlice[len(relativePathSlice)-1]
	var endpointer *sd.DefaultEndpointer
	if relativePath == "upload" {
		endpointer = sd.NewEndpointer(instancer, reqUploadFactory, s.logger)
	} else if relativePath == "info" {
		fmt.Println("match auth")
		endpointer = sd.NewEndpointer(instancer, reqInfoFactory, s.logger)
	} else {
		fmt.Println("No match method")
		return nil
	}

	balancer := lb.NewRoundRobin(endpointer)
	endpoint, _ := balancer.Endpoint()

	/*
		endpoints, err := endpointer.Endpoints()
		endpoint := endpoints[0]
	*/
	return endpoint
}
func reqUploadFactory(instanceAddr string) (endpoint.Endpoint, io.Closer, error) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		conn, err := grpc.Dial(instanceAddr, grpc.WithInsecure())
		if err != nil {
			panic("connect error")
		}
		defer conn.Close()
		req := request.(UploadRequest)
		accountClient := pb.NewUploadServiceClient(conn)
		resp, _ := accountClient.Upload(context.Background(), &pb.UploadRequest{
			UserID: req.UserID,
			Data:   req.Data,
			ImgID:  req.ImgID,
		})
		return resp, nil
	}, nil, nil
}

func reqInfoFactory(instanceAddr string) (endpoint.Endpoint, io.Closer, error) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		conn, err := grpc.Dial(instanceAddr, grpc.WithInsecure())
		if err != nil {
			panic("connect error")
		}
		defer conn.Close()
		req := request.(UploadRequest)
		accountClient := pb.NewUploadServiceClient(conn)
		resp, _ := accountClient.Info(context.Background(), &pb.UploadRequest{
			UserID: req.UserID,
			Data:   req.Data,
			ImgID:  req.ImgID,
		})
		return resp, nil
	}, nil, nil
}

func (s UploadService) proxy(c *gin.Context) {
	s.originalRequest = c.Request
	s.data = getBodyData(c.Request)

	endpoint := s.GetEndpoint(c, s.encodeRequest, s.decodeResponse)

	var request UploadRequest
	json.NewDecoder(c.Request.Body).Decode(&request)
	request.Data = s.data
	// TODO: use JWT token user id
	request.UserID = c.GetString("user")

	response, _ := endpoint(c, request)

	resp := response.(*pb.UploadResponse)
	w := c.Writer
	w.WriteHeader(int(resp.GetStatusCode()))
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(resp)
}

func getBodyData(req *http.Request) []byte {
	file, _, _ := req.FormFile("file")
	data, err := ioutil.ReadAll(file)
	defer file.Close()
	if err != nil {
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
