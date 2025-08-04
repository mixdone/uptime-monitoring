package transport

import (
	_ "github.com/mixdone/uptime-monitoring/docs"

	"github.com/gin-gonic/gin"
	"github.com/mixdone/uptime-monitoring/internal/services"
	"github.com/mixdone/uptime-monitoring/pkg/logger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *services.Services
	logger   logger.Logger
}

func NewHandler(services *services.Services, log logger.Logger) *Handler {
	return &Handler{
		services: services,
		logger:   log,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")
	{
		auth.POST("/register", h.register)
		auth.POST("/login", h.login)
		auth.POST("/refresh", h.refresh)
		auth.POST("/logout", h.authMiddleware, h.logout)
	}

	return router
}
