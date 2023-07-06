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

func (h *Handler) initCustomModulesRoutes(api *gin.RouterGroup) {
	custom := api.Group("/custom-modules")
	{
		custom.GET("/:module", h.getAllEntities)
		custom.GET("/:module/:id", h.getEntityById)
		custom.POST("/:module", h.createCustomModule)
		custom.PUT("/:module/:id", h.updateEntity)
		custom.PATCH("/:module/:id", h.updatePartlyEntity)
		custom.GET("/:module/:id/comments", h.getCustomComments)
		custom.POST("/:module/:id/comments", h.addCustomComment)
		custom.GET("/:module/:id/documents", h.getCustomDocuments)
		custom.POST("/:module/:id/documents", h.uploadCustomDocuments)
		custom.GET("/:module/:id/file/:file", h.getCustomFile)
	}
}

func (h *Handler) getAllEntities(c *gin.Context) {
	userModel := h.getValidatedUser(c)
	page, size := h.getPageAndSizeParams(c)

	if userModel == nil || page < 0 || size < 0 {
		newResponse(c, http.StatusBadRequest, "Wrong pagination params or auth user")
		return
	}
	moduleName := c.Param("module")
	if moduleName == "" {
		newResponse(c, http.StatusBadRequest, "module is empty")
		return
	}
	sortString := c.DefaultQuery("sort", "")

	records, count, err := h.services.CustomModules.GetAll(c.Request.Context(), vtiger.PaginationQueryFilter{
		Page:     page,
		PageSize: size,
		Client:   userModel.AccountId,
		Contact:  userModel.Crmid,
		Sort:     sortString,
		Search:   c.DefaultQuery("search", ""),
	}, moduleName)
	if errors.Is(service.ErrModuleNotSupported, err) {
		moduleNotSupportedResponse(c)
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, DataResponse[map[string]any]{
		Data:  records,
		Count: count,
		Page:  page,
		Size:  size,
	})
}

func (h *Handler) getEntityById(c *gin.Context) {
	userModel := h.getValidatedUser(c)
	if userModel == nil {
		newResponse(c, http.StatusBadRequest, "Wrong auth user")
		return
	}
	moduleName := c.Param("module")
	if moduleName == "" {
		newResponse(c, http.StatusBadRequest, "module is empty")
		return
	}
	id := h.getAndValidateId(c, "id")
	if id == "" {
		notPermittedResponse(c)
		return
	}
	result, err := h.services.CustomModules.GetById(c.Request.Context(), moduleName, id, *userModel)
	if errors.Is(repository.ErrRecordNotFound, err) {
		notPermittedResponse(c)
		return
	}
	if errors.Is(service.ErrModuleNotSupported, err) {
		moduleNotSupportedResponse(c)
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if errors.Is(service.ErrModuleNotSupported, err) {
		notPermittedResponse(c)
		return
	}
	res := AloneDataResponse[map[string]any]{
		Data: result,
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) createCustomModule(c *gin.Context) {
	moduleName := c.Param("module")
	if moduleName == "" {
		newResponse(c, http.StatusBadRequest, "module is empty")
		return
	}

	var inp map[string]any
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

	entity, err := h.services.CustomModules.CreateEntity(c.Request.Context(), inp, *userModel, moduleName)
	if errors.Is(service.ErrValidation, err) {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": "root", "message": err.Error()})
		return
	}
	if errors.Is(service.ErrModuleNotSupported, err) {
		notPermittedResponse(c)
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, entity)
}

func (h *Handler) updateEntity(c *gin.Context) {
	var inp map[string]any
	if err := c.BindJSON(&inp); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": fieldErr.Field(), "message": fieldErr.Error()})
			return // exit on first error
		}
	}
	id := h.getAndValidateId(c, "id")
	moduleName := c.Param("module")
	if moduleName == "" {
		newResponse(c, http.StatusBadRequest, "module is empty")
		return
	}
	userModel := h.getValidatedUser(c)
	if userModel == nil || id == "" {
		return
	}

	ticket, err := h.services.CustomModules.UpdateEntity(c.Request.Context(), inp, id, *userModel, moduleName)
	if errors.Is(service.ErrValidation, err) {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": "ticketcategories", "message": err.Error()})
		return
	}
	if errors.Is(service.ErrOperationNotPermitted, err) {
		notPermittedResponse(c)
		return
	}
	if errors.Is(service.ErrModuleNotSupported, err) {
		notPermittedResponse(c)
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusAccepted, ticket)
}

