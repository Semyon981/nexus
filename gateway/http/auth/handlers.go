package auth

import (
	"net/http"

	"github.com/Semyon981/nexus/proto/authpb"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	authclient authpb.ServiceClient
}

func NewHandler(authclient authpb.ServiceClient) *Handler {
	return &Handler{
		authclient: authclient,
	}
}

func (h *Handler) SignUp(c *gin.Context) {
	inp := new(authpb.SignUpRequest)

	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	_, err := h.authclient.SignUp(c.Request.Context(), inp)
	if err != nil {
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) SignIn(c *gin.Context) {
	inp := new(authpb.SignInRequest)
	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	token, err := h.authclient.SignIn(c.Request.Context(), inp)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, token)
}
