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
	cachedCompanyData, err := c.cache.Get(CacheCompany)
	if errors.Is(cache.ErrItemNotFound, err) || cachedCompanyData == nil {
		companyData, err := c.retrieveCompany(ctx)
		if err != nil {
			return domain.Company{}, e.Wrap("can not get a company", err)
		}
		cachedValue, err := json.Marshal(companyData)
		if err != nil {
			return domain.Company{}, err
		}
		err = c.cache.Set(CacheCompany, cachedValue, CacheManagerTtl)
		if err != nil {
			return domain.Company{}, err
		}
		return companyData, nil
	} else {
		decodedCompany := &domain.Company{}
		err = json.Unmarshal(cachedCompanyData, decodedCompany)
		if err != nil {
			if jsonErr, ok := err.(*json.SyntaxError); ok {
				problemPart := cachedCompanyData[jsonErr.Offset-10 : jsonErr.Offset+10]

				err = fmt.Errorf("%w ~ error near '%s' (offset %d)", err, problemPart, jsonErr.Offset)
			}
			return domain.Company{}, e.Wrap("can not convert caches data to company", err)
		}
		return *decodedCompany, nil
	}
}

func (c Company) retrieveCompany(ctx context.Context) (domain.Company, error) {
	return c.repository.GetCompanyInfo(ctx)
}
