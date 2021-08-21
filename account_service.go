package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type AuthResponse struct {
	StatusCode int    `json:"status"`
	JWTToken   string `json:"jwt"`
}

var as Service = NewAccountService()

func AccountGateway(r *gin.Engine) *gin.Engine {
	service := r.Group("/v1/account")
	{
		service.POST("/create", as.(accountService).proxy)
		service.GET("/auth", as.(accountService).proxy)
	}
	return r
}

func (s accountService) proxy(c *gin.Context) {
	endpoint := s.GetEndpoint(c, s.encodeRequest, s.decodeResponse)

	var request AuthRequest
	json.NewDecoder(c.Request.Body).Decode(&request)
	response, _ := endpoint(c, request)

	resp := response.(AuthResponse)
	w := c.Writer
	w.WriteHeader(resp.StatusCode)
	encoder := json.NewEncoder(w)
	encoder.Encode(resp)
}

func (s accountService) encodeRequest(ctx context.Context, request *http.Request, in interface{}) error {
	body, _ := json.Marshal(in)
	newRequest, _ := http.NewRequest(request.Method, "", bytes.NewBuffer(body))
	request.Body = newRequest.Body
	return nil
}

func (s accountService) decodeResponse(ctx context.Context, resp *http.Response) (interface{}, error) {
	var response AuthResponse
	json.NewDecoder(resp.Body).Decode(&response)
	return response, nil
}
