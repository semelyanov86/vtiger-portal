package http

import (
	"github.com/gin-gonic/gin"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	v1 "github.com/semelyanov86/vtiger-portal/internal/delivery/http/v1"
	"github.com/semelyanov86/vtiger-portal/internal/service"
	"github.com/semelyanov86/vtiger-portal/pkg/limiter"
	"net/http"
)

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

	/*	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
		if cfg.Environment != config.EnvLocal {
			docs.SwaggerInfo.Host = cfg.HTTP.Host
		}*/

	/*if h.config.Environment != "production" {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}*/

	// Init router
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.services)
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
