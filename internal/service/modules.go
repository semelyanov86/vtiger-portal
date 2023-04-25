package service

import (
	"context"
	"errors"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
)

const CacheModulesTtl = 5000

type ModulesService struct {
	repository repository.Modules
	cache      cache.Cache
}

func NewModulesService(repository repository.Modules, cache cache.Cache) ModulesService {
	return ModulesService{
		repository: repository,
		cache:      cache,
	}
}

func (m ModulesService) Describe(ctx context.Context, name string) (vtiger.Module, error) {
	module := &vtiger.Module{}
	err := GetFromCache[*vtiger.Module](name, module, m.cache)
	if err == nil {
		return *module, nil
	}

	if errors.Is(cache.ErrItemNotFound, err) {
		moduleData, err := m.retrieveModule(ctx, name)
		if err != nil {
			return moduleData, e.Wrap("can not get a module", err)
		}
		err = StoreInCache[*vtiger.Module](name, &moduleData, CacheModulesTtl, m.cache)
		if err != nil {
			return moduleData, err
		}
		return moduleData, nil
	} else {
		return *module, e.Wrap("can not convert caches data to module", err)
	}
}

func (m ModulesService) retrieveModule(ctx context.Context, name string) (vtiger.Module, error) {
	return m.repository.GetModuleInfo(ctx, name)
}
