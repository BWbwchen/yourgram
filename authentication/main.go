package authentication

import (
	authEndpoint "yourgram/authentication/endpoint"
	authSvc "yourgram/authentication/service"
	authTransport "yourgram/authentication/transport"

	httptransport "github.com/go-kit/kit/transport/http"
)

func CreateAccountHandler() *httptransport.Server {
	svc := authSvc.AuthenticateWorker{}
	return httptransport.NewServer(
		authEndpoint.MakeCreateAccountEndPoint(svc),
		authTransport.DecodeRequest,
		authTransport.EncodeResponse,
	)
}

func UserLoginHandler() *httptransport.Server {
	svc := authSvc.AuthenticateWorker{}
	return httptransport.NewServer(
		authEndpoint.MakeUserLoginEndPoint(svc),
		authTransport.DecodeRequest,
		authTransport.EncodeResponse,
	)
}
