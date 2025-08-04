package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mixdone/uptime-monitoring/internal/models/dto"
)

// @Summary Register user
// @Tags auth
// @Accept json
// @Produce json
// @Param input body dto.RegisterRequest true "user credentials"
// @Success 200 {object} dto.AuthResult
// @Failure 400 {object} map[string]string
// @Router /auth/register [post]
func (h *Handler) register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.services.Auth.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Tags auth
// @Accept json
// @Produce json
// @Param input body dto.LoginRequest true "user credentials"
// @Success 200 {object} dto.AuthResult
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (h *Handler) login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.services.Auth.Login(c.Request.Context(), req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Summary Logout user
// @Security ApiKeyAuth
// @Tags auth
// @Accept json
// @Produce json
// @Param input body dto.LogoutRequest true "tokens"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/logout [post]
func (h *Handler) logout(c *gin.Context) {
	var req dto.LogoutRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err := h.services.Auth.Logout(c.Request.Context(), userID.(int64), req)
	if err != nil {
		h.logger.Error("Logout error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "logout failed"})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Refresh user
// @Tags auth
// @Accept json
// @Produce json
// @Param input body dto.RefreshRequest true "tokens"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/refresh [post]
func (h *Handler) refresh(c *gin.Context) {
	var req dto.RefreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.services.Token.ValidateRefresh(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "expired or invalid token"})
		return
	}

	authResult, err := h.services.Auth.RefreshTokens(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, authResult)
}
