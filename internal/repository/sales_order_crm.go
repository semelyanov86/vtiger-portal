package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
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

func (m SalesOrderCrm) GetAll(ctx context.Context, filter vtiger.PaginationQueryFilter) ([]domain.SalesOrder, error) {
	items, err := m.vtiger.GetAll(ctx, filter, vtiger.QueryFieldsProps{
		DefaultSort:  "-salesorder_no",
		SearchFields: []string{"subject", "-salesorder_no", "sostatus"},
		ClientField:  "",
		AccountField: "account_id",
		TableName:    "SalesOrder",
	})
	if err != nil {
		return nil, err
	}
	orders := make([]domain.SalesOrder, 0, len(items))

	for _, data := range items {
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
