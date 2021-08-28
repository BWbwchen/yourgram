package account_svc

type Input struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Output struct {
	StatusCode int    `json:"status"`
	JWTToken   string `json:"jwt"`
}

type UserInfo Input
