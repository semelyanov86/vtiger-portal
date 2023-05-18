package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
	"strconv"
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

func (p ProductCrm) GetAll(ctx context.Context, filter PaginationQueryFilter) ([]domain.Product, error) {
	// Calculate the offset for the given page number and page size
	offset := (filter.Page - 1) * filter.PageSize
	isActive := p.getIsActiveFromFilter(filter.Filters)
	query := "SELECT id FROM Products WHERE discontinued = " + isActive + " LIMIT " + strconv.Itoa(offset) + ", " + strconv.Itoa(filter.PageSize) + ";"
	products := make([]domain.Product, 0)
	result, err := p.vtiger.Query(ctx, query)
	if err != nil {
		return products, e.Wrap("can not execute query "+query+", got error: "+result.Error.Message, err)
	}
	for _, data := range result.Result {
		product, err := domain.ConvertMapToProduct(data)
		if err != nil {
			return products, e.Wrap("can not convert map to product", err)
		}
		products = append(products, product)
	}
	return products, nil
}

func (p ProductCrm) Count(ctx context.Context, filters map[string]any) (int, error) {
	body := make(map[string]string)
	body["discontinued"] = p.getIsActiveFromFilter(filters)
	return p.vtiger.Count(ctx, "Products", body)
}

func (p ProductCrm) getIsActiveFromFilter(filters map[string]any) string {
	isActive := "1"
	if filters["discontinued"] == false || filters["discontinued"] == "false" {
		isActive = "0"
	}
	return isActive
}
