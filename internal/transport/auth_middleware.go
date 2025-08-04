package transport

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mixdone/uptime-monitoring/internal/models/errs"
)

func (h *Handler) authMiddleware(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header missing"})
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid authorization format"})
		return
	}

	accessToken := parts[1]
	userID, err := h.services.Token.ValidateAccess(accessToken)
	if err != nil {
		if errors.Is(err, errs.ErrTokenExpired) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "expired token"})
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		}
		return
	}

	c.Set("userID", userID)

	c.Next()

}
