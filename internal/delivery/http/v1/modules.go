package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
	"net/http"
)

func (h *Handler) initModulesRoutes(api *gin.RouterGroup) {
	users := api.Group("/modules")
	{
		users.GET("/:name", h.describeModule)
	}
}

func (h *Handler) describeModule(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		newResponse(c, http.StatusBadRequest, "name is empty")

		return
	}

	userModel := h.getValidatedUser(c)
	if userModel == nil {
		return
	}
	module, err := h.services.Modules.Describe(c.Request.Context(), name)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	res := AloneDataResponse[vtiger.Module]{
		Data: module,
	}
	c.JSON(http.StatusOK, res)
}
