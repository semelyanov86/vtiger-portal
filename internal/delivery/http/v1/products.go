package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"net/http"
	"strconv"
)

func (h *Handler) initProductsRoutes(api *gin.RouterGroup) {
	products := api.Group("/products")
	{
		products.GET("/", h.getAllProducts)
		products.GET("/:id", h.getProduct)
	}
}

func (h *Handler) getProduct(c *gin.Context) {
	id := h.getAndValidateId(c, "id")

	userModel := h.getValidatedUser(c)
	if userModel == nil || id == "" {
		return
	}

	product, err := h.services.Products.GetProductById(c.Request.Context(), id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	res := AloneDataResponse[domain.Product]{
		Data: product,
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) getAllProducts(c *gin.Context) {
	userModel := h.getValidatedUser(c)
	page, size := h.getPageAndSizeParams(c)

	if userModel == nil || page < 0 || size < 0 {
		newResponse(c, http.StatusBadRequest, "Wrong token or page size")
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

	products, count, err := h.services.Products.GetAll(c.Request.Context(), repository.PaginationQueryFilter{
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
	c.JSON(http.StatusOK, DataResponse[domain.Product]{
		Data:  products,
		Count: count,
		Page:  page,
		Size:  size,
	})
}
