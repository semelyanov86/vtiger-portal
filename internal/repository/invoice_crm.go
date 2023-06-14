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

type InvoiceCrm struct {
	vtiger vtiger.VtigerConnector
	config config.Config
}

func NewInvoiceCrm(config config.Config, cache cache.Cache) InvoiceCrm {
	return InvoiceCrm{
		vtiger: vtiger.NewVtigerConnector(cache, config.Vtiger.Connection, vtiger.NewWebRequest(config.Vtiger.Connection)),
		config: config,
	}
}

func (m InvoiceCrm) RetrieveById(ctx context.Context, id string) (domain.Invoice, error) {
	result, err := m.vtiger.Retrieve(ctx, id)
	if err != nil {
		return domain.Invoice{}, e.Wrap("can not retrieve invoice desk with id "+id+" got error", err)
	}
	return domain.ConvertMapToInvoice(result.Result)
}

func (m InvoiceCrm) GetAll(ctx context.Context, filter PaginationQueryFilter) ([]domain.Invoice, error) {
	// Calculate the offset for the given page number and page size
	offset := (filter.Page - 1) * filter.PageSize
	query := "SELECT * FROM Invoice WHERE "

	sort := filter.Sort
	if sort == "" {
		sort = "-ticket_no"
	}

	if filter.Search != "" {
		query += " invoice_no LIKE '%" + filter.Search + "%' OR subject LIKE '%" + filter.Search + "%' OR invoicestatus LIKE '%" + filter.Search + "%' "
	} else {
		query += "account_id = " + filter.Client + " "
	}
	query += GenerateOrderByClause(sort)
	query += " LIMIT " + strconv.Itoa(offset) + ", " + strconv.Itoa(filter.PageSize) + ";"

	invoices := make([]domain.Invoice, 0)
	result, err := m.vtiger.Query(ctx, query)
	if err != nil {
		return invoices, e.Wrap("can not execute query "+query+", got error", err)
	}
	for _, data := range result.Result {
		invoice, err := domain.ConvertMapToInvoice(data)
		if err != nil {
			return invoices, e.Wrap("can not convert map to invoice", err)
		}
		if invoice.AccountID == filter.Client {
			invoices = append(invoices, invoice)
		}
	}
	return invoices, nil
}

func (m InvoiceCrm) Count(ctx context.Context, client string) (int, error) {
	body := make(map[string]string)
	body["account_id"] = client
	return m.vtiger.Count(ctx, "Invoice", body)
}
