package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initInvoicesRoutes(api *gin.RouterGroup) {
	invoices := api.Group("/invoices")
	{
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
