package jwt_transport

import (
	"context"
	"encoding/json"
	"net/http"

	jwt_svc "yourgram/jwt/service"
)

func DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var request jwt_svc.AuthorizationRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(jwt_svc.AuthorizationResponse)
	w.WriteHeader(resp.StatusCode)
	encoder := json.NewEncoder(w)
	encoder.Encode(resp)
	return nil
}
