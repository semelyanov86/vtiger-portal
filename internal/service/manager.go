package service

import (
	"context"
	"errors"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
)

const CacheManagerTtl = 5000

type ManagerService struct {
	repository repository.Managers
	cache      cache.Cache
}

func NewManagerService(repository repository.Managers, cache cache.Cache) ManagerService {
	return ManagerService{
		repository: repository,
		cache:      cache,
	}
}

func (m ManagerService) GetManagerById(ctx context.Context, id string) (domain.Manager, error) {
	manager := &domain.Manager{}
	err := GetFromCache[*domain.Manager](id, manager, m.cache)
	if err == nil {
		return *manager, nil
	}

	if errors.Is(cache.ErrItemNotFound, err) {
		managerData, err := m.retrieveManager(ctx, id)
		if err != nil {
			return managerData, e.Wrap("can not get a manager", err)
		}
		err = StoreInCache[*domain.Manager](id, &managerData, CacheManagerTtl, m.cache)
		if err != nil {
			return managerData, err
		}
		return managerData, nil
	} else {
		return *manager, e.Wrap("can not convert caches data to manager", err)
	}
}

func (m ManagerService) retrieveManager(ctx context.Context, id string) (domain.Manager, error) {
	return m.repository.RetrieveById(ctx, id)
}
