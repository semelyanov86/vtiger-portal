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

func (s ServicesCrm) GetAll(ctx context.Context, filter PaginationQueryFilter) ([]domain.Service, error) {
	// Calculate the offset for the given page number and page size
	offset := (filter.Page - 1) * filter.PageSize
	isActive := GetIsActiveFromFilter(filter.Filters)
	query := "SELECT id FROM Services WHERE discontinued = " + isActive + " LIMIT " + strconv.Itoa(offset) + ", " + strconv.Itoa(filter.PageSize) + ";"
	services := make([]domain.Service, 0)
	result, err := s.vtiger.Query(ctx, query)
	if err != nil {
		return services, e.Wrap("can not execute query "+query+", got error: "+result.Error.Message, err)
	}
	for _, data := range result.Result {
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
