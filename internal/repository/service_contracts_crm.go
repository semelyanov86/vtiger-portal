package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
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
	items, err := s.vtiger.GetAll(ctx, filter, vtiger.QueryFieldsProps{
		DefaultSort:  "-contract_no",
		SearchFields: []string{"subject", "contract_no", "contract_type", "contract_status"},
		ClientField:  "sc_related_to",
		AccountField: "sc_related_to",
		TableName:    "ServiceContracts",
	})
	if err != nil {
		return nil, err
	}

	contracts := make([]domain.ServiceContract, 0, len(items))

	for _, data := range items {
		contract, err := domain.ConvertMapToServiceContract(data)
		if err != nil {
			return contracts, e.Wrap("can not convert map to contract", err)
		}
		if contract.ScRelatedTo == filter.Client || contract.ScRelatedTo == filter.Contact {
			contracts = append(contracts, contract)
		}
	}
	return contracts, nil
}

func (s ServiceContractCrm) Count(ctx context.Context, client string, contact string) (int, error) {
	body := make(map[string]string)
	body["sc_related_to"] = client
	body["_sc_related_to"] = contact

	return s.vtiger.Count(ctx, "ServiceContracts", body)
}
