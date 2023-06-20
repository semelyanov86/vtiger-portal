package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"net/http"
)

func (h *Handler) initCustomModulesRoutes(api *gin.RouterGroup) {
	tickets := api.Group("/custom-modules")
	{
		tickets.GET("/:module", h.getAllEntities)
		tickets.GET("/:module/:id", h.getEntityById)
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

	records, count, err := h.services.CustomModules.GetAll(c.Request.Context(), repository.PaginationQueryFilter{
		Page:     page,
		PageSize: size,
		Client:   userModel.AccountId,
		Contact:  userModel.Crmid,
		Sort:     sortString,
		Search:   c.DefaultQuery("search", ""),
	}, moduleName)
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
	res := AloneDataResponse[map[string]any]{
		Data: result,
	}
	c.JSON(http.StatusOK, res)
}
