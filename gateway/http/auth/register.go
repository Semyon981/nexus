package auth

import (
	"github.com/Semyon981/nexus/proto/authpb"
	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.Engine, authclient authpb.ServiceClient) {
	h := NewHandler(authclient)

	authEndpoints := router.Group("/auth")
	{
		authEndpoints.POST("/sign-up", h.SignUp)
		authEndpoints.POST("/sign-in", h.SignIn)
	}

	router.GET("ping", func(c *gin.Context) { c.JSON(200, "pong") })

}
