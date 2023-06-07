package service

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
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
	return h.repository.RetrieveById(ctx, id)
}

func (h Invoices) GetAll(ctx context.Context, filter repository.PaginationQueryFilter) ([]domain.Invoice, int, error) {
	invoices, err := h.repository.GetAll(ctx, filter)
	if err != nil {
		return invoices, 0, err
	}
	count, err := h.repository.Count(ctx, filter.Client)
	return invoices, count, err
}
