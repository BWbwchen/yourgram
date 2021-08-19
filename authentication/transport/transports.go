package authentication

import (
	"context"
	"encoding/json"
	"net/http"

	authsvc "yourgram/authentication/service"
)

func DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var request authsvc.AuthRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(authsvc.AuthResponse)
	w.WriteHeader(resp.StatusCode)
	encoder := json.NewEncoder(w)
	encoder.Encode(resp)
	return nil
}
