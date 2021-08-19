package main

import (
	auth "yourgram/authentication"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r = addAuthenticationService(r)

	r.Run()
}

func addAuthenticationService(r *gin.Engine) *gin.Engine {
	authentication := r.Group("/v1")
	{
		authentication.POST("/create", func(c *gin.Context) {
			auth.CreateAccountHandler().ServeHTTP(c.Writer, c.Request)
		})
		authentication.GET("/auth", func(c *gin.Context) {
			auth.UserLoginHandler().ServeHTTP(c.Writer, c.Request)
		})
	}
	return r
}
