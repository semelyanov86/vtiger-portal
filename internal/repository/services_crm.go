package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
)

type ServicesCrm struct {
	vtiger vtiger.VtigerConnector
	config config.Config
}

func NewServicesCRM(config config.Config, cache cache.Cache) ServicesCrm {
	return ServicesCrm{
		vtiger: vtiger.NewVtigerConnector(cache, config.Vtiger.Connection, vtiger.NewWebRequest(config.Vtiger.Connection)),
		config: config,
	}
}

func (s ServicesCrm) RetrieveById(ctx context.Context, id string) (domain.Service, error) {
	result, err := s.vtiger.Retrieve(ctx, id)
	if err != nil {
		return domain.Service{}, e.Wrap("can not retrieve service with id "+id+" got error:"+result.Error.Message, err)
	}
	return domain.ConvertMapToService(result.Result)
}

func (s ServicesCrm) GetAll(ctx context.Context, filter vtiger.PaginationQueryFilter) ([]domain.Service, error) {
	items, err := s.vtiger.GetByWhereClause(ctx, filter, "discontinued", GetIsActiveFromFilter(filter.Filters), "Services")
	if err != nil {
		return nil, err
	}
	services := make([]domain.Service, 0, len(items))

	for _, data := range items {
		product, err := domain.ConvertMapToService(data)
		if err != nil {
			return services, e.Wrap("can not convert map to product", err)
		}
		services = append(services, product)
	}
	return services, nil
}

func (p ServicesCrm) Count(ctx context.Context, filters map[string]any) (int, error) {
	body := make(map[string]string)
	body["discontinued"] = GetIsActiveFromFilter(filters)
	return p.vtiger.Count(ctx, "Services", body)
}
