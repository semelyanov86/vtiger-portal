package v1

import (
	"github.com/gin-gonic/gin"
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

	userModel := h.services.Context.ContextGetUser(c)
	if userModel.Crmid == "" || userModel.Id < 1 {
		anonymousResponse(c)
		return
	}
	module, err := h.services.Modules.Describe(c.Request.Context(), name)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, module)
}
