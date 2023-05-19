package service

import (
	"context"
	"errors"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
)

const CacheServicesTtl = 5000

type ServicesService struct {
	repository repository.Service
	cache      cache.Cache
	currency   CurrencyService
	module     ModulesService
	config     config.Config
}

func NewServicesService(repository repository.Service, cache cache.Cache, currency CurrencyService, module ModulesService, config config.Config) ServicesService {
	return ServicesService{
		repository: repository,
		cache:      cache,
		currency:   currency,
		module:     module,
		config:     config,
	}
}

func (s ServicesService) GetAll(ctx context.Context, filter repository.PaginationQueryFilter) ([]domain.Service, int, error) {
	services, err := s.repository.GetAll(ctx, filter)
	if err != nil {
		return services, 0, err
	}
	count, err := s.repository.Count(ctx, filter.Filters)
	if err != nil {
		return services, count, err
	}
	fullServices := make([]domain.Service, len(services))
	for i, service := range services {
		newService, err := s.GetServiceById(ctx, service.Id)
		if err != nil {
			return fullServices, count, e.Wrap("can not get service with id "+service.Id, err)
		}
		fullServices[i] = newService
	}
	return fullServices, count, err
}

func (s ServicesService) GetServiceById(ctx context.Context, id string) (domain.Service, error) {
	service := &domain.Service{}
	err := GetFromCache[*domain.Service](id, service, s.cache)
	if err == nil {
		return *service, nil
	}

	if errors.Is(cache.ErrItemNotFound, err) {
		serviceData, err := s.repository.RetrieveById(ctx, id)
		if err != nil {
			return serviceData, e.Wrap("can not get a service", err)
		}
		if serviceData.CurrencyId != "" {
			currency, err := s.currency.GetCurrencyById(ctx, serviceData.CurrencyId)
			if err != nil {
				return serviceData, e.Wrap("can not get a currency by id "+serviceData.CurrencyId, err)
			}
			serviceData.Currency = currency
		}
		err = StoreInCache[*domain.Service](id, &serviceData, CacheServicesTtl, s.cache)
		if err != nil {
			return serviceData, err
		}
		return serviceData, nil
	} else {
		return *service, e.Wrap("can not convert caches data to service", err)
	}
}
