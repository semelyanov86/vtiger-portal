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

func (s ServiceContractCrm) Count(ctx context.Context, client string) (int, error) {
	body := make(map[string]string)
	body["sc_related_to"] = client
	return s.vtiger.Count(ctx, "ServiceContract", body)
}
