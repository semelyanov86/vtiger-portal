package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
)

type AccountCrm struct {
	vtiger vtiger.VtigerConnector
	config config.Config
}

func NewAccountCrm(config config.Config, cache cache.Cache) AccountCrm {
	return AccountCrm{
		vtiger: vtiger.NewVtigerConnector(cache, config.Vtiger.Connection, vtiger.NewWebRequest(config.Vtiger.Connection)),
		config: config,
	}
}

func (a AccountCrm) RetrieveById(ctx context.Context, id string) (domain.Account, error) {
	result, err := a.vtiger.Retrieve(ctx, id)
	if err != nil {
		return domain.Account{}, e.Wrap("can not retrieve account with id "+id, err)
	}
	account, err := domain.ConvertMapToAccount(result.Result)
	if err != nil {
		return account, e.Wrap("can not convert account to map", err)
	}
	return account, nil
}
