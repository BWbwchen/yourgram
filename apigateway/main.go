package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(CORS())

	r = AccountGateway(r)
	r = UploadGateway(r)
	r = ViewImageGateway(r)
	//	r = JWTGateway(r)

	r.Run()
}
