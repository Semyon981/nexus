package http

import (
	"github.com/Semyon981/nexus/gateway/auth"
	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.Engine, uc auth.UseCase) {
	h := NewHandler(uc)

	authEndpoints := router.Group("/auth")
	{
		authEndpoints.POST("/sign-up", h.SignUp)
		authEndpoints.POST("/sign-in", h.SignIn)
	}

	router.GET("ping", func(c *gin.Context) { c.JSON(200, "pong") })

}
