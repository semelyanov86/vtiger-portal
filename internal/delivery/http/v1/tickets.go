package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/service"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
	"net/http"
)

func (h *Handler) initTicketsRoutes(api *gin.RouterGroup) {
	tickets := api.Group("/tickets")
	{
		tickets.GET("/", h.getAllTickets)
		tickets.POST("/", h.createTicket)
		tickets.GET("/:id", h.getTicket)
		tickets.PUT("/:id", h.updateTicket)
		tickets.PATCH("/:id", h.updatePartlyTicket)
		tickets.GET("/:id/comments", h.getComments)
		tickets.POST("/:id/comments", h.addComment)
		tickets.GET("/:id/documents", h.getDocuments)
		tickets.POST("/:id/documents", h.uploadTicketDocuments)
		tickets.GET("/:id/file/:file", h.getFile)
	}
}

func (h *Handler) getTicket(c *gin.Context) {
	id := h.getAndValidateId(c, "id")

	userModel := h.getValidatedUser(c)
	if userModel == nil || id == "" {
		return
	}

	ticket, err := h.services.HelpDesk.GetHelpDeskById(c.Request.Context(), id, *userModel)
	if errors.Is(service.ErrOperationNotPermitted, err) {
		notPermittedResponse(c)
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	res := AloneDataResponse[domain.HelpDesk]{
		Data: ticket,
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) getComments(c *gin.Context) {
	id := h.getAndValidateId(c, "id")

	userModel := h.getValidatedUser(c)

	if id == "" || userModel == nil {
		return
	}

	comments, err := h.services.HelpDesk.GetRelatedComments(c.Request.Context(), id, *userModel)
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

type createCommentInput struct {
	Commentcontent string `json:"commentcontent" binding:"required"`
}

func (h *Handler) addComment(c *gin.Context) {
	id := h.getAndValidateId(c, "id")

	userModel := h.getValidatedUser(c)

	if id == "" || userModel == nil {
		return
	}
	var inp createCommentInput
	if err := c.BindJSON(&inp); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": fieldErr.Field(), "message": fieldErr.Error()})
			return // exit on first error
		}
	}
	comment, err := h.services.HelpDesk.AddComment(c.Request.Context(), inp.Commentcontent, id, *userModel)
	if errors.Is(service.ErrOperationNotPermitted, err) {
		notPermittedResponse(c)
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, comment)
}

func (h *Handler) getDocuments(c *gin.Context) {
	id := h.getAndValidateId(c, "id")

	userModel := h.getValidatedUser(c)

	if id == "" || userModel == nil {
		return
	}

	documents, err := h.services.HelpDesk.GetRelatedDocuments(c.Request.Context(), id, *userModel)
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

func (h *Handler) uploadTicketDocuments(c *gin.Context) {
	id := h.getAndValidateId(c, "id")
	userModel := h.getValidatedUser(c)

	if id == "" || userModel == nil {
		notPermittedResponse(c)
		return
	}
	_, err := h.services.HelpDesk.GetHelpDeskById(c.Request.Context(), id, *userModel)
	if errors.Is(service.ErrOperationNotPermitted, err) {
		notPermittedResponse(c)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	document, err := h.services.Documents.AttachFile(c.Request.Context(), file, id, *userModel, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, AloneDataResponse[domain.Document]{Data: document})
}

func (h *Handler) getFile(c *gin.Context) {
	id := h.getAndValidateId(c, "id")

	userModel := h.getValidatedUser(c)
	fileId := h.getAndValidateId(c, "file")

	if id == "" || userModel == nil || fileId == "" {
		return
	}

	_, err := h.services.HelpDesk.GetHelpDeskById(c.Request.Context(), id, *userModel)
	if errors.Is(service.ErrOperationNotPermitted, err) {
		notPermittedResponse(c)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	res := AloneDataResponse[vtiger.File]{
		Data: file,
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) getAllTickets(c *gin.Context) {
	userModel := h.getValidatedUser(c)
	page, size := h.getPageAndSizeParams(c)

	if userModel == nil || page < 0 || size < 0 {
		newResponse(c, http.StatusBadRequest, "Wrong pagination params or auth user")
		return
	}
	sortString := c.DefaultQuery("sort", "-ticket_no")

	if !isSortStringValid(sortString, []string{"ticket_no", "ticket_title", "ticketstatus", "hours", "days", "ticketcategories"}) {
		newResponse(c, http.StatusUnprocessableEntity, "sort value "+sortString+" is not allowed")
		return
	}

	tickets, count, err := h.services.HelpDesk.GetAll(c.Request.Context(), vtiger.PaginationQueryFilter{
		Page:     page,
		PageSize: size,
		Client:   userModel.AccountId,
		Contact:  userModel.Crmid,
		Sort:     sortString,
		Search:   c.DefaultQuery("search", ""),
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

func (h *Handler) updateTicket(c *gin.Context) {
	var inp service.CreateTicketInput
	if err := c.BindJSON(&inp); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": fieldErr.Field(), "message": fieldErr.Error()})
			return // exit on first error
		}
	}
	id := h.getAndValidateId(c, "id")
	userModel := h.getValidatedUser(c)
	if userModel == nil || id == "" {
		return
	}

	ticket, err := h.services.HelpDesk.UpdateTicket(c.Request.Context(), inp, id, *userModel)
	if errors.Is(service.ErrValidation, err) {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": "ticketcategories", "message": err.Error()})
		return
	}
	if errors.Is(service.ErrOperationNotPermitted, err) {
		notPermittedResponse(c)
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusAccepted, ticket)
}

func (h *Handler) updatePartlyTicket(c *gin.Context) {
	var inp map[string]any
	if err := c.ShouldBindJSON(&inp); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": "ticketstatus", "message": "Incorrect value"})
		return // exit on first error
	}
	id := h.getAndValidateId(c, "id")
	userModel := h.getValidatedUser(c)
	if userModel == nil || id == "" {
		return
	}

	ticket, err := h.services.HelpDesk.Revise(c.Request.Context(), inp, id, *userModel)
	if errors.Is(service.ErrValidation, err) {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": "ticketstatus", "message": err.Error()})
		return
	}
	if errors.Is(service.ErrOperationNotPermitted, err) {
		notPermittedResponse(c)
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusAccepted, ticket)
}
