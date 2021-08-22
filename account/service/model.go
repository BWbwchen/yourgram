package account_svc

type AuthRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type AuthResponse struct {
	StatusCode int    `json:"status"`
	JWTToken   string `json:"jwt"`
}

type UserInfo AuthRequest
