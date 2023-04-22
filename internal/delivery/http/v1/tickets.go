package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/internal/service"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) initTicketsRoutes(api *gin.RouterGroup) {
	tickets := api.Group("/tickets")
	{
		tickets.GET("/", h.getAllTickets)
		tickets.GET("/:id", h.getTicket)
		tickets.GET("/:id/comments", h.getComments)
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

func (h *Handler) getComments(c *gin.Context) {
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

	comments, err := h.services.HelpDesk.GetRelatedComments(c.Request.Context(), id, userModel.AccountId)
	if errors.Is(service.ErrOperationNotPermitted, err) {
		notPermittedResponse(c)
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, DataResponse[domain.Comment]{
		Data:  comments,
		Count: len(comments),
		Page:  1,
		Size:  100,
	})
}

func (h *Handler) getAllTickets(c *gin.Context) {
	userModel := h.services.Context.ContextGetUser(c)
	if userModel.Crmid == "" || userModel.Id < 1 {
		anonymousResponse(c)
		return
	}
	if userModel.AccountId == "" {
		notPermittedResponse(c)
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid page number"})
		return
	}

	size, err := strconv.Atoi(c.DefaultQuery("size", strconv.Itoa(h.config.Vtiger.Business.DefaultPagination)))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid page size"})
		return
	}
	if size < 1 {
		size = 100
	}

	tickets, count, err := h.services.HelpDesk.GetAll(c.Request.Context(), repository.TicketsQueryFilter{
		Page:     page,
		PageSize: size,
		Client:   userModel.AccountId,
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, DataResponse[domain.HelpDesk]{
		Data:  tickets,
		Count: count,
		Page:  page,
		Size:  size,
	})
}
