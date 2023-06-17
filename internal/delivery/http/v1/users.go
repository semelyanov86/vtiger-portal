package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/internal/service"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
	"net/http"
)

func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.POST("/", h.userSignUp)
		users.POST("/login", h.signIn)
		users.GET("/my", h.getUserInfo)
		users.GET("/settings", h.getUserSettings)
		users.PATCH("/settings", h.updateUserSettings)
		users.PUT("/my", h.updateUserInfo)
		users.GET("/my/documents", h.getUserDocuments)
		users.GET("/my/account", h.getAccountData)
		users.POST("/restore", h.sendRestoreToken)
		users.PUT("/password", h.resetPassword)
		users.GET("/all", h.usersFromAccount)
		users.GET("/:id/file/:file", h.getUserFile)
		users.DELETE("/:id/documents/:document", h.deleteUserFile)
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

	user, err := h.services.Users.SignUp(c.Request.Context(), inp, h.config)

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

type userSettingInput map[string]bool

func (h *Handler) updateUserSettings(c *gin.Context) {
	userModel := h.getValidatedUser(c)
	if userModel == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Auth Error", "field": "crmid", "message": "User is not found in auth process"})
		return
	}
	var inp userSettingInput
	if err := c.BindJSON(&inp); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": fieldErr.Field(), "message": fieldErr.Error()})
			return // exit on first error
		}
	}

	fields := h.config.Vtiger.Business.UserSettingsFields
	for _, field := range fields {
		value, ok := inp[field]
		if ok {
			err := h.services.Users.ChangeUserSetting(c.Request.Context(), userModel.Crmid, field, value)
			if err != nil {
				newResponse(c, http.StatusInternalServerError, err.Error())
				return
			}
		}
	}

	c.JSON(http.StatusAccepted, userModel)
}

func (h *Handler) updateUserInfo(c *gin.Context) {
	userModel := h.getValidatedUser(c)
	if userModel == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Auth Error", "field": "crmid", "message": "User is not found in auth process"})
		return
	}
	var inp service.UserUpdateInput
	if err := c.BindJSON(&inp); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": fieldErr.Field(), "message": fieldErr.Error()})
			return // exit on first error
		}
	}

	user, err := h.services.Users.Update(c.Request.Context(), userModel.Id, inp)

	if err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) {
			newResponse(c, http.StatusUnprocessableEntity, err.Error())

			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}
	go h.getUserInfo(c)
	c.JSON(http.StatusAccepted, user)
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
	userModel := h.getValidatedUser(c)
	if userModel == nil {
		return
	}
	user, err := h.services.Users.GetUserById(c.Request.Context(), userModel.Id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	res := AloneDataResponse[domain.User]{
		Data: *user,
	}
	c.JSON(http.StatusOK, res)
}

func (h Handler) getUserSettings(c *gin.Context) {
	userModel := h.getValidatedUser(c)
	if userModel == nil {
		return
	}
	settings, err := h.services.Users.GetUserSettings(c.Request.Context(), userModel.Crmid)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, settings)
}

func (h Handler) sendRestoreToken(c *gin.Context) {
	type UserEmailInput struct {
		Email string `json:"email" binding:"required,email,max=64"`
	}
	var inp UserEmailInput
	if err := c.BindJSON(&inp); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": fieldErr.Field(), "message": fieldErr.Error()})
			return
		}
	}
	err := h.services.Tokens.SendPasswordResetToken(c.Request.Context(), inp.Email)
	if errors.Is(repository.ErrRecordNotFound, err) {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": "email", "message": "User with this email not found"})
		return
	}
	if errors.Is(service.ErrUserIsNotActive, err) {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": "is_active", "message": "This user was disabled in portal"})
		return
	}

	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Token successfully created, please check an email"})
}

func (h Handler) resetPassword(c *gin.Context) {
	var inp service.PasswordResetInput
	if err := c.BindJSON(&inp); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": fieldErr.Field(), "message": fieldErr.Error()})
			return
		}
	}
	user, err := h.services.Users.ResetUserPassword(c.Request.Context(), inp)
	if errors.Is(repository.ErrRecordNotFound, err) {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": "email", "message": "User with this token not found"})
		return
	}
	if errors.Is(service.ErrUserIsNotActive, err) {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": "is_active", "message": "This user was disabled in portal"})
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusAccepted, user)
}

func (h Handler) usersFromAccount(c *gin.Context) {
	userModel := h.getValidatedUser(c)
	page, size := h.getPageAndSizeParams(c)

	if userModel == nil || page < 0 || size < 0 {
		newResponse(c, http.StatusBadRequest, "Wrong token or page size")
		return
	}
	users, count, err := h.services.Users.FindContactsFromAccount(c.Request.Context(), repository.PaginationQueryFilter{
		Page:     page,
		PageSize: size,
		Client:   userModel.AccountId,
		Contact:  userModel.Crmid,
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, DataResponse[domain.User]{
		Data:  users,
		Count: count,
		Page:  page,
		Size:  size,
	})
}

func (h *Handler) getUserDocuments(c *gin.Context) {
	userModel := h.getValidatedUser(c)

	if userModel == nil {
		return
	}

	documents, err := h.services.Documents.GetRelated(c.Request.Context(), userModel.Crmid)
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

func (h *Handler) getAccountData(c *gin.Context) {
	userModel := h.getValidatedUser(c)

	if userModel == nil {
		return
	}

	account, err := h.services.Accounts.GetAccountById(c.Request.Context(), userModel.AccountId)
	if errors.Is(service.ErrOperationNotPermitted, err) {
		notPermittedResponse(c)
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, AloneDataResponse[domain.Account]{
		Data: account,
	})
}

func (h *Handler) getUserFile(c *gin.Context) {
	id := h.getAndValidateId(c, "id")

	userModel := h.getValidatedUser(c)
	fileId := h.getAndValidateId(c, "file")

	if id == "" || userModel == nil || fileId == "" {
		return
	}

	if userModel.Crmid != id {
		notPermittedResponse(c)
		return
	}

	file, err := h.services.Documents.GetFile(c.Request.Context(), fileId, userModel.AccountId)

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

func (h *Handler) deleteUserFile(c *gin.Context) {
	id := h.getAndValidateId(c, "id")

	userModel := h.getValidatedUser(c)
	fileId := h.getAndValidateId(c, "document")

	if id == "" || userModel == nil || fileId == "" {
		return
	}

	if userModel.Crmid != id {
		notPermittedResponse(c)
		return
	}

	err := h.services.Documents.DeleteFile(c.Request.Context(), fileId, id)

	if errors.Is(service.ErrOperationNotPermitted, err) {
		notPermittedResponse(c)
		return
	}

	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
