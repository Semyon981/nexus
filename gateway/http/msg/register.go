package msg

import (
	"github.com/Semyon981/nexus/proto/msgpb"
	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.RouterGroup, msgclient msgpb.ServiceClient) {
	h := NewHandler(msgclient)

	authEndpoints := router.Group("/msg")
	{
		authEndpoints.POST("/send", h.SendMessage)
		authEndpoints.GET("/get", h.GetMessages)
	}

}
