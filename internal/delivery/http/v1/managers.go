package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (h *Handler) initManagersRoutes(api *gin.RouterGroup) {
	users := api.Group("/managers")
	{
		users.GET("/:id", h.getById)
	}
}

func (h *Handler) getById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		newResponse(c, http.StatusBadRequest, "code is empty")

		return
	}

	if !strings.Contains(id, "x") {
		newResponse(c, http.StatusUnprocessableEntity, "wrong id")

		return
	}

	userModel := h.services.Context.ContextGetUser(c)
	if userModel.Crmid == "" || userModel.Id < 1 {
		anonymousResponse(c)
		return
	}
	user, err := h.services.Managers.GetManagerById(c.Request.Context(), id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, user)
}
