package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
)

type CompanyCrm struct {
	vtiger vtiger.VtigerConnector
	config config.Config
}

func NewCompanyCrm(config config.Config, cache cache.Cache) CompanyCrm {
	return CompanyCrm{
		vtiger: vtiger.NewVtigerConnector(cache, config.Vtiger.Connection, vtiger.NewWebRequest(config.Vtiger.Connection)),
		config: config,
	}
}

func (m CompanyCrm) GetCompanyInfo(ctx context.Context) (domain.Company, error) {
	result, err := m.vtiger.Retrieve(ctx, m.config.Vtiger.Business.CompanyId)
	if err != nil {
		return domain.Company{}, e.Wrap("can not retrieve company with id "+m.config.Vtiger.Business.CompanyId, err)
	}
	company := domain.ConvertMapToCompany(result.Result)
	return company, nil
}
