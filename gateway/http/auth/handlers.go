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

type signUpInput struct {
	Number   string `json:"number" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Lastname string `json:"lastname" binding:"required"`
}

func (h *Handler) SignUp(c *gin.Context) {
	inp := new(signUpInput)
	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	_, err := h.authclient.SignUp(c.Request.Context(), &authpb.SignUpRequest{Number: inp.Number, Password: inp.Password, Name: inp.Name, Lastname: inp.Lastname})
	if err != nil {
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	c.Status(http.StatusOK)
}

type signInResponse struct {
	Token string `json:"token"`
}

type signInInput struct {
	Number   string `json:"number" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) SignIn(c *gin.Context) {
	inp := new(signInInput)
	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	token, err := h.authclient.SignIn(c.Request.Context(), &authpb.SignInRequest{Number: inp.Number, Password: inp.Password})
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, signInResponse{Token: token.Token})
}
