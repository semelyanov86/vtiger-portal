package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"net/http"
)

func (h *Handler) initFaqsRoutes(api *gin.RouterGroup) {
	tickets := api.Group("/faqs")
	{
		tickets.GET("/", h.getAllFaqs)
	}
}

func (h *Handler) getAllFaqs(c *gin.Context) {
	userModel := h.getValidatedUser(c)
	page, size := h.getPageAndSizeParams(c)

	if userModel == nil || page < 0 || size < 0 {
		return
	}

	faqs, count, err := h.services.Faqs.GetAll(c.Request.Context(), repository.PaginationQueryFilter{
		Page:     page,
		PageSize: size,
		Client:   userModel.AccountId,
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, DataResponse[domain.Faq]{
		Data:  faqs,
		Count: count,
		Page:  page,
		Size:  size,
	})
}
