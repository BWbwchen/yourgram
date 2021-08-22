package jwt_svc

type AuthorizationRequest struct {
	UserData string `json:"userdata"`
	JWTToken string `json:"jwt"`
}

type AuthorizationResponse struct {
	StatusCode int    `json:"status"`
	Return     string `json:"return"`
}
