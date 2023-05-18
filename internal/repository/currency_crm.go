package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
)

type CurrencyCrm struct {
	vtiger vtiger.VtigerConnector
	config config.Config
}

func NewCurrencyCrm(config config.Config, cache cache.Cache) CurrencyCrm {
	return CurrencyCrm{
		vtiger: vtiger.NewVtigerConnector(cache, config.Vtiger.Connection, vtiger.NewWebRequest(config.Vtiger.Connection)),
		config: config,
	}
}

func (c CurrencyCrm) RetrieveById(ctx context.Context, id string) (domain.Currency, error) {
	result, err := c.vtiger.Retrieve(ctx, id)
	if err != nil {
		return domain.Currency{}, e.Wrap("can not retrieve currency with id "+id, err)
	}
	return domain.ConvertMapToCurrency(result.Result)
}
