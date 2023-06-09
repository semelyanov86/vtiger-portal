package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"net/http"
)

func (h *Handler) initStatisticsRoutes(api *gin.RouterGroup) {
	projects := api.Group("/statistics")
	{
		projects.GET("/", h.getAllStatistics)
		projects.GET("/tasks", h.getProgressTasks)
	}
}

func (h *Handler) getProgressTasks(c *gin.Context) {
	userModel := h.getValidatedUser(c)
	if userModel == nil {
		newResponse(c, http.StatusBadRequest, "no data about authenticated user")
		return
	}

	tasks, err := h.services.Statistics.GetInProgressTasks(c.Request.Context(), *userModel)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, DataResponse[domain.ProjectTask]{
		Data:  tasks,
		Count: len(tasks),
		Page:  1,
		Size:  100,
	})
}

func (h *Handler) getAllStatistics(c *gin.Context) {
	userModel := h.getValidatedUser(c)
	if userModel == nil {
		newResponse(c, http.StatusBadRequest, "no data about authenticated user")
		return
	}

	stat, err := h.services.Statistics.GetStatistics(c.Request.Context(), *userModel)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	res := AloneDataResponse[domain.Statistics]{
		Data: stat,
	}
	c.JSON(http.StatusOK, res)
}
