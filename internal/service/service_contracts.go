package service

import (
	"context"
	"errors"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
)

const CacheServiceContractTtl = 5000

type ServiceContracts struct {
	repository repository.ServiceContract
	cache      cache.Cache
	document   DocumentServiceInterface
	module     ModulesService
	config     config.Config
}

func NewServiceContractsService(repository repository.ServiceContract, cache cache.Cache, document DocumentServiceInterface, module ModulesService, config config.Config) ServiceContracts {
	return ServiceContracts{
		repository: repository,
		cache:      cache,
		document:   document,
		module:     module,
		config:     config,
	}
}

func (s ServiceContracts) GetServiceContractById(ctx context.Context, id string) (domain.ServiceContract, error) {
	serviceContract := &domain.ServiceContract{}
	err := GetFromCache[*domain.ServiceContract](id, serviceContract, s.cache)
	if err == nil {
		return *serviceContract, nil
	}

	if errors.Is(cache.ErrItemNotFound, err) {
		scData, err := s.repository.RetrieveById(ctx, id)
		if err != nil {
			return scData, e.Wrap("can not get a serviceContract", err)
		}
		err = StoreInCache[*domain.ServiceContract](id, &scData, CacheServiceContractTtl, s.cache)
		if err != nil {
			return scData, err
		}
		return scData, nil
	} else {
		return *serviceContract, e.Wrap("can not convert caches data to serviceContract", err)
	}
}

func (s ServiceContracts) GetAll(ctx context.Context, filter vtiger.PaginationQueryFilter) ([]domain.ServiceContract, int, error) {
	serviceContracts, err := s.repository.GetAll(ctx, filter)
	if err != nil {
		return serviceContracts, 0, err
	}
	count, err := s.repository.Count(ctx, filter.Client, filter.Contact)
	return serviceContracts, count, err
}
