package service

import (
	"context"
	"errors"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
)

const CacheAccountTtl = 5000

type AccountService struct {
	repository repository.Account
	cache      cache.Cache
}

func NewAccountService(repository repository.Account, cache cache.Cache) AccountService {
	return AccountService{
		repository: repository,
		cache:      cache,
	}
}

func (a AccountService) GetAccountById(ctx context.Context, id string) (domain.Account, error) {
	account := &domain.Account{}
	err := GetFromCache[*domain.Account](id, account, a.cache)
	if err == nil {
		return *account, nil
	}

	if errors.Is(cache.ErrItemNotFound, err) {
		accountData, err := a.repository.RetrieveById(ctx, id)
		if err != nil {
			return accountData, e.Wrap("can not get a account", err)
		}
		err = StoreInCache[*domain.Account](id, &accountData, CacheAccountTtl, a.cache)
		if err != nil {
			return accountData, err
		}
		return accountData, nil
	} else {
		return *account, e.Wrap("can not convert caches data to account", err)
	}
}
