package service

import (
	"context"
	"errors"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
)

const CacheCurrencyTtl = 500000

type CurrencyService struct {
	repository repository.Currency
	cache      cache.Cache
}

func NewCurrencyService(repository repository.Currency, cache cache.Cache) CurrencyService {
	return CurrencyService{
		repository: repository,
		cache:      cache,
	}
}

func (c CurrencyService) GetCurrencyById(ctx context.Context, id string) (domain.Currency, error) {
	currency := &domain.Currency{}
	err := GetFromCache[*domain.Currency](id, currency, c.cache)
	if err == nil {
		return *currency, nil
	}

	if errors.Is(cache.ErrItemNotFound, err) {
		currencyData, err := c.retrieveCurrency(ctx, id)
		if err != nil {
			return currencyData, e.Wrap("can not get a currency", err)
		}
		err = StoreInCache[*domain.Currency](id, &currencyData, CacheCurrencyTtl, c.cache)
		if err != nil {
			return currencyData, err
		}
		return currencyData, nil
	} else {
		return *currency, e.Wrap("can not convert caches data to currency", err)
	}
}

func (c CurrencyService) retrieveCurrency(ctx context.Context, id string) (domain.Currency, error) {
	return c.repository.RetrieveById(ctx, id)
}
