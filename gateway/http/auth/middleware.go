package auth

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Semyon981/nexus/proto/identifierpb"
	"github.com/gin-gonic/gin"
)

const CtxUserKey = "id_users"

type AuthMiddleware struct {
	identclient identifierpb.ServiceClient
}

func NewAuthMiddleware(identclient identifierpb.ServiceClient) gin.HandlerFunc {
	return (&AuthMiddleware{
		identclient: identclient,
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

	resp, err := m.identclient.JwtParse(c.Request.Context(), &identifierpb.JwtParseRequest{Token: headerParts[1]})
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	id, err := strconv.ParseInt(resp.Subject, 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Set(CtxUserKey, id)
}
