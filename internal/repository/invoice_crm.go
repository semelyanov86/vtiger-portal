package repository

import (
	"context"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
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
		return domain.Invoice{}, e.Wrap("can not retrieve invoice desk with id "+id+" got error:"+result.Error.Message, err)
	}
	return domain.ConvertMapToInvoice(result.Result)
}
