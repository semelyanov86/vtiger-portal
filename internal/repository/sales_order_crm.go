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

type SalesOrderCrm struct {
	vtiger vtiger.VtigerConnector
	config config.Config
}

func NewSalesOrderCrm(config config.Config, cache cache.Cache) SalesOrderCrm {
	return SalesOrderCrm{
		vtiger: vtiger.NewVtigerConnector(cache, config.Vtiger.Connection, vtiger.NewWebRequest(config.Vtiger.Connection)),
		config: config,
	}
}

func (m SalesOrderCrm) RetrieveById(ctx context.Context, id string) (domain.SalesOrder, error) {
	result, err := m.vtiger.Retrieve(ctx, id)
	if err != nil {
		return domain.SalesOrder{}, e.Wrap("can not retrieve sales order with id "+id+" got error", err)
	}
	return domain.ConvertMapToSalesOrder(result.Result)
}

func (m SalesOrderCrm) GetAll(ctx context.Context, filter PaginationQueryFilter) ([]domain.SalesOrder, error) {
	// Calculate the offset for the given page number and page size
	offset := (filter.Page - 1) * filter.PageSize
	query := "SELECT * FROM SalesOrder WHERE "

	sort := filter.Sort
	if sort == "" {
		sort = "-salesorder_no"
	}

	if filter.Search != "" {
		query += " salesorder_no LIKE '%" + filter.Search + "%' OR subject LIKE '%" + filter.Search + "%' OR sostatus LIKE '%" + filter.Search + "%' "
	} else {
		query += "account_id = " + filter.Client + " "
	}
	query += GenerateOrderByClause(sort)
	query += " LIMIT " + strconv.Itoa(offset) + ", " + strconv.Itoa(filter.PageSize) + ";"

	orders := make([]domain.SalesOrder, 0)
	result, err := m.vtiger.Query(ctx, query)
	if err != nil {
		return orders, e.Wrap("can not execute query "+query+", got error", err)
	}
	for _, data := range result.Result {
		salesOrder, err := domain.ConvertMapToSalesOrder(data)
		if err != nil {
			return orders, e.Wrap("can not convert map to sales order", err)
		}
		if salesOrder.AccountID == filter.Client {
			orders = append(orders, salesOrder)
		}
	}
	return orders, nil
}

func (m SalesOrderCrm) Count(ctx context.Context, client string) (int, error) {
	body := make(map[string]string)
	body["account_id"] = client
	return m.vtiger.Count(ctx, "SalesOrder", body)
}
