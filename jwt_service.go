package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type JWTService struct {
	GateWayStruct
}

func NewJWTService() Service {
	as := JWTService{
		GateWayStruct: NewService(),
	}
	as.serviceName = "jwt_service"
	as.tags = []string{"primary"}
	fmt.Println(as)
	return as
}

type AuthorizationRequest struct {
	UserData string `json:"userdata"`
	JWTToken string `json:"jwt"`
}

type AuthorizationResponse struct {
	StatusCode int    `json:"status"`
	Return     string `json:"return"`
}

var js Service = NewJWTService()

func JWTGateway(r *gin.Engine) *gin.Engine {
	service := r.Group("/v1/jwt")
	{
		service.POST("/create", js.(JWTService).proxy)
		service.GET("/auth", js.(JWTService).proxy)
	}
	return r
}

func (s JWTService) proxy(c *gin.Context) {
	fmt.Println(c.Request)
	endpoint := s.GetEndpoint(c, s.encodeRequest, s.decodeResponse)

	var request AuthorizationRequest
	json.NewDecoder(c.Request.Body).Decode(&request)
	response, _ := endpoint(c, request)

	resp := response.(AuthorizationResponse)
	w := c.Writer
	w.WriteHeader(resp.StatusCode)
	encoder := json.NewEncoder(w)
	encoder.Encode(resp)
}

func (s JWTService) encodeRequest(ctx context.Context, request *http.Request, in interface{}) error {
	body, _ := json.Marshal(in)
	newRequest, _ := http.NewRequest(request.Method, "", bytes.NewBuffer(body))
	request.Body = newRequest.Body
	return nil
}

func (s JWTService) decodeResponse(ctx context.Context, resp *http.Response) (interface{}, error) {
	var response AuthorizationResponse
	json.NewDecoder(resp.Body).Decode(&response)
	return response, nil
}
