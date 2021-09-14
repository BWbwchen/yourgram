package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"yourgram/pb"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"

	"google.golang.org/grpc"
)

type viewImageService struct {
	GateWayStruct
}

func NewViewImageService() Service {
	nvs := viewImageService{
		GateWayStruct: NewService(),
	}
	nvs.serviceName = "view_image_service"
	nvs.tags = []string{"primary"}
	fmt.Println(nvs)
	return nvs
}

type viewImageReq struct {
	UserID string `json:"UserID"`
}

type viewImageResp struct {
	ImgURLs []string `json:"ImgURLs"`
}

var nvs Service = NewViewImageService()

func ViewImageGateway(r *gin.Engine) *gin.Engine {
	service := r.Group("/api/v1/img")
	service.Use(JWTAuthMiddleWare())
	{
		service.GET("/getImage/:image_owner", nvs.(viewImageService).proxy)
	}
	return r
}

func (s viewImageService) GetEndpoint(c *gin.Context, enc httptransport.EncodeRequestFunc, dec httptransport.DecodeResponseFunc) endpoint.Endpoint {
	// get service instance
	instancer := consul.NewInstancer(s.client, s.logger, s.serviceName, s.tags, true)

	endpointer := sd.NewEndpointer(instancer, reqCreateGetImageFactory, s.logger)

	balancer := lb.NewRoundRobin(endpointer)
	endpoint, _ := balancer.Endpoint()

	return endpoint
}

func reqCreateGetImageFactory(instanceAddr string) (endpoint.Endpoint, io.Closer, error) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		conn, err := grpc.Dial(instanceAddr, grpc.WithInsecure())
		if err != nil {
			panic("connect error")
		}
		defer conn.Close()
		req := request.(viewImageReq)
		viewImageClient := pb.NewViewImageServiceClient(conn)
		resp, _ := viewImageClient.GetImage(context.Background(), &pb.ViewImageRequest{
			UserID: req.UserID,
		})
		return resp, nil
	}, nil, nil
}

func (s viewImageService) proxy(c *gin.Context) {
	endpoint := s.GetEndpoint(c, s.encodeRequest, s.decodeResponse)

	var request viewImageReq
	json.NewDecoder(c.Request.Body).Decode(&request)
	request.UserID = c.Param("image_owner")
	s.logger.Log("Welcome!", c.GetString("user"))
	response, _ := endpoint(c, request)

	resp := response.(*pb.ViewImageResponse)
	w := c.Writer
	//w.WriteHeader(int(resp.GetStatusCode()))
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(resp)
}

func (s viewImageService) encodeRequest(ctx context.Context, request *http.Request, in interface{}) error {
	return nil
}

func (s viewImageService) decodeResponse(ctx context.Context, resp *http.Response) (interface{}, error) {
	var response viewImageResp
	json.NewDecoder(resp.Body).Decode(&response)
	return response, nil
}
