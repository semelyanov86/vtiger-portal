package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/internal/service"
	"net/http"
)

func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.POST("/", h.userSignUp)
		users.POST("/login", h.signIn)
		users.GET("/my", h.getUserInfo)
	}
}

func (h *Handler) userSignUp(c *gin.Context) {
	var inp service.UserSignUpInput
	if err := c.BindJSON(&inp); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": fieldErr.Field(), "message": fieldErr.Error()})
			return // exit on first error
		}
	}

	user, err := h.services.Users.SignUp(c.Request.Context(), inp)

	if err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) {
			newResponse(c, http.StatusUnprocessableEntity, err.Error())

			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}
	c.JSON(http.StatusCreated, user)
}

func (h *Handler) signIn(c *gin.Context) {
	var inp service.UserSignInInput
	if err := c.BindJSON(&inp); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": fieldErr.Field(), "message": fieldErr.Error()})
			return
		}
	}

	token, err := h.services.Tokens.CreateAuthToken(c.Request.Context(), inp.Email, inp.Password)

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": "email", "message": "User with this email not found"})

			return
		}
		if errors.Is(err, service.ErrPasswordDoesNotMatch) {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": "password", "message": "Password you passed to us is incorrect"})

			return
		}
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, token)
}

func (h Handler) getUserInfo(c *gin.Context) {
	userModel := h.services.Context.ContextGetUser(c)
	if userModel.Crmid == "" || userModel.Id < 1 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Anonymous Access", "field": "crmid", "message": "Got anonymous user from token or user without crmid"})
		return
	}
	user, err := h.services.Users.GetUserById(c.Request.Context(), userModel.Id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, user)
}