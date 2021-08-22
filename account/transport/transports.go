package account_transport

import (
	"context"
	"encoding/json"
	"net/http"
	account_svc "yourgram/account/service"
)

func DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var request account_svc.AuthRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(account_svc.AuthResponse)
	w.WriteHeader(resp.StatusCode)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(resp)
	return nil
}
