package main

import (
	"net/http"
	"os"
	"strconv"

	account_endp "yourgram/account/endpoint"
	account_svc "yourgram/account/service"
	account_trans "yourgram/account/transport"

	"github.com/gin-gonic/gin"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
)

func registerService() {
	// register the service to consul
	config := api.DefaultConfig()
	config.Address = os.Getenv("consul_url")

	reg := api.AgentServiceRegistration{}
	reg.Name = "account_service"
	reg.ID = reg.Name + uuid.NewString()
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

	client.Agent().ServiceRegister(&reg)
}

func CreateAccountHandler() *httptransport.Server {
	svc := account_svc.AuthenticateWorker{}
	return httptransport.NewServer(
		account_endp.MakeCreateAccountEndPoint(svc),
		account_trans.DecodeRequest,
		account_trans.EncodeResponse,
	)
}

func UserLoginHandler() *httptransport.Server {
	svc := account_svc.AuthenticateWorker{}
	return httptransport.NewServer(
		account_endp.MakeUserLoginEndPoint(svc),
		account_trans.DecodeRequest,
		account_trans.EncodeResponse,
	)
}

func main() {
	r := gin.Default()

	r.POST("/create", func(c *gin.Context) {
		CreateAccountHandler().ServeHTTP(c.Writer, c.Request)
	})
	r.GET("/auth", func(c *gin.Context) {
		UserLoginHandler().ServeHTTP(c.Writer, c.Request)
	})
	r.GET("/health", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	registerService()
	r.Run()
}
