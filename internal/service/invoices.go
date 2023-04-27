package service

import (
	"context"
	"errors"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
)

const CacheInvoiceTtl = 5000

type Invoices struct {
	repository repository.Invoice
	cache      cache.Cache
	module     ModulesService
	config     config.Config
}

func NewInvoiceService(repository repository.Invoice, cache cache.Cache, module ModulesService, config config.Config) Invoices {
	return Invoices{
		repository: repository,
		cache:      cache,
		module:     module,
		config:     config,
	}
}

func (h Invoices) GetInvoiceById(ctx context.Context, id string) (domain.Invoice, error) {
	invoice := &domain.Invoice{}
	err := GetFromCache[*domain.Invoice](id, invoice, h.cache)
	if err == nil {
		return *invoice, nil
	}

	if errors.Is(cache.ErrItemNotFound, err) {
		invoiceData, err := h.repository.RetrieveById(ctx, id)
		if err != nil {
			return invoiceData, e.Wrap("can not get a invoice", err)
		}
		err = StoreInCache[*domain.Invoice](id, &invoiceData, CacheInvoiceTtl, h.cache)
		if err != nil {
			return invoiceData, err
		}
		return invoiceData, nil
	} else {
		return *invoice, e.Wrap("can not convert caches data to invoice", err)
	}
}
