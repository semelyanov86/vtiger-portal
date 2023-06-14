package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"net/http"
)

func (h *Handler) initSearchRoutes(api *gin.RouterGroup) {
	tickets := api.Group("/search")
	{
		tickets.GET("", h.globalSearch)
	}
}

func (h *Handler) globalSearch(c *gin.Context) {
	query := c.DefaultQuery("search", "")

	if len(query) < 3 {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": "query", "message": "Search query should be more then 3 symbols"})
		return
	}

	userModel := h.getValidatedUser(c)
	if userModel == nil {
		return
	}

	searches, err := h.services.Searches.GlobalSearch(c.Request.Context(), query, *userModel)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, DataResponse[domain.Search]{
		Data:  searches,
		Count: len(searches),
		Page:  1,
		Size:  100,
	})
}
