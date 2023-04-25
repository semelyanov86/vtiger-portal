package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/internal/service"
	"net/http"
)

func (h *Handler) initTicketsRoutes(api *gin.RouterGroup) {
	tickets := api.Group("/tickets")
	{
		tickets.GET("/", h.getAllTickets)
		tickets.POST("/", h.createTicket)
		tickets.GET("/:id", h.getTicket)
		tickets.GET("/:id/comments", h.getComments)
		tickets.GET("/:id/documents", h.getDocuments)
		tickets.GET("/:id/file/:file", h.getFile)
	}
}

func (h *Handler) getTicket(c *gin.Context) {
	id := h.getAndValidateId(c, "id")

	userModel := h.getValidatedUser(c)
	if userModel == nil || id == "" {
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
	id := h.getAndValidateId(c, "id")

	userModel := h.getValidatedUser(c)

	if id == "" || userModel == nil {
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

func (h *Handler) getDocuments(c *gin.Context) {
	id := h.getAndValidateId(c, "id")

	userModel := h.getValidatedUser(c)

	if id == "" || userModel == nil {
		return
	}

	documents, err := h.services.HelpDesk.GetRelatedDocuments(c.Request.Context(), id, userModel.AccountId)
	if errors.Is(service.ErrOperationNotPermitted, err) {
		notPermittedResponse(c)
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, DataResponse[domain.Document]{
		Data:  documents,
		Count: len(documents),
		Page:  1,
		Size:  100,
	})
}

func (h *Handler) getFile(c *gin.Context) {
	id := h.getAndValidateId(c, "id")

	userModel := h.getValidatedUser(c)
	fileId := h.getAndValidateId(c, "file")

	if id == "" || userModel == nil || fileId == "" {
		return
	}

	file, err := h.services.Documents.GetFile(c.Request.Context(), fileId, id)

	if errors.Is(service.ErrOperationNotPermitted, err) {
		notPermittedResponse(c)
		return
	}

	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, file)
}

func (h *Handler) getAllTickets(c *gin.Context) {
	userModel := h.getValidatedUser(c)
	page, size := h.getPageAndSizeParams(c)

	if userModel == nil || page < 0 || size < 0 {
		return
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

func (h *Handler) createTicket(c *gin.Context) {
	var inp service.CreateTicketInput
	if err := c.BindJSON(&inp); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": fieldErr.Field(), "message": fieldErr.Error()})
			return // exit on first error
		}
	}
	userModel := h.getValidatedUser(c)
	if userModel == nil {
		return
	}

	ticket, err := h.services.HelpDesk.CreateTicket(c.Request.Context(), inp, *userModel)
	if errors.Is(service.ErrValidation, err) {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": "ticketcategories", "message": err.Error()})
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, ticket)
}
