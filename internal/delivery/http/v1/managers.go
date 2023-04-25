package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initManagersRoutes(api *gin.RouterGroup) {
	users := api.Group("/managers")
	{
		users.GET("/:id", h.getById)
	}
}

func (h *Handler) getById(c *gin.Context) {
	id := h.getAndValidateId(c, "id")

	userModel := h.getValidatedUser(c)

	if id == "" || userModel == nil {
		return
	}

	user, err := h.services.Managers.GetManagerById(c.Request.Context(), id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, user)
}
