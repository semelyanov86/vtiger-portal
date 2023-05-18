package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initProductsRoutes(api *gin.RouterGroup) {
	products := api.Group("/products")
	{
		products.GET("/:id", h.getProduct)
	}
}

func (h *Handler) getProduct(c *gin.Context) {
	id := h.getAndValidateId(c, "id")

	userModel := h.getValidatedUser(c)
	if userModel == nil || id == "" {
		return
	}

	ticket, err := h.services.Products.GetProductById(c.Request.Context(), id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, ticket)
}
