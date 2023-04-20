package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/service"
)

type Handler struct {
	services *service.Services
	config   *config.Config
}

func NewHandler(services *service.Services, config *config.Config) *Handler {
	return &Handler{
		services: services,
		config:   config,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initUsersRoutes(v1)
		h.initManagersRoutes(v1)
	}
}
