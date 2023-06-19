package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/service"
	"net/http"
	"strconv"
	"strings"
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
		h.initModulesRoutes(v1)
		h.initCompanyRoutes(v1)
		h.initTicketsRoutes(v1)
		h.initFaqsRoutes(v1)
		h.initInvoicesRoutes(v1)
		h.initSalesOrdersRoutes(v1)
		h.initServiceContractsRoutes(v1)
		h.initProductsRoutes(v1)
		h.initServicesRoutes(v1)
		h.initProjectsRoutes(v1)
		h.initStatisticsRoutes(v1)
		h.initLeadsRoutes(v1)
		h.initOtpRoutes(v1)
		h.initSearchRoutes(v1)
		h.initPaymentRoutes(v1)
		h.initNotificationsRoutes(v1)
		h.initCustomModulesRoutes(v1)
	}
}

func (h *Handler) getAndValidateId(c *gin.Context, field string) string {
	if field == "" {
		field = "id"
	}
	id := c.Param(field)
	if id == "" {
		newResponse(c, http.StatusBadRequest, "code is empty")

		return ""
	}

	if !strings.Contains(id, "x") {
		newResponse(c, http.StatusUnprocessableEntity, "wrong id")

		return ""
	}
	return id
}

func (h *Handler) getValidatedUser(c *gin.Context) *domain.User {
	userModel := h.services.Context.ContextGetUser(c)
	if userModel.Crmid == "" || userModel.Id < 1 {
		anonymousResponse(c)
		return nil
	}
	if userModel.AccountId == "" {
		notPermittedResponse(c)
		return nil
	}
	return userModel
}

func (h *Handler) getPageAndSizeParams(c *gin.Context) (int, int) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid page number"})
		return -1, -1
	}

	size, err := strconv.Atoi(c.DefaultQuery("size", strconv.Itoa(h.config.Vtiger.Business.DefaultPagination)))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid page size"})
		return -1, -1
	}
	if size < 1 {
		size = 100
	}
	return page, size
}

func isSortStringValid(sort string, allowedFields []string) bool {
	fields := strings.Split(sort, ",")

	for _, field := range fields {
		fieldName := strings.TrimPrefix(field, "-")
		if !contains(allowedFields, fieldName) {
			return false
		}
	}

	return true
}

func contains(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}
