package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
	"strconv"
)

type ServiceContractCrm struct {
	vtiger vtiger.VtigerConnector
	config config.Config
}

func NewServiceContractCrm(config config.Config, cache cache.Cache) ServiceContractCrm {
	return ServiceContractCrm{
		vtiger: vtiger.NewVtigerConnector(cache, config.Vtiger.Connection, vtiger.NewWebRequest(config.Vtiger.Connection)),
		config: config,
	}
}

func (s ServiceContractCrm) RetrieveById(ctx context.Context, id string) (domain.ServiceContract, error) {
	result, err := s.vtiger.Retrieve(ctx, id)
	if err != nil {
		return domain.ServiceContract{}, e.Wrap("can not retrieve service contract with id "+id+" got error:"+result.Error.Message, err)
	}
	return domain.ConvertMapToServiceContract(result.Result)
}

func (s ServiceContractCrm) GetAll(ctx context.Context, filter vtiger.PaginationQueryFilter) ([]domain.ServiceContract, error) {
	// Calculate the offset for the given page number and page size
	offset := (filter.Page - 1) * filter.PageSize
	query := "SELECT * FROM ServiceContracts WHERE sc_related_to = " + filter.Client + " OR sc_related_to = " + filter.Contact + " LIMIT " + strconv.Itoa(offset) + ", " + strconv.Itoa(filter.PageSize) + ";"
	serviceContracts := make([]domain.ServiceContract, 0)
	result, err := s.vtiger.Query(ctx, query)
	if err != nil {
		return serviceContracts, e.Wrap("can not execute query "+query+", got error", err)
	}
	for _, data := range result.Result {
		serviceContract, err := domain.ConvertMapToServiceContract(data)
		if err != nil {
			return serviceContracts, e.Wrap("can not convert map to service contracts", err)
		}
		serviceContracts = append(serviceContracts, serviceContract)
	}
	return serviceContracts, nil
}

func (s ServiceContractCrm) Count(ctx context.Context, client string, contact string) (int, error) {
	body := make(map[string]string)
	body["sc_related_to"] = client
	body["_sc_related_to"] = contact

	return s.vtiger.Count(ctx, "ServiceContracts", body)
}
