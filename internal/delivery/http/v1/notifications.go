package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"net/http"
	"strconv"
)

func (h *Handler) initNotificationsRoutes(api *gin.RouterGroup) {
	invoices := api.Group("/notifications")
	{
		invoices.GET("", h.getAllNotifications)
		invoices.DELETE("/:id", h.markNotificationRead)
	}
}

func (h *Handler) getAllNotifications(c *gin.Context) {
	userModel := h.getValidatedUser(c)
	if userModel == nil {
		return
	}
	notifications, err := h.services.Notifications.GetNotificationsByUserId(c.Request.Context(), userModel.Crmid)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	res := DataResponse[domain.Notification]{
		Data:  notifications,
		Count: len(notifications),
		Page:  1,
		Size:  100,
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) markNotificationRead(c *gin.Context) {
	userModel := h.getValidatedUser(c)
	if userModel == nil {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil || id < 1 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid notification number"})
		return
	}
	err = h.services.Notifications.MarkNotificationRead(c.Request.Context(), int64(id), userModel.Crmid)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
