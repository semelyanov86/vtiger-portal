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
