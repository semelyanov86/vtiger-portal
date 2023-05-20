package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"net/http"
	"strconv"
)

func (h *Handler) initServicesRoutes(api *gin.RouterGroup) {
	products := api.Group("/services")
	{
		products.GET("/", h.getAllServices)
		products.GET("/:id", h.getService)
	}
}

func (h *Handler) getService(c *gin.Context) {
	id := h.getAndValidateId(c, "id")

	userModel := h.getValidatedUser(c)
	if userModel == nil || id == "" {
		newResponse(c, http.StatusBadRequest, "Wrong token or page size")
		return
	}

	service, err := h.services.Services.GetServiceById(c.Request.Context(), id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	res := AloneDataResponse[domain.Service]{
		Data: service,
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) getAllServices(c *gin.Context) {
	userModel := h.getValidatedUser(c)
	page, size := h.getPageAndSizeParams(c)

	if userModel == nil || page < 0 || size < 0 {
		return
	}

	discontinuedFilter := c.Query("filter[discontinued]")
	discontinued := true
	if discontinuedFilter != "" {
		tmpDiscontinued, err := strconv.ParseBool(discontinuedFilter)
		if err != nil {
			discontinued = true
			newResponse(c, http.StatusBadRequest, "Wrong filter value for discontinued, expected boolean")
			return
		}
		discontinued = tmpDiscontinued
	}

	services, count, err := h.services.Services.GetAll(c.Request.Context(), repository.PaginationQueryFilter{
		Page:     page,
		PageSize: size,
		Client:   userModel.AccountId,
		Contact:  userModel.Crmid,
		Filters: map[string]any{
			"discontinued": discontinued,
		},
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, DataResponse[domain.Service]{
		Data:  services,
		Count: count,
		Page:  page,
		Size:  size,
	})
}
