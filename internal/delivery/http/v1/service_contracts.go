package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initServiceContractsRoutes(api *gin.RouterGroup) {
	tickets := api.Group("/service-contracts")
	{
		tickets.GET("/:id", h.getServiceContract)
	}
}

func (h *Handler) getServiceContract(c *gin.Context) {
	id := h.getAndValidateId(c, "id")

	userModel := h.getValidatedUser(c)
	if userModel == nil || id == "" {
		return
	}

	serviceContract, err := h.services.ServiceContracts.GetServiceContractById(c.Request.Context(), id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if userModel.AccountId != serviceContract.ScRelatedTo && userModel.Crmid != serviceContract.ScRelatedTo {
		notPermittedResponse(c)
		return
	}

	c.JSON(http.StatusOK, serviceContract)
}
