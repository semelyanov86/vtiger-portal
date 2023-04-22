package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (h *Handler) initTicketsRoutes(api *gin.RouterGroup) {
	users := api.Group("/tickets")
	{
		users.GET("/:id", h.getTicket)
	}
}

func (h *Handler) getTicket(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		newResponse(c, http.StatusBadRequest, "code is empty")

		return
	}

	if !strings.Contains(id, "x") {
		newResponse(c, http.StatusUnprocessableEntity, "wrong id")

		return
	}

	userModel := h.services.Context.ContextGetUser(c)
	if userModel.Crmid == "" || userModel.Id < 1 {
		anonymousResponse(c)
		return
	}
	ticket, err := h.services.HelpDesk.GetHelpDeskById(c.Request.Context(), id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if userModel.AccountId != ticket.ParentID {
		notPermittedResponse(c)
		return
	}
	c.JSON(http.StatusOK, ticket)
}
