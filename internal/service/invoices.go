package service

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
)

type Invoices struct {
	repository repository.Invoice
	cache      cache.Cache
	module     ModulesService
	config     config.Config
	currency   CurrencyService
}

func NewInvoiceService(repository repository.Invoice, cache cache.Cache, module ModulesService, config config.Config, currency CurrencyService) Invoices {
	return Invoices{
		repository: repository,
		cache:      cache,
		module:     module,
		config:     config,
		currency:   currency,
	}
}

func (h Invoices) GetInvoiceById(ctx context.Context, id string) (domain.Invoice, error) {
	invoice, err := h.repository.RetrieveById(ctx, id)
	if err != nil {
		return invoice, err
	}
	if invoice.CurrencyID != "" {
		currency, err := h.currency.GetCurrencyById(ctx, invoice.CurrencyID)
		if err != nil {
			return invoice, e.Wrap("can not get a currency by id "+invoice.CurrencyID, err)
		}
		invoice.Currency = currency
	}
	return invoice, nil
}

func (h Invoices) GetAll(ctx context.Context, filter repository.PaginationQueryFilter) ([]domain.Invoice, int, error) {
	invoices, err := h.repository.GetAll(ctx, filter)
	if err != nil {
		return invoices, 0, err
	}
	count, err := h.repository.Count(ctx, filter.Client)
	for i, invoice := range invoices {
		if invoice.CurrencyID != "" {
			currency, err := h.currency.GetCurrencyById(ctx, invoice.CurrencyID)
			if err != nil {
				return invoices, count, e.Wrap("can not get a currency by id "+invoice.CurrencyID, err)
			}
			invoices[i].Currency = currency
		}
	}
	return invoices, count, err
}
