package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"net/http"
)

func (h *Handler) initSalesOrdersRoutes(api *gin.RouterGroup) {
	invoices := api.Group("/sales-orders")
	{
		invoices.GET("/", h.getAllSalesOrders)
		invoices.GET("/:id", h.getSalesOrder)
	}
}

func (h *Handler) getSalesOrder(c *gin.Context) {
	id := h.getAndValidateId(c, "id")

	userModel := h.getValidatedUser(c)
	if userModel == nil || id == "" {
		return
	}

	salesOrder, err := h.services.SalesOrders.GetSalesOrderById(c.Request.Context(), id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if userModel.AccountId != salesOrder.AccountID {
		notPermittedResponse(c)
		return
	}
	res := AloneDataResponse[domain.SalesOrder]{
		Data: salesOrder,
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) getAllSalesOrders(c *gin.Context) {
	userModel := h.getValidatedUser(c)
	page, size := h.getPageAndSizeParams(c)

	if userModel == nil || page < 0 || size < 0 {
		return
	}
	sortString := c.DefaultQuery("sort", "-salesorder_no")

	if !isSortStringValid(sortString, []string{"salesorder_no", "id", "subject", "sostatus", "hdnGrandTotal"}) {
		newResponse(c, http.StatusUnprocessableEntity, "sort value "+sortString+" is not allowed")
		return
	}

	salesOrders, count, err := h.services.SalesOrders.GetAll(c.Request.Context(), repository.PaginationQueryFilter{
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
	c.JSON(http.StatusOK, DataResponse[domain.SalesOrder]{
		Data:  salesOrders,
		Count: count,
		Page:  page,
		Size:  size,
	})
}
