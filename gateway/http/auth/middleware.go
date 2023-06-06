package http

import (
	"net/http"
	"strings"

	"github.com/Semyon981/nexus/gateway/auth"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	usecase auth.UseCase
}

func NewAuthMiddleware(usecase auth.UseCase) gin.HandlerFunc {
	return (&AuthMiddleware{
		usecase: usecase,
	}).Handle
}

func (m *AuthMiddleware) Handle(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if headerParts[0] != "Bearer" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	id_users, err := m.usecase.ParseToken(c.Request.Context(), headerParts[1])
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set(auth.CtxUserKey, id_users)
}
