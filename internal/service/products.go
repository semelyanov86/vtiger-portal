package service

import (
	"context"
	"errors"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
)

const CacheProductTtl = 5000

type ProductService struct {
	repository repository.Product
	cache      cache.Cache
	currency   CurrencyService
	document   repository.Document
	module     ModulesService
	config     config.Config
}

func NewProductService(repository repository.Product, cache cache.Cache, currency CurrencyService, document repository.Document, module ModulesService, config config.Config) ProductService {
	return ProductService{
		repository: repository,
		cache:      cache,
		currency:   currency,
		document:   document,
		module:     module,
		config:     config,
	}
}

func (p ProductService) GetAll(ctx context.Context, filter vtiger.PaginationQueryFilter) ([]domain.Product, int, error) {
	products, err := p.repository.GetAll(ctx, filter)
	if err != nil {
		return products, 0, err
	}
	count, err := p.repository.Count(ctx, filter.Filters)
	if err != nil {
		return products, count, err
	}
	fullProducts := make([]domain.Product, len(products))
	for i, product := range products {
		newProduct, err := p.GetProductById(ctx, product.Id)
		if err != nil {
			return fullProducts, count, e.Wrap("can not get product with id "+product.Id, err)
		}
		fullProducts[i] = newProduct
	}
	return fullProducts, count, err
}

func (p ProductService) GetProductById(ctx context.Context, id string) (domain.Product, error) {
	product := &domain.Product{}
	err := GetFromCache[*domain.Product](id, product, p.cache)
	if err == nil {
		return *product, nil
	}

	if errors.Is(cache.ErrItemNotFound, err) {
		productData, err := p.repository.RetrieveById(ctx, id)
		if err != nil {
			return productData, e.Wrap("can not get a product", err)
		}
		if productData.CurrencyId != "" {
			currency, err := p.currency.GetCurrencyById(ctx, productData.CurrencyId)
			if err != nil {
				return productData, e.Wrap("can not get a currency by id "+productData.CurrencyId, err)
			}
			productData.Currency = currency
		}
		if productData.Imageattachmentids != "" {
			file, err := p.document.RetrieveFile(ctx, productData.Imageattachmentids)
			if err == nil && file.Filecontents != "" {
				productData.Imagecontent = file.Filecontents
			}
		}
		err = StoreInCache[*domain.Product](id, &productData, CacheProductTtl, p.cache)
		if err != nil {
			return productData, err
		}
		return productData, nil
	} else {
		return *product, e.Wrap("can not convert caches data to product", err)
	}
}
