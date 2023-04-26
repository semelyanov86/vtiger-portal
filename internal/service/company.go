package service

import (
	"context"
	"errors"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
)

const CacheCompanyTtl = 50000

const CacheCompany = "COMPANY_INFO"

type Company struct {
	repository repository.Company
	cache      cache.Cache
}

func NewCompanyService(repository repository.Company, cache cache.Cache) Company {
	return Company{
		repository: repository,
		cache:      cache,
	}
}

func (c Company) GetCompany(ctx context.Context) (domain.Company, error) {
	company := &domain.Company{}
	err := GetFromCache[*domain.Company](CacheCompany, company, c.cache)
	if err == nil {
		return *company, nil
	}

	if errors.Is(cache.ErrItemNotFound, err) {
		companyData, err := c.retrieveCompany(ctx)
		if err != nil {
			return domain.Company{}, e.Wrap("can not get a company", err)
		}
		err = StoreInCache[*domain.Company](CacheCompany, &companyData, CacheCompanyTtl, c.cache)
		if err != nil {
			return domain.Company{}, err
		}
		return companyData, nil
	} else {
		return domain.Company{}, e.Wrap("can not convert caches data to company", err)
	}
}

func (c Company) retrieveCompany(ctx context.Context) (domain.Company, error) {
	return c.repository.GetCompanyInfo(ctx)
}
