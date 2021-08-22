package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	upload_endp "yourgram/upload/endpoint"
	upload_svc "yourgram/upload/service"
	upload_trans "yourgram/upload/transport"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/consul"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
)

func registerService() *consul.Registrar {
	// register the service to consul
	config := api.DefaultConfig()
	config.Address = os.Getenv("consul_url")

	reg := api.AgentServiceRegistration{}
	reg.Name = "upload_service"
	reg.ID = reg.Name + uuid.New().String()
	reg.Address = os.Getenv("localIP")
	reg.Port, _ = strconv.Atoi(os.Getenv("PORT"))
	reg.Tags = []string{"primary"}

	check := api.AgentServiceCheck{}
	check.Interval = "9s"
	check.HTTP = "http://" + os.Getenv("localIP") + ":" + strconv.Itoa(reg.Port) + "/health"

	reg.Check = &check

	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}
	consulClient := consul.NewClient(client)

	return consul.NewRegistrar(consulClient, &reg, log.NewLogfmtLogger(os.Stdout))
}

func UploadHandler() *httptransport.Server {
	svc := upload_svc.UploadWorker{}
	return httptransport.NewServer(
		upload_endp.MakeUploadEndPoint(svc),
		upload_trans.DecodeRequest,
		upload_trans.EncodeResponse,
	)
}

func InfoHandler() *httptransport.Server {
	svc := upload_svc.UploadWorker{}
	return httptransport.NewServer(
		upload_endp.MakeInfoEndPoint(svc),
		upload_trans.DecodeRequest,
		upload_trans.EncodeResponse,
	)
}

func main() {
	r := gin.Default()

	r.POST("/upload", func(c *gin.Context) {
		UploadHandler().ServeHTTP(c.Writer, c.Request)
	})
	r.GET("/info", func(c *gin.Context) {
		InfoHandler().ServeHTTP(c.Writer, c.Request)
	})
	r.GET("/health", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	registar := registerService()
	registar.Register()

	errc := make(chan error)
	go func() {
		registar.Register()
		errc <- r.Run()
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	<-errc
	registar.Deregister()
}
