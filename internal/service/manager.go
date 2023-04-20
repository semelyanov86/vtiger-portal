package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	cachedManagerData, err := m.cache.Get(id)
	if errors.Is(cache.ErrItemNotFound, err) || cachedManagerData == nil {
		managerData, err := m.retrieveManager(ctx, id)
		if err != nil {
			return domain.Manager{}, e.Wrap("can not get a manager", err)
		}
		cachedValue, err := json.Marshal(managerData)
		if err != nil {
			return domain.Manager{}, err
		}
		err = m.cache.Set(id, cachedValue, CacheManagerTtl)
		if err != nil {
			return domain.Manager{}, err
		}
		return managerData, nil
	} else {
		decodedManager := &domain.Manager{}
		err = json.Unmarshal(cachedManagerData, decodedManager)
		if err != nil {
			if jsonErr, ok := err.(*json.SyntaxError); ok {
				problemPart := cachedManagerData[jsonErr.Offset-10 : jsonErr.Offset+10]

				err = fmt.Errorf("%w ~ error near '%s' (offset %d)", err, problemPart, jsonErr.Offset)
			}
			return domain.Manager{}, e.Wrap("can not convert caches data to manager", err)
		}
		return *decodedManager, nil
	}
}

func (m ManagerService) retrieveManager(ctx context.Context, id string) (domain.Manager, error) {
	return m.repository.RetrieveById(ctx, id)
}
