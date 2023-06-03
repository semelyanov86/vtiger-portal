package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/service"
	"net/http"
)

func (h *Handler) initLeadsRoutes(api *gin.RouterGroup) {
	leads := api.Group("/leads")
	{
		leads.POST("/", h.createLead)
	}
}

func (h *Handler) createLead(c *gin.Context) {
	userModel := h.getValidatedUser(c)

	if userModel == nil {
		return
	}
	var inp domain.Lead
	if err := c.BindJSON(&inp); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": fieldErr.Field(), "message": fieldErr.Error()})
			return // exit on first error
		}
	}
	inp.Description = "From Contact " + userModel.Crmid + " - " + userModel.FirstName + " " + userModel.LastName
	lead, err := h.services.Leads.Create(c.Request.Context(), inp)
	if errors.Is(service.ErrOperationNotPermitted, err) {
		notPermittedResponse(c)
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, lead)
}
