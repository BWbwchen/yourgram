package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r = AccountGateway(r)
	r = UploadGateway(r)
	//	r = JWTGateway(r)

	r.Run()
}
