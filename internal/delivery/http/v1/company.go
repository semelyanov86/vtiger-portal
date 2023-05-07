package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"net/http"
	"strings"
)

func (h *Handler) initCompanyRoutes(api *gin.RouterGroup) {
	users := api.Group("/company")
	{
		users.GET("/", h.getCompany)
	}
}

func (h *Handler) getCompany(c *gin.Context) {
	id := h.config.Vtiger.Business.CompanyId
	if id == "" {
		newResponse(c, http.StatusBadRequest, "there is no company code in config")

		return
	}

	if !strings.Contains(id, "x") {
		newResponse(c, http.StatusUnprocessableEntity, "wrong company id")

		return
	}

	company, err := h.services.Company.GetCompany(c.Request.Context())
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	resp := AloneDataResponse[domain.Company]{
		Data: company,
	}
	c.JSON(http.StatusOK, resp)
}
