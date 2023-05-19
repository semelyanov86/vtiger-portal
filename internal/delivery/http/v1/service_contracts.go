package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"net/http"
)

func (h *Handler) initServiceContractsRoutes(api *gin.RouterGroup) {
	tickets := api.Group("/service-contracts")
	{
		tickets.GET("/", h.getAllServiceContracts)
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
	res := AloneDataResponse[domain.ServiceContract]{
		Data: serviceContract,
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) getAllServiceContracts(c *gin.Context) {
	userModel := h.getValidatedUser(c)
	page, size := h.getPageAndSizeParams(c)

	if userModel == nil || page < 0 || size < 0 {
		return
	}

	serviceContracts, count, err := h.services.ServiceContracts.GetAll(c.Request.Context(), repository.PaginationQueryFilter{
		Page:     page,
		PageSize: size,
		Client:   userModel.AccountId,
		Contact:  userModel.Crmid,
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, DataResponse[domain.ServiceContract]{
		Data:  serviceContracts,
		Count: count,
		Page:  page,
		Size:  size,
	})
}
