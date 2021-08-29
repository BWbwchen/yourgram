package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"yourgram/upload/health"
	"yourgram/upload/pb"

	upload_endp "yourgram/upload/endpoint"
	upload_svc "yourgram/upload/service"
	upload_trans "yourgram/upload/transport"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/sd/consul"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	check.GRPC = os.Getenv("localIP") + ":" + strconv.Itoa(reg.Port) + "/Health/Check"

	reg.Check = &check

	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}
	consulClient := consul.NewClient(client)

	return consul.NewRegistrar(consulClient, &reg, log.NewLogfmtLogger(os.Stdout))
}

func main() {
	var logger log.Logger
	logger = log.NewJSONLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	listener, _ := net.Listen("tcp", ":"+os.Getenv("PORT"))

	addservice := upload_svc.NewService(logger)
	addendpoint := upload_endp.MakeEndpoints(addservice)
	grpcServer := upload_trans.NewGRPCServer(addendpoint, logger)

	healthService := health.NewService(logger)

	registar := registerService()

	errc := make(chan error)
	go func() {
		baseServer := grpc.NewServer()
		pb.RegisterUploadServiceServer(baseServer, grpcServer)
		pb.RegisterHealthServer(baseServer, healthService)
		level.Info(logger).Log("msg", "Server started successfully 🚀")

		registar.Register()
		reflection.Register(baseServer)
		baseServer.Serve(listener)

		registar.Register()
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGALRM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	level.Error(logger).Log("exit", <-errc)
	registar.Deregister()

}
