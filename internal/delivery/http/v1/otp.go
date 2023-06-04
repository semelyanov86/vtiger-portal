package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/semelyanov86/vtiger-portal/internal/service"
	"net/http"
)

func (h *Handler) initOtpRoutes(api *gin.RouterGroup) {
	users := api.Group("/otp")
	{
		users.GET("/generate", h.generateOtp)
		users.PATCH("/verify", h.verifyOtp)
	}
}

func (h *Handler) generateOtp(c *gin.Context) {
	userModel := h.getValidatedUser(c)

	if userModel == nil {
		newResponse(c, http.StatusUnauthorized, "User with this ID not found")
		return
	}
	payload := service.OTPInput{UserId: userModel.Id}

	otp, err := h.services.Auth.GenerateOtp(c.Request.Context(), payload)
	if errors.Is(service.ErrUserNotFound, err) {
		newResponse(c, http.StatusUnauthorized, "User with this ID not found")
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, otp)
}

func (h *Handler) verifyOtp(c *gin.Context) {
	userModel := h.getValidatedUser(c)

	if userModel == nil {
		newResponse(c, http.StatusUnauthorized, "User with this ID not found")
		return
	}

	var payload *service.OTPInput

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "message": err.Error()})
		return
	}
	payload.UserId = userModel.Id
	result, err := h.services.Auth.VerifyOtp(c.Request.Context(), *payload)
	if errors.Is(service.ErrTokenNotExist, err) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token Error", "message": err.Error()})
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusAccepted, result)
}
