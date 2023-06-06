package http

import (
	"net/http"

	"github.com/Semyon981/nexus/proto/userspb"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	c userspb.UserServiceClient
}

func NewHandler(c userspb.UserServiceClient) *Handler {
	return &Handler{
		c: c,
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

	if err := h.useCase.SignUp(c.Request.Context(), inp.Number, inp.Password, inp.Name, inp.Lastname); err != nil {
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

	token, err := h.useCase.SignIn(c.Request.Context(), inp.Number, inp.Password)
	if err != nil {
		if err == auth.ErrUserNotFound {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, signInResponse{Token: token})
}
