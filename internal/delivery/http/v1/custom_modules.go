package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/internal/service"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
	"net/http"
)

func (h *Handler) initCustomModulesRoutes(api *gin.RouterGroup) {
	tickets := api.Group("/custom-modules")
	{
		tickets.GET("/:module", h.getAllEntities)
		tickets.GET("/:module/:id", h.getEntityById)
		tickets.POST("/:module", h.createCustomModule)
		tickets.PUT("/:module/:id", h.updateEntity)
		tickets.PATCH("/:module/:id", h.updatePartlyEntity)
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
		notPermittedResponse(c)
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
