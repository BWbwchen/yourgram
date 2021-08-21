package main

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/hashicorp/consul/api"
)

type Service interface {
	GetEndpoint(c *gin.Context, enc httptransport.EncodeRequestFunc, dec httptransport.DecodeResponseFunc) endpoint.Endpoint
	encodeRequest(ctx context.Context, request *http.Request, in interface{}) error
	decodeResponse(ctx context.Context, resp *http.Response) (interface{}, error)
}

type GateWayStruct struct {
	serviceName string
	tags        []string
	logger      log.Logger
	client      consul.Client
}

func NewService() GateWayStruct {
	config := api.DefaultConfig()
	config.Address = os.Getenv("consul_url")
	api_client, _ := api.NewClient(config)
	ret := GateWayStruct{
		serviceName: "None",
		tags:        []string{},
		logger:      log.NewLogfmtLogger(os.Stdout),
		client:      consul.NewClient(api_client),
	}
	return ret
}

func (gw GateWayStruct) GetEndpoint(c *gin.Context, enc httptransport.EncodeRequestFunc, dec httptransport.DecodeResponseFunc) endpoint.Endpoint {
	// get service instance
	instancer := consul.NewInstancer(gw.client, gw.logger, gw.serviceName, gw.tags, true)

	relativePathSlice := strings.Split(c.Request.URL.Path, "/")
	relativePath := relativePathSlice[len(relativePathSlice)-1]

	factory := func(service_url string) (endpoint.Endpoint, io.Closer, error) {
		tart, _ := url.Parse("http://" + service_url + "/" + relativePath)
		return httptransport.NewClient(c.Request.Method, tart, enc, dec).Endpoint(), nil, nil
	}

	endpointer := sd.NewEndpointer(instancer, factory, gw.logger)
	endpoints, _ := endpointer.Endpoints()
	endpoint := endpoints[0]
	return endpoint
}
