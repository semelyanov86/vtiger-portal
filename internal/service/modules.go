package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	cachedModuleData, err := m.cache.Get(name)
	if errors.Is(cache.ErrItemNotFound, err) || cachedModuleData == nil {
		moduleData, err := m.retrieveModule(ctx, name)
		if err != nil {
			return vtiger.Module{}, e.Wrap("can not get a module "+name, err)
		}
		cachedValue, err := json.Marshal(moduleData)
		if err != nil {
			return vtiger.Module{}, err
		}
		err = m.cache.Set(name, cachedValue, CacheModulesTtl)
		if err != nil {
			return vtiger.Module{}, err
		}
		return moduleData, nil
	} else {
		decodedModule := &vtiger.Module{}
		err = json.Unmarshal(cachedModuleData, decodedModule)
		if err != nil {
			if jsonErr, ok := err.(*json.SyntaxError); ok {
				problemPart := cachedModuleData[jsonErr.Offset-10 : jsonErr.Offset+10]

				err = fmt.Errorf("%w ~ error near '%s' (offset %d)", err, problemPart, jsonErr.Offset)
			}
			return vtiger.Module{}, e.Wrap("can not convert caches data to manager", err)
		}
		return *decodedModule, nil
	}
}

func (m ModulesService) retrieveModule(ctx context.Context, name string) (vtiger.Module, error) {
	return m.repository.GetModuleInfo(ctx, name)
}
