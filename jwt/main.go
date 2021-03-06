package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	jwt_endp "yourgram/jwt/endpoint"
	jwt_svc "yourgram/jwt/service"
	jwt_trans "yourgram/jwt/transport"

	"github.com/go-kit/kit/log"

	"github.com/gin-gonic/gin"
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
	reg.Name = "jwt_service"
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

func CreateJWTHandler() *httptransport.Server {
	svc := jwt_svc.AuthorizationWorker{}
	return httptransport.NewServer(
		jwt_endp.MakeCreateJWTEndPoint(svc),
		jwt_trans.DecodeRequest,
		jwt_trans.EncodeResponse,
	)
}

func VerifyJWTHandler() *httptransport.Server {
	svc := jwt_svc.AuthorizationWorker{}
	return httptransport.NewServer(
		jwt_endp.MakeVerifyJWTEndPoint(svc),
		jwt_trans.DecodeRequest,
		jwt_trans.EncodeResponse,
	)
}

func main() {
	r := gin.Default()

	r.POST("/create", func(c *gin.Context) {
		CreateJWTHandler().ServeHTTP(c.Writer, c.Request)
	})
	r.GET("/auth", func(c *gin.Context) {
		VerifyJWTHandler().ServeHTTP(c.Writer, c.Request)
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
