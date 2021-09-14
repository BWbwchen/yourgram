package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

type accountService struct {
	GateWayStruct
}

func NewAccountService() Service {
	as := accountService{
		GateWayStruct: NewService(),
	}
	as.serviceName = "account_service"
	as.tags = []string{"primary"}
	fmt.Println(as)
	return as
}

type AuthRequest struct {
	Email    string `json:"Email"`
	Name     string `json:"Name"`
	Password string `json:"Password"`
}

type AuthResponse struct {
	StatusCode int    `json:"StatusCode"`
	JWTToken   string `json:"JWTToken"`
	Email      string `json:"Email"`
	Name       string `json:"Name"`
}

var as Service = NewAccountService()

func AccountGateway(r *gin.Engine) *gin.Engine {
	service := r.Group("/api/v1/account")
	{
		service.POST("/create", as.(accountService).proxy)
		service.POST("/auth", as.(accountService).proxy)
	}
	return r
}

func (s accountService) GetEndpoint(c *gin.Context, enc httptransport.EncodeRequestFunc, dec httptransport.DecodeResponseFunc) endpoint.Endpoint {
	// get service instance
	instancer := consul.NewInstancer(s.client, s.logger, s.serviceName, s.tags, true)
	fmt.Println(instancer)

	relativePathSlice := strings.Split(c.Request.URL.Path, "/")
	relativePath := relativePathSlice[len(relativePathSlice)-1]
	var endpointer *sd.DefaultEndpointer
	if relativePath == "create" {
		endpointer = sd.NewEndpointer(instancer, reqCreateAccountFactory, s.logger)
	} else if relativePath == "auth" {
		fmt.Println("match auth")
		endpointer = sd.NewEndpointer(instancer, reqUserLoginFactory, s.logger)
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

func reqCreateAccountFactory(instanceAddr string) (endpoint.Endpoint, io.Closer, error) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		conn, err := grpc.Dial(instanceAddr, grpc.WithInsecure())
		if err != nil {
			panic("connect error")
		}
		defer conn.Close()
		req := request.(AuthRequest)
		accountClient := pb.NewAuthServiceClient(conn)
		resp, _ := accountClient.CreateAccount(context.Background(), &pb.AuthRequest{Email: req.Email, Name: req.Name, Password: req.Password})
		return resp, nil
	}, nil, nil
}

func reqUserLoginFactory(instanceAddr string) (endpoint.Endpoint, io.Closer, error) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		conn, err := grpc.Dial(instanceAddr, grpc.WithInsecure())
		if err != nil {
			panic("connect error")
		}
		defer conn.Close()
		req := request.(AuthRequest)
		accountClient := pb.NewAuthServiceClient(conn)
		resp, _ := accountClient.UserLogin(context.Background(), &pb.AuthRequest{Email: req.Email, Name: req.Name, Password: req.Password})
		return resp, nil
	}, nil, nil
}

func (s accountService) proxy(c *gin.Context) {
	endpoint := s.GetEndpoint(c, s.encodeRequest, s.decodeResponse)

	var request AuthRequest
	json.NewDecoder(c.Request.Body).Decode(&request)
	response, _ := endpoint(c, request)

	resp := response.(*pb.AuthResponse)
	w := c.Writer
	w.WriteHeader(int(resp.GetStatusCode()))
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(resp)
}

func (s accountService) encodeRequest(ctx context.Context, request *http.Request, in interface{}) error {
	body, _ := json.Marshal(in)
	bodyReader := bytes.NewReader(body)
	request.Body = io.NopCloser(bodyReader)
	return nil
}

func (s accountService) decodeResponse(ctx context.Context, resp *http.Response) (interface{}, error) {
	var response AuthResponse
	json.NewDecoder(resp.Body).Decode(&response)
	return response, nil
}
