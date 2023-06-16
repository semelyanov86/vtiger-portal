package service

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
)

type SalesOrders struct {
	repository        repository.SalesOrderCrm
	invoiceRepository repository.Invoice
	cache             cache.Cache
	module            ModulesService
	config            config.Config
	currency          CurrencyService
}

func NewSalesOrderService(repository repository.SalesOrderCrm, cache cache.Cache, module ModulesService, config config.Config, currency CurrencyService, invoice repository.Invoice) SalesOrders {
	return SalesOrders{
		repository:        repository,
		cache:             cache,
		module:            module,
		config:            config,
		currency:          currency,
		invoiceRepository: invoice,
	}
}

func (h SalesOrders) GetSalesOrderById(ctx context.Context, id string) (domain.SalesOrder, error) {
	salesOrder, err := h.repository.RetrieveById(ctx, id)
	if err != nil {
		return salesOrder, err
	}
	if salesOrder.CurrencyID != "" {
		currency, err := h.currency.GetCurrencyById(ctx, salesOrder.CurrencyID)
		if err != nil {
			return salesOrder, e.Wrap("can not get a currency by id "+salesOrder.CurrencyID, err)
		}
		salesOrder.Currency = currency
	}
	invoices, err := h.invoiceRepository.GetFromSalesOrder(ctx, id)
	if err != nil {
		return salesOrder, e.Wrap("can not get invoices from sales order", err)
	}
	salesOrder.Invoices = invoices
	return salesOrder, nil
}

func (h SalesOrders) GetAll(ctx context.Context, filter repository.PaginationQueryFilter) ([]domain.SalesOrder, int, error) {
	salesOrders, err := h.repository.GetAll(ctx, filter)
	if err != nil {
		return salesOrders, 0, err
	}
	count, err := h.repository.Count(ctx, filter.Client)
	for i, salesOrder := range salesOrders {
		if salesOrder.CurrencyID != "" {
			currency, err := h.currency.GetCurrencyById(ctx, salesOrder.CurrencyID)
			if err != nil {
				return salesOrders, count, e.Wrap("can not get a currency by id "+salesOrder.CurrencyID, err)
			}
			salesOrders[i].Currency = currency
		}
	}
	return salesOrders, count, err
}
