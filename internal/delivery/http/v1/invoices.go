package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"net/http"
)

func (h *Handler) initInvoicesRoutes(api *gin.RouterGroup) {
	invoices := api.Group("/invoices")
	{
		invoices.GET("/", h.getAllInvoices)
		invoices.GET("/:id", h.getInvoice)
	}
}

func (h *Handler) getInvoice(c *gin.Context) {
	id := h.getAndValidateId(c, "id")

	userModel := h.getValidatedUser(c)
	if userModel == nil || id == "" {
		return
	}

	invoice, err := h.services.Invoices.GetInvoiceById(c.Request.Context(), id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if userModel.AccountId != invoice.AccountID {
		notPermittedResponse(c)
		return
	}
	c.JSON(http.StatusOK, invoice)
}

func (h *Handler) getAllInvoices(c *gin.Context) {
	userModel := h.getValidatedUser(c)
	page, size := h.getPageAndSizeParams(c)

	if userModel == nil || page < 0 || size < 0 {
		return
	}

	invoices, count, err := h.services.Invoices.GetAll(c.Request.Context(), repository.PaginationQueryFilter{
		Page:     page,
		PageSize: size,
		Client:   userModel.AccountId,
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, DataResponse[domain.Invoice]{
		Data:  invoices,
		Count: count,
		Page:  page,
		Size:  size,
	})
}