func (h *Handler) updatePartlyEntity(c *gin.Context) {
	var inp map[string]any
	if err := c.ShouldBindJSON(&inp); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": "root", "message": "Incorrect value"})
		return // exit on first error
	}
	id := h.getAndValidateId(c, "id")
	moduleName := c.Param("module")
	if moduleName == "" {
		newResponse(c, http.StatusBadRequest, "module is empty")
		return
	}
	userModel := h.getValidatedUser(c)
	if userModel == nil || id == "" {
		return
	}

	ticket, err := h.services.CustomModules.Revise(c.Request.Context(), inp, id, *userModel, moduleName)
	if errors.Is(service.ErrValidation, err) {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": "root", "message": err.Error()})
		return
	}
	if errors.Is(service.ErrOperationNotPermitted, err) {
		notPermittedResponse(c)
		return
	}
	if errors.Is(service.ErrModuleNotSupported, err) {
		notPermittedResponse(c)
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusAccepted, ticket)
}

func (h *Handler) getCustomComments(c *gin.Context) {
	id := h.getAndValidateId(c, "id")

	userModel := h.getValidatedUser(c)

	if id == "" || userModel == nil {
		return
	}

	moduleName := c.Param("module")
	if moduleName == "" {
		newResponse(c, http.StatusBadRequest, "module is empty")
		return
	}

	comments, err := h.services.CustomModules.GetRelatedComments(c.Request.Context(), id, moduleName, *userModel)
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

func (h *Handler) addCustomComment(c *gin.Context) {
	id := h.getAndValidateId(c, "id")

	userModel := h.getValidatedUser(c)

	if id == "" || userModel == nil {
		return
	}

	moduleName := c.Param("module")
	if moduleName == "" {
		newResponse(c, http.StatusBadRequest, "module is empty")
		return
	}

	var inp createCommentInput
	if err := c.BindJSON(&inp); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": fieldErr.Field(), "message": fieldErr.Error()})
			return // exit on first error
		}
	}
	comment, err := h.services.CustomModules.AddComment(c.Request.Context(), inp.Commentcontent, id, moduleName, *userModel)
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

func (h *Handler) getCustomDocuments(c *gin.Context) {
	id := h.getAndValidateId(c, "id")

	userModel := h.getValidatedUser(c)

	if id == "" || userModel == nil {
		return
	}

	moduleName := c.Param("module")
	if moduleName == "" {
		newResponse(c, http.StatusBadRequest, "module is empty")
		return
	}

	documents, err := h.services.CustomModules.GetRelatedDocuments(c.Request.Context(), id, moduleName, *userModel)
	if errors.Is(service.ErrOperationNotPermitted, err) {
		notPermittedResponse(c)
		return
	}
	if errors.Is(service.ErrModuleNotSupported, err) {
		moduleNotSupportedResponse(c)
		return
	}
	if errors.Is(repository.ErrRecordNotFound, err) {
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

func (h *Handler) uploadCustomDocuments(c *gin.Context) {
	id := h.getAndValidateId(c, "id")
	userModel := h.getValidatedUser(c)

	if id == "" || userModel == nil {
		notPermittedResponse(c)
		return
	}
	moduleName := c.Param("module")
	if moduleName == "" {
		newResponse(c, http.StatusBadRequest, "module is empty")
		return
	}

	_, err := h.services.CustomModules.GetById(c.Request.Context(), moduleName, id, *userModel)
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

func (h *Handler) getCustomFile(c *gin.Context) {
	id := h.getAndValidateId(c, "id")

	userModel := h.getValidatedUser(c)
	fileId := h.getAndValidateId(c, "file")

	if id == "" || userModel == nil || fileId == "" {
		return
	}

	moduleName := c.Param("module")
	if moduleName == "" {
		newResponse(c, http.StatusBadRequest, "module is empty")
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
