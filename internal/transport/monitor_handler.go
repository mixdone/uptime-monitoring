package transport

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mixdone/uptime-monitoring/internal/models"
	"github.com/mixdone/uptime-monitoring/internal/models/dto"
)

// @Summary Create a new monitor
// @Tags monitors
// @Accept json
// @Produce json
// @Param monitor body dto.MonitorRequest true "monitor create request"
// @Success 200 {object} dto.MonitorResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /monitors [post]
func (h *Handler) createMonitor(c *gin.Context) {
	var req dto.MonitorRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	monitor := models.Monitor{
		UserID:           userID.(int64),
		Name:             req.Name,
		Type:             req.Type,
		Target:           req.Target,
		Timeout:          req.Timeout,
		Interval:         req.Interval,
		IsActive:         req.IsActive,
		RequestSpec:      req.RequestSpec,
		ExpectedResponse: req.ExpectedResponse,
	}

	id, err := h.services.Monitor.CreateMonitor(c.Request.Context(), monitor)

	if err != nil {
		h.logger.WithError(err).Error("Failed to create monitor")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create monitor"})
		return
	}

	c.JSON(http.StatusOK, dto.MonitorResponse{ID: id})
}

// @Summary Get monitor by ID
// @Tags monitors
// @Produce json
// @Param id path int true "Monitor ID"
// @Success 200 {object} models.Monitor
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /monitors/{id} [get]
func (h *Handler) getMonitor(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Monitor ID"})
		return
	}

	monitor, err := h.services.Monitor.GetMonitor(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "monitor not found"})
		return
	}

	c.JSON(http.StatusOK, monitor)
}

// @Summary Get all user's monitors
// @Tags monitors
// @Produce json
// @Success 200 {object} []models.Monitor
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /monitors [get]
func (h *Handler) getAllUserMonitor(c *gin.Context) {

	id, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	monitors, err := h.services.Monitor.GetAllUserMonitors(c.Request.Context(), id.(int64))
	if err != nil {
		h.logger.WithError(err).Error("Failed to fetch user monitors")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch all user monitors"})
		return
	}

	c.JSON(http.StatusOK, monitors)
}

// @Summary Update monitor
// @Tags monitors
// @Accept json
// @Produce json
// @Param id path int true "Monitor ID"
// @Param monitor body dto.MonitorRequest true "monitor update request"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /monitors/{id} [put]
func (h *Handler) updateMonitor(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid monitor ID"})
		return
	}

	var req dto.MonitorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	monitor := models.Monitor{
		ID:               id,
		UserID:           userID.(int64),
		Name:             req.Name,
		Type:             req.Type,
		Target:           req.Target,
		Timeout:          req.Timeout,
		Interval:         req.Interval,
		IsActive:         req.IsActive,
		RequestSpec:      req.RequestSpec,
		ExpectedResponse: req.ExpectedResponse,
	}

	if err := h.services.Monitor.UpdateMonitor(c.Request.Context(), monitor); err != nil {
		h.logger.WithError(err).Error("Failed to update monitor")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update monitor"})
		return
	}

	c.Status(http.StatusNoContent)

}

// @Summary Delete monitor
// @Tags monitors
// @Accept json
// @Produce json
// @Param id path int true "Monitor ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /monitors/{id} [delete]
func (h *Handler) deleteMonitor(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid monitor ID"})
		return
	}

	if err := h.services.Monitor.DeleteMonitor(c.Request.Context(), id); err != nil {
		h.logger.WithError(err).Error("Failed to delete monitor")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete monitor"})
		return
	}

	c.Status(http.StatusNoContent)
}
