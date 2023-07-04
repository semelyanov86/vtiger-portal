package http

import (
	_ "embed"
	"github.com/flowchartsman/swaggerui"
	"github.com/gin-gonic/gin"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	v1 "github.com/semelyanov86/vtiger-portal/internal/delivery/http/v1"
	"github.com/semelyanov86/vtiger-portal/internal/service"
	"github.com/semelyanov86/vtiger-portal/pkg/limiter"
	"net/http"
)

//go:embed swagger.yaml
var spec []byte

type Handler struct {
	services *service.Services
	config   *config.Config
}

func NewHandler(services *service.Services, config *config.Config) *Handler {
	return &Handler{services: services, config: config}
}

func (h *Handler) Init() *gin.Engine {
	// Init gin handler
	router := gin.Default()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
		limiter.Limit(int(h.config.Limiter.Rps), h.config.Limiter.Burst, h.config.Limiter.TTL),
		h.corsMiddleware,
		h.authenticate,
	)

	router.GET("/swagger/*w", gin.WrapH(http.StripPrefix("/swagger", swaggerui.Handler(spec))))

	// Init router
	router.GET("/api/v1/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "active"})
	})

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.services, h.config)
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
