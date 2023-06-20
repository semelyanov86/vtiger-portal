package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
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
	res := AloneDataResponse[domain.Invoice]{
		Data: invoice,
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) getAllInvoices(c *gin.Context) {
	userModel := h.getValidatedUser(c)
	page, size := h.getPageAndSizeParams(c)

	if userModel == nil || page < 0 || size < 0 {
		return
	}
	sortString := c.DefaultQuery("sort", "-invoice_no")

	if !isSortStringValid(sortString, []string{"invoice_no", "id", "subject", "invoicestatus", "invoicedate", "hdnGrandTotal"}) {
		newResponse(c, http.StatusUnprocessableEntity, "sort value "+sortString+" is not allowed")
		return
	}

	invoices, count, err := h.services.Invoices.GetAll(c.Request.Context(), vtiger.PaginationQueryFilter{
		Page:     page,
		PageSize: size,
		Client:   userModel.AccountId,
		Sort:     sortString,
		Search:   c.DefaultQuery("search", ""),
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
