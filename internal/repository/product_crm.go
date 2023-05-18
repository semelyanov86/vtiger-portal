package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
)

type ProductCrm struct {
	vtiger vtiger.VtigerConnector
	config config.Config
}

func NewProductCrm(config config.Config, cache cache.Cache) ProductCrm {
	return ProductCrm{
		vtiger: vtiger.NewVtigerConnector(cache, config.Vtiger.Connection, vtiger.NewWebRequest(config.Vtiger.Connection)),
		config: config,
	}
}

func (p ProductCrm) RetrieveById(ctx context.Context, id string) (domain.Product, error) {
	result, err := p.vtiger.Retrieve(ctx, id)
	if err != nil {
		return domain.Product{}, e.Wrap("can not retrieve product with id "+id+" got error:"+result.Error.Message, err)
	}
	return domain.ConvertMapToProduct(result.Result)
}
